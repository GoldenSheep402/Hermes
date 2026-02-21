package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/torrent/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type torrentFile struct {
	stdao.Std[*model.TorrentFile]
}

func (f *torrentFile) Init(db *gorm.DB) error {
	return f.Std.Init(db)
}
