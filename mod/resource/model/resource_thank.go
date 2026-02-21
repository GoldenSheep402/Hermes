package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// ResourceThank records a user "thanking" a resource (similar to upvote/like).
type ResourceThank struct {
	stdao.Model
	ResourceID string `gorm:"type:char(26);index;not null" json:"resource_id"`
	UserID     string `gorm:"type:char(26);index;not null" json:"user_id"`
}
