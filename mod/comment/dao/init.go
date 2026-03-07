package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Comment = &comment{}
)

func Init(db *gorm.DB, rds *redis.Client) error {
	if err := Comment.Init(db); err != nil {
		return err
	}
	return nil
}
