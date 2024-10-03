package dao

import "github.com/redis/go-redis/v9"

type code struct {
	rds *redis.Client
}

func (c *code) Init(rds *redis.Client) error {
	c.rds = rds
	return nil
}
