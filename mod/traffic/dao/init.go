package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	UserTraffic     = &userTraffic{}
	TransferHistory = &transferHistory{}
	TorrentStats    = &torrentStats{}
)

func Init(db *gorm.DB, rds *redis.Client) error {
	if err := UserTraffic.Init(db); err != nil {
		return err
	}
	if err := TransferHistory.Init(db); err != nil {
		return err
	}
	if err := TorrentStats.Init(db); err != nil {
		return err
	}
	return nil
}
