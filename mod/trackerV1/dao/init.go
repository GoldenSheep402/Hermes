package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Peer      = &peer{}
	TrackerV1 = &trackerV1{}
	Sum       = &sum{}
	SingleSum = &singleSum{}
)

func Init(db *gorm.DB, rds *redis.Client) error {
	err := Peer.Init(db, rds)
	if err != nil {
		return err
	}

	err = TrackerV1.Init(rds)
	if err != nil {
		return err
	}

	err = Sum.Init(db, rds)
	if err != nil {
		return err
	}

	err = SingleSum.Init(db, rds)
	if err != nil {
		return err
	}
	return nil
}
