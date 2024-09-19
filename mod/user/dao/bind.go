package dao

import (
	"context"
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/redis/go-redis/v9"
	"time"
)

type bind struct {
	stdao.Std[*model.Bind]
}

// func (b *bind) Init(db *gorm.DB) error {
//	b.Std = stdao.Create(&model.Bind{})
//	return b.Std.Init(db)
// }

var BindMark = &bindMark{prefix: "user:bind:mark"}

type bindMark struct {
	rds    *redis.Client
	prefix string
}

func (v *bindMark) Init(rds *redis.Client) {
	v.rds = rds
}

func (v *bindMark) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return v.rds.Set(ctx, v.prefix+":"+key, value, expiration)
}

func (v *bindMark) Get(ctx context.Context, key string) *redis.StringCmd {
	return v.rds.Get(ctx, v.prefix+":"+key)
}

func (v *bindMark) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return v.rds.Expire(ctx, v.prefix+":"+key, expiration)
}
