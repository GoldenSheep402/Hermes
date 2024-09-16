package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/torrent/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type file struct {
	stdao.Std[*model.File]
}

func (f *file) Init(db *gorm.DB) error {
	return f.Std.Init(db)
}
