package model

import "github.com/GoldenSheep402/Hermes/pkg/stdao"

type Sum struct {
	stdao.Model
	UID          string `gorm:"column:uid;type:varchar(255);not null" json:"uid"`
	RealUpload   int64  `gorm:"column:real_upload;type:bigint;not null" json:"real_upload"`
	RealDownload int64  `gorm:"column:real_download;type:bigint;not null" json:"real_download"`
	AddUpload    int64  `gorm:"column:add_upload;type:bigint;not null" json:"add_upload"`
	AddDownload  int64  `gorm:"column:add_download;type:bigint;not null" json:"add_download"`
}
