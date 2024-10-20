package model

import "github.com/GoldenSheep402/Hermes/pkg/stdao"

type GroupMembershipMetadata struct {
	stdao.Model
	MembershipID    string `gorm:"type:char(26);not null;column:membership_id" json:"membership_id"`
	GroupMetadataID string `gorm:"type:varchar(255);not null;column:group_metadata_id" json:"group_metadata_id"`
	Key             string `gorm:"type:varchar(255);not null;column:key" json:"key"`
	Value           string `gorm:"type:varchar(255);not null;column:value" json:"value"`
	Type            string `gorm:"type:varchar(255);not null;column:type" json:"type"`
	Description     string `gorm:"type:varchar(255);not null;column:description" json:"description"`
	Order           int    `gorm:"type:int;not null;column:order" json:"order"`
}
