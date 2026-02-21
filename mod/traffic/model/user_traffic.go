package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// UserTraffic stores aggregate traffic statistics for a user.
type UserTraffic struct {
	stdao.Model
	UserID        string `gorm:"type:char(26);uniqueIndex;not null" json:"user_id"`
	RealUpload    int64  `gorm:"not null;default:0" json:"real_upload"`
	RealDownload  int64  `gorm:"not null;default:0" json:"real_download"`
	BonusUpload   int64  `gorm:"not null;default:0" json:"bonus_upload"`
	BonusDownload int64  `gorm:"not null;default:0" json:"bonus_download"`
}
