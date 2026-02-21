package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// Category represents a content category (e.g., Movies, Books, Software).
type Category struct {
	stdao.Model
	Name        string  `gorm:"uniqueIndex;size:64;not null" json:"name"`
	Slug        string  `gorm:"uniqueIndex;size:64;not null" json:"slug"`
	Description string  `gorm:"size:512" json:"description"`
	ParentID    *string `gorm:"type:char(26);index" json:"parent_id"`
	Icon        string  `gorm:"size:256" json:"icon"`
	SortOrder   int     `gorm:"not null;default:0" json:"sort_order"`
	IsEnabled   bool    `gorm:"not null;default:true" json:"is_enabled"`
}
