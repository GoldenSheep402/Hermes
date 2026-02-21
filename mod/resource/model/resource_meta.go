package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// ResourceMeta stores metadata in EAV (Entity-Attribute-Value) format.
// The keys are defined by CategoryMetaTemplate per category.
type ResourceMeta struct {
	stdao.Model
	ResourceID string `gorm:"type:char(26);index;not null" json:"resource_id"`
	Key        string `gorm:"size:128;not null" json:"key"`
	Value      string `gorm:"type:text;not null" json:"value"`
}
