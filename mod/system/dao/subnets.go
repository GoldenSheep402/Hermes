package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/system/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type subnet struct {
	stdao.Std[*model.Subnet]
}

func (s *subnet) Init(db *gorm.DB) error {
	return s.Std.Init(db)
}
