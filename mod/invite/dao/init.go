package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	InviteCode = &inviteCode{}
)

func Init(db *gorm.DB, rds *redis.Client) error {
	if err := InviteCode.Init(db); err != nil {
		return err
	}
	return nil
}
