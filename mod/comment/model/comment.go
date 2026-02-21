package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// Comment represents a user comment on a resource.
type Comment struct {
	stdao.Model
	ResourceID string  `gorm:"type:char(26);index;not null" json:"resource_id"`
	AuthorID   string  `gorm:"type:char(26);index;not null" json:"author_id"`
	Body       string  `gorm:"type:text;not null" json:"body"`
	ParentID   *string `gorm:"type:char(26);index" json:"parent_id"` // for nested replies
}
