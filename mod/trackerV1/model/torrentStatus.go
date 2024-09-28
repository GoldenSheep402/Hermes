package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type TorrentStatus struct {
	stdao.Model
	TorrentID     string `gorm:"column:torrent_id; type:varchar(255); not null; index:torrent_id"`
	UploadCount   int    `gorm:"column:upload_count; type:int; not null; default:0"`
	UploadSum     int64  `gorm:"column:upload_sum; type:int; not null; default:0"`
	DownloadCount int    `gorm:"column:download_count; type:int; not null; default:0"`
	DownloadSum   int64  `gorm:"column:download_sum; type:int; not null; default:0"`
	SeedingCount  int    `gorm:"column:seeding_count; type:int; not null; default:0"`
}
