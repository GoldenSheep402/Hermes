package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// Setting stores site-wide configuration as key-value pairs.
type Setting struct {
	stdao.Model
	Key   string `gorm:"uniqueIndex;size:128;not null" json:"key"`
	Value string `gorm:"type:text;not null" json:"value"`
	Type  string `gorm:"size:32;not null;default:'string'" json:"type"` // string, int, bool, json
	Desc  string `gorm:"size:256" json:"desc"`
}
