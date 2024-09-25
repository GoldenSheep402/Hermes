package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type TorrentMetadata struct {
	stdao.Model
	TorrentID  string `json:"torrent_id"`
	CategoryID string `json:"category_id"`
	MetadataID string `json:"metadata_id"`
	Value      string `json:"value"`
}
