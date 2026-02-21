package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Peer   = &peer{}
	Snatch = &snatch{}
)

func Init(db *gorm.DB, rds *redis.Client) error {
	if err := Peer.Init(db, rds); err != nil {
		return err
	}
	if err := Snatch.Init(db); err != nil {
		return err
	}
	return nil
}
