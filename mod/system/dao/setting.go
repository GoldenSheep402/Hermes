package dao

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/system/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	SettingKeyTrackerList = "tracker.list"
)

type setting struct {
	stdao.Std[*model.Setting]
	rds *redis.Client
}

func (s *setting) Init(db *gorm.DB, rds *redis.Client) error {
	s.rds = rds
	return s.Std.Init(db)
}

func (s *setting) GetByKey(ctx context.Context, key string) (*model.Setting, error) {
	var val model.Setting
	if err := s.GetTxFromCtx(ctx).WithContext(ctx).
		Where(clause.Eq{Column: clause.Column{Name: "key"}, Value: key}).
		First(&val).Error; err != nil {
		return nil, err
	}
	return &val, nil
}

func (s *setting) DeleteByKey(ctx context.Context, key string) error {
	return s.GetTxFromCtx(ctx).WithContext(ctx).
		Where(clause.Eq{Column: clause.Column{Name: "key"}, Value: key}).
		Delete(&model.Setting{}).Error
}
