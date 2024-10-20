package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type groupMetadata struct {
	stdao.Std[*model.GroupMetadata]
}

func (g *groupMetadata) Init(db *gorm.DB) error {
	return g.Std.Init(db)
}
