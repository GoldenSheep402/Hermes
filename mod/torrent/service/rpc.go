package service

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/torrent/dao"
	"github.com/GoldenSheep402/Hermes/mod/torrent/model"
	torrentV1 "github.com/GoldenSheep402/Hermes/pkg/proto/torrent/v1"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/GoldenSheep402/Hermes/pkg/torrent"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ torrentV1.TorrentServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	torrentV1.UnimplementedTorrentServiceServer
}

func (s *S) UploadTorrent(ctx context.Context, req *torrentV1.UploadTorrentRequest) (*torrentV1.UploadTorrentResponse, error) {
	if len(req.TorrentData) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Empty torrent data")
	}

	// 1. Parse torrent
	parsed, err := torrent.Parse(req.TorrentData)
	if err != nil {
		s.Log.Errorw("failed to parse torrent", "error", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid torrent file")
	}

	// 2. Validate hash doesn't already exist
	existing, err := dao.Torrent.GetByHash(ctx, parsed.InfoHash)
	if err == nil && existing != nil {
		return nil, status.Error(codes.AlreadyExists, "Torrent already exists")
	} else if err != nil && status.Code(err) != codes.NotFound {
		s.Log.Errorw("failed to check existing torrent", "error", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	// TODO: Get real UploaderID from context (JWT/Auth module)
	uploaderID := "TODO_UPLOADER_ID"

	// 3. Map to Model
	torrentModel := &model.Torrent{
		Model:        stdao.Model{ID: ulid.Make().String()},
		InfoHash:     parsed.InfoHash,
		UploaderID:   uploaderID,
		Name:         parsed.Name,
		Size:         parsed.Size,
		PieceLength:  parsed.PieceLength,
		PieceCount:   parsed.PieceCount,
		IsSingleFile: len(parsed.Files) == 1,
		FileCount:    len(parsed.Files),
		IsActive:     true,
	}

	// 4. Map files
	var files []model.TorrentFile
	for _, f := range parsed.Files {
		files = append(files, model.TorrentFile{
			Model:     stdao.Model{ID: ulid.Make().String()},
			TorrentID: torrentModel.ID,
			Path:      f.Path,
			Size:      f.Size,
		})
	}

	// 5. Save to DB
	id, err := dao.Torrent.Create(ctx, torrentModel, files)
	if err != nil {
		s.Log.Errorw("failed to create torrent", "error", err)
		return nil, err // DAO returns proper grpc status
	}

	// 6. Return response
	return &torrentV1.UploadTorrentResponse{
		TorrentId: id,
		InfoHash:  parsed.InfoHash,
	}, nil
}
