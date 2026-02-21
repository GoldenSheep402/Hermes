package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	User      = &user{}
	UserGroup = &userGroup{}
)

func Init(DB *gorm.DB, rds *redis.Client) error {
	if err := User.Init(DB, rds); err != nil {
		return err
	}

	if err := UserGroup.Init(DB); err != nil {
		return err
	}

	return nil
}
