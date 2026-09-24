package dao

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

const (
	passkeyCacheTTL    = 10 * time.Minute
	passkeyNegativeTTL = 60 * time.Second
	passkeyCacheKeyPref = "tracker:passkey:"
)

func passkeyCacheKey(passkey string) string {
	return passkeyCacheKeyPref + passkey
}

// GetByPasskeyCached resolves a user by passkey via Redis, falling back to Postgres.
// The cached entry only carries ID / IsEnabled / Passkey — enough for announce auth.
func (u *user) GetByPasskeyCached(ctx context.Context, passkey string) (*model.User, error) {
	if passkey == "" {
		return nil, gorm.ErrRecordNotFound
	}

	if u.rds != nil {
		key := passkeyCacheKey(passkey)
		values, err := u.rds.HGetAll(ctx, key).Result()
		if err == nil && len(values) > 0 {
			if values["missing"] == "1" {
				return nil, gorm.ErrRecordNotFound
			}
			userID := values["user_id"]
			if userID == "" {
				return nil, gorm.ErrRecordNotFound
			}
			return &model.User{
				Model:     stdao.Model{ID: userID},
				Passkey:   passkey,
				IsEnabled: values["is_enabled"] == "1",
			}, nil
		}
	}

	user, err := u.GetByPasskey(ctx, passkey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.cachePasskeyMiss(ctx, passkey)
		}
		return nil, err
	}
	u.cachePasskeyHit(ctx, passkey, user.ID, user.IsEnabled)
	return user, nil
}

func (u *user) cachePasskeyHit(ctx context.Context, passkey, userID string, enabled bool) {
	if u.rds == nil || passkey == "" || userID == "" {
		return
	}
	enabledVal := "0"
	if enabled {
		enabledVal = "1"
	}
	key := passkeyCacheKey(passkey)
	_ = u.rds.HSet(ctx, key, map[string]interface{}{
		"user_id":    userID,
		"is_enabled": enabledVal,
		"missing":    "0",
	}).Err()
	_ = u.rds.Expire(ctx, key, passkeyCacheTTL).Err()
}

func (u *user) cachePasskeyMiss(ctx context.Context, passkey string) {
	if u.rds == nil || passkey == "" {
		return
	}
	key := passkeyCacheKey(passkey)
	_ = u.rds.HSet(ctx, key, map[string]interface{}{
		"missing": "1",
	}).Err()
	_ = u.rds.Expire(ctx, key, passkeyNegativeTTL).Err()
}

// InvalidatePasskeyCache removes a cached passkey entry.
func (u *user) InvalidatePasskeyCache(ctx context.Context, passkey string) {
	if u.rds == nil || passkey == "" {
		return
	}
	_ = u.rds.Del(ctx, passkeyCacheKey(passkey)).Err()
}

// SetRedis sets the Redis client (mainly for tests).
func (u *user) SetRedis(rds *redis.Client) {
	u.rds = rds
}
