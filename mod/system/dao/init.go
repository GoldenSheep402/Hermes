package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Setting = &setting{}
	Subnet  = &subnet{}
	Tracker = &tracker{}
)

func Init(db *gorm.DB, rds *redis.Client) error {
	if err := Setting.Init(db, rds); err != nil {
		return err
	}
	if err := Subnet.Init(db); err != nil {
		return err
	}
	if err := Tracker.Init(db, rds); err != nil {
		return err
	}
	return nil
}
