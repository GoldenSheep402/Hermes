package model

import (
	"time"

	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// Resource is the core publishing entity of the PT site.
// It binds a Torrent to a Category and carries all user-facing information.
type Resource struct {
	stdao.Model
	Title        string     `gorm:"size:512;not null;index" json:"title"`
	Subtitle     string     `gorm:"size:512" json:"subtitle"`
	Description  string     `gorm:"type:text" json:"description"`
	CategoryID   string     `gorm:"type:char(26);index;not null" json:"category_id"`
	TorrentID    string     `gorm:"type:char(26);uniqueIndex;not null" json:"torrent_id"`
	UploaderID   string     `gorm:"type:char(26);index;not null" json:"uploader_id"`
	Status       int        `gorm:"not null;default:0" json:"status"` // 0=pending, 1=approved, 2=dead, 3=banned
	IsSticky     bool       `gorm:"not null;default:false" json:"is_sticky"`
	IsFree       bool       `gorm:"not null;default:false" json:"is_free"`
	FreeUntil    *time.Time `json:"free_until"`
	DoubleUpload bool       `gorm:"not null;default:false" json:"double_upload"`
	DoubleUntil  *time.Time `json:"double_until"`
	ViewCount    int        `gorm:"not null;default:0" json:"view_count"`
	CommentCount int        `gorm:"not null;default:0" json:"comment_count"`
	ThankCount   int        `gorm:"not null;default:0" json:"thank_count"`
}

// Resource status constants.
const (
	ResourceStatusPending  = 0
	ResourceStatusApproved = 1
	ResourceStatusDead     = 2
	ResourceStatusBanned   = 3
)
