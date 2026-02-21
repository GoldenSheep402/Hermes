package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/system/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type setting struct {
	stdao.Std[*model.Setting]
	rds *redis.Client
}

func (s *setting) Init(db *gorm.DB, rds *redis.Client) error {
	s.rds = rds
	return s.Std.Init(db)
}
