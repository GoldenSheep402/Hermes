package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type TorrentMetadata struct {
	stdao.Model
	TorrentID  string `json:"torrent_id"`
	MetadataID string `json:"metadata_id"`
	Value      string `json:"value"`
}

func (t *TorrentMetadata) BeforeCreate(_ *gorm.DB) error {
	t.ID = ulid.Make().String()
	return nil
}
