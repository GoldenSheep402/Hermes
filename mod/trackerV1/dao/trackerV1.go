package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/model/trackerV1Values"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/redis/go-redis/v9"
	"time"
)

type trackerV1 struct {
	rds *redis.Client
}

func (t *trackerV1) Init(rds *redis.Client) error {
	t.rds = rds
	return nil
}

func (t *trackerV1) CheckKey(ctx context.Context, key string) (string, error) {
	val, err := t.rds.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		_, err := userDao.User.CheckKey(ctx, key)
		if err != nil {
			return "", err
		}

		// TODO: status check
		return t.setKeyWithTTL(ctx, key, trackerV1Values.OK, trackerV1Values.TTL)
	}

	if err != nil {
		return "", err
	}

	if val == trackerV1Values.OK {
		if err := t.rds.Expire(ctx, key, trackerV1Values.TTL).Err(); err != nil {
			return "", err
		}
		return trackerV1Values.OK, nil
	} else if val == trackerV1Values.Banned {
		return trackerV1Values.Banned, nil
	}

	return "", fmt.Errorf("unknown key status: %s", val)
}

func (t *trackerV1) setKeyWithTTL(ctx context.Context, key string, value string, ttl time.Duration) (string, error) {
	if err := t.rds.Set(ctx, key, value, ttl).Err(); err != nil {
		return "", err
	}
	return value, nil
}
