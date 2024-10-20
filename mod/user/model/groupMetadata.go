package model

import "github.com/GoldenSheep402/Hermes/pkg/stdao"

type GroupMetadata struct {
	stdao.Model
	GID         string `gorm:"type:char(26);not null;column:gid" json:"gid"`
	Key         string `gorm:"type:varchar(255);not null;column:key" json:"key"`
	Value       string `gorm:"type:varchar(255);not null;column:value" json:"value"`
	Type        string `gorm:"type:varchar(255);not null;column:type" json:"type"`
	Description string `gorm:"type:varchar(255);not null;column:description" json:"description"`
	Order       int    `gorm:"type:int;not null;column:order" json:"order"`
}
