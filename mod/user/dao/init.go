package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Bind                    = &bind{}
	User                    = &user{}
	Group                   = &group{}
	GroupMetadata           = &groupMetadata{}
	GroupMembership         = &groupMembership{}
	GroupMembershipMetadata = &groupMembershipMetadata{}
)

func Init(DB *gorm.DB, rds *redis.Client) error {
	err := Bind.Init(DB)
	if err != nil {
		return err
	}

	err = User.Init(DB)
	if err != nil {
		return err
	}

	err = Group.Init(DB)
	if err != nil {
		return err
	}

	err = GroupMetadata.Init(DB)
	if err != nil {
		return err
	}

	err = GroupMembership.Init(DB)
	if err != nil {
		return err
	}

	err = GroupMembershipMetadata.Init(DB)
	if err != nil {
		return err
	}

	return nil
}
