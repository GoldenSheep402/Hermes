package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Category = &category{}
	Metadata = &metadata{}
)

func Init(DB *gorm.DB, rds *redis.Client) error {
	err := Category.Init(DB)
	if err != nil {
		return err
	}

	err = Metadata.Init(DB)
	if err != nil {
		return err
	}

	return nil
}
