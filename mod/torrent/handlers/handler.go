package handlers

import (
	"bytes"
	"context"
	systemDao "github.com/GoldenSheep402/Hermes/mod/system/dao"
	torrentDao "github.com/GoldenSheep402/Hermes/mod/torrent/dao"
	torrentModel "github.com/GoldenSheep402/Hermes/mod/torrent/model"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/anacrolix/torrent/bencode"
	"github.com/juanjiTech/jin"
	"net/http"
	"strings"
)

func Registry(jinE *jin.Engine) {
	trackerGroup := jinE.Group("/torrent")

	trackerGroup.Group("/download").
		GET("/:key", DownloadHandler)
}

func DownloadHandler(c *jin.Context) {
	ctx := context.Background()
	key, ok := c.Params.Get("key")
	if !ok {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := userDao.User.CheckKey(ctx, key)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	torrentID := c.Request.URL.Query().Get("id")

	torrent, torrentFile, err := torrentDao.Torrent.Get(ctx, torrentID)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	_torrentFull := &torrentModel.BencodeTorrent{
		//Announce:  conf.Get().TrackerV1.Endpoint + key,
		CreatedBy: torrent.CreatedBy,
		Comment:   torrent.Comment,
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

	var announceList []string

	trackers, err := systemDao.InnetTracker.GetTrackers(ctx)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, tracker := range trackers {
		if tracker.Enable {
			announceList = append(announceList, tracker.Address+"/tracker/announce/key/"+user.Key)
		}
	}

	if len(announceList) > 0 {
		_torrentFull.AnnounceList = &[][]string{announceList}
	}

	if torrent.NameUTF8 != "" {
		_torrentFull.Info.NameUTF8 = &torrent.NameUTF8
	}

	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var buf bytes.Buffer
	encoder := bencode.NewEncoder(&buf)
	err = encoder.Encode(_torrentFull)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	c.Writer.Header().Set("Content-Type", "application/x-bittorrent")
	c.Writer.Header().Set("Content-Disposition", `attachment; filename=`+torrent.Name+`.torrent"`)
	c.Writer.Write(buf.Bytes())
	return
}
