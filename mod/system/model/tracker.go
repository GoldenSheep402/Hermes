package model

import "github.com/GoldenSheep402/Hermes/pkg/stdao"

type Tracker struct {
	stdao.Model
	Address string `json:"address"`
	Enable  bool   `json:"enable"`
}

// TableName keeps compatibility with existing schema created before renaming.
func (Tracker) TableName() string {
	return "innet_trackers"
}
