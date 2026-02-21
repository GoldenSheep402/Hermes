package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// BonusLog records all bonus point transactions for a user.
type BonusLog struct {
	stdao.Model
	UserID      string `gorm:"type:char(26);index;not null" json:"user_id"`
	Amount      int64  `gorm:"not null" json:"amount"` // positive = earned, negative = spent
	Reason      string `gorm:"size:32;not null" json:"reason"`
	Description string `gorm:"size:512" json:"description"`
	RelatedID   string `gorm:"type:char(26)" json:"related_id"`
}

// Bonus reason constants.
const (
	BonusReasonSeeding        = "seeding"
	BonusReasonExchangeUpload = "exchange_upload"
	BonusReasonGift           = "gift"
	BonusReasonSystem         = "system"
)
