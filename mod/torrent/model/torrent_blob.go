package model

import "github.com/GoldenSheep402/Hermes/pkg/stdao"

// TorrentBlob stores raw .torrent bytes separately from the main torrent table
// to keep common torrent queries light.
type TorrentBlob struct {
	stdao.Model
	TorrentID string `gorm:"type:char(26);uniqueIndex;not null" json:"torrent_id"`
	RawData   []byte `gorm:"not null" json:"-"`
}
