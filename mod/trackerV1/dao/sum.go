package dao

import (
	"context"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/model"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/model/trackerV1Values"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
)

type sum struct {
	stdao.Std[model.Sum]
	rds *redis.Client
}

func (s *sum) Init(db *gorm.DB, rds *redis.Client) error {
	err := s.Std.Init(db)
	if err != nil {
		return err
	}

	s.rds = rds

	return nil
}

func (s *sum) Sync() error {
	ctx := context.Background()

	// hash map
	torrentSeedingKeyPrefix := "TorrentSeeding:"
	torrentSeedingKeys, err := s.rds.Keys(ctx, torrentSeedingKeyPrefix+"*").Result()
	if err != nil {
		return err
	}

	for _, key := range torrentSeedingKeys {
		//torrentID := key[len(torrentSeedingKeyPrefix):]
		dataTmp := make(map[string]map[string]string)
		fields, err := s.rds.HGetAll(ctx, key).Result()
		if err != nil {
			return err
		}

		dataTmp[key] = fields
		var downloading int64
		var seeding int64
		var complete int64
		for _, value := range fields {
			if value == strconv.Itoa(trackerV1Values.Downloading) {
				downloading++
			} else if value == strconv.Itoa(trackerV1Values.Seeding) {
				seeding++
			} else if value == strconv.Itoa(trackerV1Values.Finished) {
				complete++
			}
		}
	}

	return nil
}
