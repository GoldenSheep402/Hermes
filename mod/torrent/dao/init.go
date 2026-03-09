package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Torrent      = &torrent{}
	TorrentFile  = &torrentFile{}
	TorrentBlob  = &torrentBlob{}
	TorrentPiece = &torrentPiece{}
)

func Init(db *gorm.DB, rds *redis.Client) error {
	if err := Torrent.Init(db); err != nil {
		return err
	}
	if err := TorrentFile.Init(db); err != nil {
		return err
	}
	if err := TorrentBlob.Init(db); err != nil {
		return err
	}
	if err := TorrentPiece.Init(db); err != nil {
		return err
	}
	return nil
}
