package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type groupMembership struct {
	stdao.Std[*model.GroupMembership]
}

func (g *groupMembership) Init(db *gorm.DB) error {
	return g.Std.Init(db)
}
