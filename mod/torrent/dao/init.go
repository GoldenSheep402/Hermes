package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Torrent         = &torrent{}
	TorrentMetadata = &torrentMetadata{}
	File            = &file{}
)

func Init(db *gorm.DB, rds *redis.Client) error {
	err := Torrent.Init(db)
	if err != nil {
		return err
	}

	err = File.Init(db)
	if err != nil {
		return err
	}

	err = TorrentMetadata.Init(db)
	if err != nil {
		return err
	}

	return nil
}
