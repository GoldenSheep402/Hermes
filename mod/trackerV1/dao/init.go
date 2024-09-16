package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Peer = &peer{}
)

func Init(db *gorm.DB, rds *redis.Client) error {
	err := Peer.Init(db, rds)
	if err != nil {
		return err
	}
	return nil
}
