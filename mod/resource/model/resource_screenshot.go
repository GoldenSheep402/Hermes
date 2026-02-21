package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// ResourceScreenshot stores screenshot/preview images for a resource.
type ResourceScreenshot struct {
	stdao.Model
	ResourceID string `gorm:"type:char(26);index;not null" json:"resource_id"`
	URL        string `gorm:"size:1024;not null" json:"url"`
	SortOrder  int    `gorm:"not null;default:0" json:"sort_order"`
}
