package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/category/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type metadata struct {
	stdao.Std[*model.Metadata]
}

func (m *metadata) Init(db *gorm.DB) error {
	return m.Std.Init(db)
}
