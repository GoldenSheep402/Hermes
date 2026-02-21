package model

import (
	"time"

	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type User struct {
	stdao.Model
	Username    string     `gorm:"uniqueIndex;size:32;not null" json:"username"`
	Email       string     `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Password    string     `gorm:"not null" json:"-"`
	Salt        string     `gorm:"not null" json:"-"`
	Passkey     string     `gorm:"uniqueIndex;size:64;not null" json:"passkey"`
	Avatar      string     `gorm:"size:512" json:"avatar"`
	IsAdmin     bool       `gorm:"not null;default:false" json:"is_admin"`
	IsEnabled   bool       `gorm:"not null;default:true" json:"is_enabled"`
	GroupID     string     `gorm:"type:char(26);index" json:"group_id"`
	BonusPoints int64      `gorm:"not null;default:0" json:"bonus_points"`
	Uploaded    int64      `gorm:"not null;default:0" json:"uploaded"`
	Downloaded  int64      `gorm:"not null;default:0" json:"downloaded"`
	SeedTime    int64      `gorm:"not null;default:0" json:"seed_time"`
	InviteCount int        `gorm:"not null;default:0" json:"invite_count"`
	LastLogin   *time.Time `json:"last_login"`
}
