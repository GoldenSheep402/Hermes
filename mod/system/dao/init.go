package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Setting = &setting{}
)

func Init(db *gorm.DB, rds *redis.Client) error {
	if err := Setting.Init(db, rds); err != nil {
		return err
	}
	return nil
}
