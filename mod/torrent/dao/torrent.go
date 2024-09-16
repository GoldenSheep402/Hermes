package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/torrent/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type torrent struct {
	stdao.Std[*model.Torrent]
}

func (t *torrent) Init(db *gorm.DB) error {
	return t.Std.Init(db)
}
