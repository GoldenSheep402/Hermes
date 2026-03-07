package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Category             = &category{}
	CategoryMetaTemplate = &categoryMetaTemplate{}
)

func Init(db *gorm.DB, rdb *redis.Client) error {
	if err := Category.Std.Init(db); err != nil {
		return err
	}
	if err := CategoryMetaTemplate.Std.Init(db); err != nil {
		return err
	}
	return nil
}
