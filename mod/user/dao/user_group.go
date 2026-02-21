package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type userGroup struct {
	stdao.Std[*model.UserGroup]
}

func (g *userGroup) Init(db *gorm.DB) error {
	return g.Std.Init(db)
}
