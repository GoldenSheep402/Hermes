package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// TorrentStats stores aggregate statistics for a torrent.
type TorrentStats struct {
	stdao.Model
	TorrentID     string `gorm:"type:char(26);uniqueIndex;not null" json:"torrent_id"`
	SeedCount     int    `gorm:"not null;default:0" json:"seed_count"`
	LeechCount    int    `gorm:"not null;default:0" json:"leech_count"`
	SnatchCount   int    `gorm:"not null;default:0" json:"snatch_count"`
	TotalUpload   int64  `gorm:"not null;default:0" json:"total_upload"`
	TotalDownload int64  `gorm:"not null;default:0" json:"total_download"`
}
