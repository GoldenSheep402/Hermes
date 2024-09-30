package model

import "github.com/GoldenSheep402/Hermes/pkg/stdao"

type SingleSum struct {
	stdao.Model
	UID       string `gorm:"column:uid;type:varchar(255);not null" json:"uid"`
	TorrentID string `gorm:"column:torrent_id;type:varchar(255);not null" json:"torrent_id"`
	Upload    int64  `gorm:"column:upload;type:bigint;not null" json:"upload"`
	Download  int64  `gorm:"column:download;type:bigint;not null" json:"download"`
	IsFinish  bool   `gorm:"column:is_finish;type:bool;not null" json:"is_finish"`
}
