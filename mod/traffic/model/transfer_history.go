package model

import (
	"time"

	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// TransferHistory tracks per-torrent transfer statistics for each user.
type TransferHistory struct {
	stdao.Model
	UserID     string    `gorm:"type:char(26);index;not null" json:"user_id"`
	TorrentID  string    `gorm:"type:char(26);index;not null" json:"torrent_id"`
	Uploaded   int64     `gorm:"not null;default:0" json:"uploaded"`
	Downloaded int64     `gorm:"not null;default:0" json:"downloaded"`
	SeedTime   int64     `gorm:"not null;default:0" json:"seed_time"`
	IsFinished bool      `gorm:"not null;default:false" json:"is_finished"`
	IsActive   bool      `gorm:"not null;default:false" json:"is_active"`
	LastAction time.Time `gorm:"not null" json:"last_action"`
}
