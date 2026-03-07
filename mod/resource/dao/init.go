package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Resource           = &resource{}
	ResourceMeta       = &resourceMeta{}
	ResourceScreenshot = &resourceScreenshot{}
	ResourceThank      = &resourceThank{}
	Tag                = &tag{}
	ResourceTag        = &resourceTag{}
)

func Init(db *gorm.DB, rdb *redis.Client) error {
	if err := Resource.Std.Init(db); err != nil {
		return err
	}
	if err := ResourceMeta.Std.Init(db); err != nil {
		return err
	}
	if err := ResourceScreenshot.Std.Init(db); err != nil {
		return err
	}
	if err := ResourceThank.Std.Init(db); err != nil {
		return err
	}
	if err := Tag.Std.Init(db); err != nil {
		return err
	}
	if err := ResourceTag.Std.Init(db); err != nil {
		return err
	}
	return nil
}
