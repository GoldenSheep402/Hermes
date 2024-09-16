package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/category/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type category struct {
	stdao.Std[*model.Category]
}

func (c *category) Init(db *gorm.DB) error {
	return c.Std.Init(db)
}
