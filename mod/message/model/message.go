package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// Message represents a private message or system notification.
type Message struct {
	stdao.Model
	SenderID   string `gorm:"type:char(26);index;not null" json:"sender_id"`
	ReceiverID string `gorm:"type:char(26);index;not null" json:"receiver_id"`
	Subject    string `gorm:"size:256;not null" json:"subject"`
	Body       string `gorm:"type:text;not null" json:"body"`
	IsRead     bool   `gorm:"not null;default:false" json:"is_read"`
	Type       int    `gorm:"not null;default:0" json:"type"` // 0=private, 1=system
}

// Message type constants.
const (
	MessageTypePrivate = 0
	MessageTypeSystem  = 1
)
