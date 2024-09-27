package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type sum struct {
	stdao.Std[model.Sum]
}

func (s *sum) Init(db *gorm.DB, rds *redis.Client) error {
	err := s.Std.Init(db)
	if err != nil {
		return err
	}

	return nil
}
