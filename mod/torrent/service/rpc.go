package service

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/GoldenSheep402/Hermes/conf"
	categoryDao "github.com/GoldenSheep402/Hermes/mod/category/dao"
	torrentDao "github.com/GoldenSheep402/Hermes/mod/torrent/dao"
	torrentModel "github.com/GoldenSheep402/Hermes/mod/torrent/model"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	torrentV1 "github.com/GoldenSheep402/Hermes/pkg/proto/torrent/v1"
	"github.com/anacrolix/torrent/bencode"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"strings"
	"time"
)

var _ torrentV1.TorrentServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	torrentV1.UnimplementedTorrentServiceServer
}

func (s *S) GetTorrentV1(ctx context.Context, req *torrentV1.GetTorrentV1Request) (*torrentV1.GetTorrentV1Response, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	// TODO: rbac

	metadata, err := torrentDao.Torrent.GetTorrentMetadata(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var resp torrentV1.GetTorrentV1Response
	for _, _meta := range metadata {
		resp.Metadata = append(resp.Metadata, &torrentV1.TorrentMetaData{
			Id:          _meta.ID,
			CategoryId:  _meta.CategoryID,
			TorrentId:   req.Id,
			Key:         _meta.Key,
			Order:       int32(_meta.Order),
			Description: _meta.Description,
			Type:        _meta.Type,
			Value:       _meta.Value,
		})
	}

	return &resp, nil
}

func (s *S) GetTorrentV1List(ctx context.Context, req *torrentV1.GetTorrentV1ListRequest) (*torrentV1.GetTorrentV1ListResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	// TODO: rbac
	torrents, err := torrentDao.Torrent.GetTorrentList(ctx, req.CategoryId, req.Id, int(req.Limit))
	if err != nil {
		return nil, err
	}

	var resp torrentV1.GetTorrentV1ListResponse
	for _, _torrent := range torrents {
		resp.Torrents = append(resp.Torrents, &torrentV1.TorrentBasic{
			Id:   _torrent.ID,
			Name: _torrent.Name,
			// Description:  _torrent.Description,
			CategoryId: _torrent.CategoryID,
		})
		category, _, err := categoryDao.Category.Get(ctx, _torrent.CategoryID)
		if err != nil {
			return nil, err
		}
		resp.Torrents[len(resp.Torrents)-1].CategoryName = category.Name
	}

	return &resp, nil
}

func (s *S) CreateTorrentV1(ctx context.Context, req *torrentV1.CreateTorrentV1Request) (*torrentV1.CreateTorrentV1Response, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	// TODO: rbac

	user, err := userDao.User.GetInfo(ctx, UID)
	if err != nil {
		return nil, err
	}

	// Here torrent.Data is byte[] type
	decoder := bencode.NewDecoder(bytes.NewReader(req.Torrent.Data))
	bencodeTorrent := &torrentModel.BencodeTorrent{}
	err = decoder.Decode(bencodeTorrent)
	if err != nil {
		return nil, err
	}

	marshaledInfo, err := bencode.Marshal(bencodeTorrent.Info)
	if err != nil {
		return nil, err
	}

	trackerAddress := conf.Get().TrackerV1.Endpoint

	now := time.Now()
	isPrivate := true
	var files []torrentModel.File
	torrent := &torrentModel.Torrent{
		CategoryID: req.CategoryId,

		InfoHash:     fmt.Sprintf("%x", sha1.Sum(marshaledInfo)),
		CreatorID:    UID,
		Announce:     trackerAddress + "/tracker/announce/key/" + user.Key,
		CreatedBy:    bencodeTorrent.CreatedBy,
		CreationDate: &now,
		// Name:         bencodeTorrent.Info.Name,
		// NameUTF8:     *bencodeTorrent.Info.NameUTF8,
		Length:      bencodeTorrent.Info.Length,
		Pieces:      []byte(bencodeTorrent.Info.Pieces),
		PieceLength: bencodeTorrent.Info.PieceLength,
		Private:     &isPrivate,
		Source:      bencodeTorrent.Info.Source,
		Md5sum:      bencodeTorrent.Info.Md5sum,
	}

	// TODO
	if req.Name != "" {
		torrent.Name = req.Name
	}

	//if bencodeTorrent.Info.NameUTF8 != nil {
	//	if *bencodeTorrent.Info.NameUTF8 == "" {
	//		torrent.NameUTF8 = torrent.Name
	//	} else {
	//		torrent.NameUTF8 = *bencodeTorrent.Info.NameUTF8
	//	}
	//} else {
	//	torrent.NameUTF8 = torrent.Name
	//}

	// If the torrent is a single file torrent
	if bencodeTorrent.Info.Files == nil {
		path := ""
		if torrent.Name != "" {
			path = strings.Join([]string{bencodeTorrent.Info.Name}, "/")
		}

		pathUTF8 := ""
		if torrent.NameUTF8 != "" {
			pathUTF8 = strings.Join([]string{*bencodeTorrent.Info.NameUTF8}, "/")
		}

		files = append(files, torrentModel.File{
			TorrentID: torrent.InfoHash,
			Length:    *torrent.Length,
			Path:      path,
			PathUTF8:  pathUTF8,
		})
	} else {
		for _, fileInfo := range *bencodeTorrent.Info.Files {
			path := ""
			if torrent.Name == "" {
				path = strings.Join(fileInfo.Path, "/")
			} else {
				path = strings.Join(fileInfo.Path, "/")
			}

			pathUTF8 := ""
			if bencodeTorrent.Info.NameUTF8 != nil {
				pathUTF8 = strings.Join([]string{*bencodeTorrent.Info.NameUTF8}, "/")
			}
			files = append(files, torrentModel.File{
				TorrentID: torrent.InfoHash,
				Length:    fileInfo.Length,
				Path:      path,
				PathUTF8:  pathUTF8,
			})
		}
	}

	metas := make([]torrentModel.TorrentMetadata, len(req.Metadata))
	for i, meta := range req.Metadata {
		metas[i] = torrentModel.TorrentMetadata{
			CategoryID: req.CategoryId,
			MetadataID: meta.Id,
			Value:      meta.Value,
		}
	}

	// Create the torrent
	// Link with files and metas
	id, err := torrentDao.Torrent.Create(ctx, torrent, files, metas)
	if err != nil {
		return nil, err
	}

	return &torrentV1.CreateTorrentV1Response{
		Id: id,
	}, nil
}

func (s *S) DownloadTorrentV1(ctx context.Context, req *torrentV1.DownloadTorrentV1Request) (*torrentV1.DownloadTorrentV1Response, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	// TODO: rbac

	user, err := userDao.User.GetInfo(ctx, UID)
	if err != nil {
		return nil, err
	}

	torrent, torrentFile, err := torrentDao.Torrent.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	_torrentFull := &torrentModel.BencodeTorrent{
		Announce:  conf.Get().TrackerV1.Endpoint + user.Key,
		CreatedBy: torrent.CreatedBy,
		CreatedAt: func(i int64) *int { v := int(i); return &v }(torrent.CreatedAt.Unix()),
		Info: torrentModel.BencodeInfo{
			Files: func() *[]torrentModel.FileInfo {
				if len(torrentFile) == 0 {
					return nil
				}
				files := make([]torrentModel.FileInfo, len(torrentFile))
				for i, file := range torrentFile {
					var path []string
					if file.Path != "" {
						path = strings.Split(file.Path, "/")
					}

					var pathUTF8 []string
					if file.PathUTF8 != "" {
						pathUTF8 = strings.Split(file.PathUTF8, "/")
					}

					files[i] = torrentModel.FileInfo{
						Length:   file.Length,
						Path:     path,
						PathUTF8: pathUTF8,
					}
				}
				return &files
			}(),
			Name:        torrent.Name,
			Length:      torrent.Length,
			Md5sum:      torrent.Md5sum,
			Pieces:      string(torrent.Pieces),
			PieceLength: torrent.PieceLength,
			Private: func(b *bool) *int {
				if b == nil {
					return nil
				}
				v := 0
				if *b {
					v = 1
				}
				return &v
			}(torrent.Private),
			Source: torrent.Source,
		},
	}

	if torrent.NameUTF8 != "" {
		_torrentFull.Info.NameUTF8 = &torrent.NameUTF8
	}

	marshaledInfo, err := bencode.Marshal(_torrentFull.Info)
	if err != nil {
		return nil, err
	}
	hash := fmt.Sprintf("%x", sha1.Sum(marshaledInfo))

	var buf bytes.Buffer
	encoder := bencode.NewEncoder(&buf)
	err = encoder.Encode(_torrentFull)
	if err != nil {
		return nil, err
	}

	filePath := "/tmp/torrent_debug.torrent"
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// 将编码后的数据写入文件
	_, err = file.Write(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to write to file: %v", err)
	}

	base64Data := base64.StdEncoding.EncodeToString(buf.Bytes())

	return &torrentV1.DownloadTorrentV1Response{
		Data: base64Data,
	}, nil
}
