package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Message = &message{}
)

func Init(db *gorm.DB, rds *redis.Client) error {
	if err := Message.Init(db); err != nil {
		return err
	}
	return nil
}
