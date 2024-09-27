package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/model/trackerV1Values"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/redis/go-redis/v9"
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

	} else if err != nil {
		userDao.User.DB()
	} else {
		if val == trackerV1Values.OK {
			return trackerV1Values.OK, nil
			if err := t.rds.Expire(ctx, key, trackerV1Values.TTL).Err(); err != nil {

			}
		}
		if err := t.rds.Set(ctx, key, newValue, ttl).Err(); err != nil {
			fmt.Println("error updating key value:", err)
		} else {
			fmt.Println("key value updated successfully and TTL extended")
		}
	}

}
