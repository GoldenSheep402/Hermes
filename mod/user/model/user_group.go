package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type UserGroup struct {
	stdao.Model
	Name            string  `gorm:"uniqueIndex;size:64;not null" json:"name"`
	Description     string  `gorm:"size:512" json:"description"`
	Level           int     `gorm:"not null;default:0" json:"level"`
	MinUpload       int64   `gorm:"not null;default:0" json:"min_upload"`
	MinRatio        float64 `gorm:"not null;default:0" json:"min_ratio"`
	MinSeedTime     int64   `gorm:"not null;default:0" json:"min_seed_time"`
	MaxDownloads    int     `gorm:"not null;default:-1" json:"max_downloads"`
	CanUpload       bool    `gorm:"not null;default:true" json:"can_upload"`
	CanInvite       bool    `gorm:"not null;default:false" json:"can_invite"`
	IsImmuneToRatio bool    `gorm:"not null;default:false" json:"is_immune_to_ratio"`
	Color           string  `gorm:"size:16" json:"color"`
	Icon            string  `gorm:"size:512" json:"icon"`
}
