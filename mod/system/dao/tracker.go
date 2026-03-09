package dao

import (
	"context"
	"errors"
	"github.com/GoldenSheep402/Hermes/mod/system/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type tracker struct {
	stdao.Std[*model.Tracker]
	rds *redis.Client
}

func (it *tracker) Init(db *gorm.DB, rds *redis.Client) error {
	it.rds = rds
	return it.Std.Init(db)
}

func (it *tracker) GetTrackers(ctx context.Context) ([]model.Tracker, error) {
	key := "Trackers"

	if it.rds == nil {
		var trackers []model.Tracker
		if it.Std.DB() == nil {
			return nil, errors.New("tracker dao not initialized")
		}
		if err := it.Std.DB().WithContext(ctx).Find(&trackers).Error; err != nil {
			return nil, err
		}
		return trackers, nil
	}

	result, err := it.rds.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		var trackers []model.Tracker
		if err := it.Std.DB().WithContext(ctx).Find(&trackers).Error; err != nil {
			return nil, err
		}

		for _, tracker := range trackers {
			err := it.rds.HSet(ctx, key, tracker.Address, tracker.Enable).Err()
			if err != nil {
				return nil, err
			}
		}

		it.rds.Expire(ctx, key, 24*time.Hour)

		return trackers, nil
	}

	var trackers []model.Tracker
	for addr, enable := range result {
		isEnabled, err := strconv.ParseBool(enable)
		if err != nil {
			return nil, err
		}
		trackers = append(trackers, model.Tracker{
			Address: addr,
			Enable:  isEnabled,
		})
	}

	return trackers, nil
}

func (it *tracker) ClearTrackers(ctx context.Context) error {
	key := "Trackers"
	if it.rds == nil {
		return nil
	}
	return it.rds.Del(ctx, key).Err()
}
