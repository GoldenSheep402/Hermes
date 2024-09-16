package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type group struct {
	stdao.Std[*model.Group]
}

func (g *group) Init(db *gorm.DB) error {
	return g.Std.Init(db)
}
