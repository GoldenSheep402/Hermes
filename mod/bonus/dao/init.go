package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	BonusLog = &bonusLog{}
)

func Init(db *gorm.DB, rds *redis.Client) error {
	if err := BonusLog.Init(db); err != nil {
		return err
	}
	return nil
}
