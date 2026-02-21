package model

import (
	"time"

	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// Snatch records a user's completion (snatch) of a torrent download.
type Snatch struct {
	stdao.Model
	TorrentID  string     `gorm:"type:char(26);index;not null" json:"torrent_id"`
	UserID     string     `gorm:"type:char(26);index;not null" json:"user_id"`
	Uploaded   int64      `gorm:"not null;default:0" json:"uploaded"`
	Downloaded int64      `gorm:"not null;default:0" json:"downloaded"`
	SeedTime   int64      `gorm:"not null;default:0" json:"seed_time"`
	IsActive   bool       `gorm:"not null;default:false" json:"is_active"`
	FinishedAt *time.Time `json:"finished_at"`
	LastAction time.Time  `gorm:"not null" json:"last_action"`
}
