package dao

import (
	"context"
	"github.com/GoldenSheep402/Hermes/mod/system/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type innetTracker struct {
	stdao.Std[*model.InnetTracker]
	rds *redis.Client
}

func (it *innetTracker) Init(db *gorm.DB, rds *redis.Client) error {
	it.rds = rds
	return it.Std.Init(db)
}

func (it *innetTracker) GetTrackers(ctx context.Context) ([]model.InnetTracker, error) {
	key := "Trackers"

	result, err := it.rds.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		var trackers []model.InnetTracker
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

	var trackers []model.InnetTracker
	for addr, enable := range result {
		isEnabled, err := strconv.ParseBool(enable)
		if err != nil {
			return nil, err
		}
		trackers = append(trackers, model.InnetTracker{
			Address: addr,
			Enable:  isEnabled,
		})
	}

	return trackers, nil
}

func (it *innetTracker) ClearTrackers(ctx context.Context) error {
	key := "Trackers"
	return it.rds.Del(ctx, key).Err()
}
