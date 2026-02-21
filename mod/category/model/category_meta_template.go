package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// CategoryMetaTemplate defines the schema of metadata fields for a category.
type CategoryMetaTemplate struct {
	stdao.Model
	CategoryID   string `gorm:"type:char(26);index;not null" json:"category_id"`
	Key          string `gorm:"size:64;not null" json:"key"`
	Label        string `gorm:"size:128;not null" json:"label"`
	Type         string `gorm:"size:32;not null" json:"type"`
	Required     bool   `gorm:"not null;default:false" json:"required"`
	Options      string `gorm:"type:text" json:"options"`
	SortOrder    int    `gorm:"not null;default:0" json:"sort_order"`
	DefaultValue string `gorm:"size:512" json:"default_value"`
}
