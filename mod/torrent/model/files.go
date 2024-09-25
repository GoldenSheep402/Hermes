package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type File struct {
	stdao.Model
	TorrentID string `json:"torrent_id"`
	Length    uint64 `json:"length"`
	Path      string `json:"path"`
	PathUTF8  string `json:"path_utf_8"`
}

func (f *File) BeforeCreate(_ *gorm.DB) error {
	f.ID = ulid.Make().String()
	return nil
}
