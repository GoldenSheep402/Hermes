package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// Tag is a freeform label that can be attached to resources.
type Tag struct {
	stdao.Model
	Name string `gorm:"uniqueIndex;size:64;not null" json:"name"`
}

// ResourceTag is the many-to-many join table between Resource and Tag.
type ResourceTag struct {
	stdao.Model
	ResourceID string `gorm:"type:char(26);index;not null" json:"resource_id"`
	TagID      string `gorm:"type:char(26);index;not null" json:"tag_id"`
}
