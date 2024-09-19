package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type groupMembershipMetadata struct {
	stdao.Std[*model.GroupMembershipMetadata]
}

func (g *groupMembershipMetadata) Init(db *gorm.DB) error {
	return g.Std.Init(db)
}
