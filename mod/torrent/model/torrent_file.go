package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// TorrentFile stores individual file entries within a torrent.
type TorrentFile struct {
	stdao.Model
	TorrentID string `gorm:"type:char(26);index;not null" json:"torrent_id"`
	Path      string `gorm:"size:1024;not null" json:"path"`
	Size      int64  `gorm:"not null" json:"size"`
}
