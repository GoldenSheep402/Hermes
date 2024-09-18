package model

import "github.com/GoldenSheep402/Hermes/pkg/stdao"

type GroupMembership struct {
	stdao.Model
	UID string `gorm:"type:char(26);not null;column:uid" json:"uid"`
	GID string `gorm:"type:char(26);not null;column:gid" json:"gid"`
}
