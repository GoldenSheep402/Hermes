package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoldenSheep402/Hermes/mod/auth/model/codeValues"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

var (
	Code = &code{}
)

func Init(db *gorm.DB, rds *redis.Client) error {
	if err := Code.Init(rds); err != nil {
		return err
	}

	return nil
}

func (c *code) SetCodeWithEmail(ctx context.Context, email string, code string) error {
	return c.rds.Set(ctx, email, code, time.Minute*5).Err()
}

func (c *code) CheckCodeWithAttempts(ctx context.Context, email string, code string) (int, error) {
	attemptsKey := email + "_attempts"

	attempts, err := c.rds.Incr(ctx, attemptsKey).Result()
	if err != nil {
		return codeValues.TooManyAttempts, err
	}

	if attempts == 1 {
		c.rds.Expire(ctx, attemptsKey, time.Minute*30)
	}

	if attempts > 5 {
		return codeValues.TooManyAttempts, fmt.Errorf("too many attempts")
	}

	storedCode, err := c.rds.Get(ctx, email).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return codeValues.Wrong, fmt.Errorf("code not found")
		}
		return codeValues.Wrong, err
	}

	if code != storedCode {
		return codeValues.Wrong, fmt.Errorf("wrong code")
	}

	c.rds.Del(ctx, email)
	c.rds.Del(ctx, attemptsKey)

	return codeValues.Right, nil
}
