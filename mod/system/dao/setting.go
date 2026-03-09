package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/GoldenSheep402/Hermes/mod/system/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	SettingKeyTrackerList           = "tracker.list"
	SettingKeyTrackerFlushInterval  = "tracker.flush_interval"
	SettingKeyTrackerFlushBatchSize = "tracker.flush_batch_size"
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

func (s *setting) GetIntByKey(ctx context.Context, key string) (int, error) {
	item, err := s.GetByKey(ctx, key)
	if err != nil {
		return 0, err
	}
	val, err := strconv.Atoi(strings.TrimSpace(item.Value))
	if err != nil {
		return 0, fmt.Errorf("setting %s invalid int value %q: %w", key, item.Value, err)
	}
	return val, nil
}

func (s *setting) UpdateOrCreate(ctx context.Context, item *model.Setting) error {
	return s.GetTxFromCtx(ctx).WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "key"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"value":      item.Value,
				"type":       item.Type,
				"desc":       item.Desc,
				"updated_at": time.Now(),
			}),
		}).
		Create(item).Error
}
