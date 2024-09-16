package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/torrent/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type torrentMetadata struct {
	stdao.Std[*model.TorrentMetadata]
}

func (t *torrentMetadata) Init(db *gorm.DB) error {
	return t.Std.Init(db)
}
