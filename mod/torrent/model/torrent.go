package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"time"
)

type Torrent struct {
	stdao.Model
	CategoryID string `json:"category_id,omitempty"`
	InfoHash   string `json:"info_hash"`
	CreatorID  string `json:"creator_id"`

	// The following fields are from the torrent file
	Announce     string     `json:"announce"`
	CreatedBy    *string    `json:"created_by,omitempty"`
	CreationDate *time.Time `json:"creation_date,omitempty"`
	Name         string     `json:"name"`
	Length       *uint64    `json:"length,omitempty"`
	Md5sum       *string    `json:"md5sum,omitempty"`
	Pieces       []byte     `json:"pieces"`
	PieceLength  uint64     `json:"piece_length"`
	Private      *bool      `json:"private,omitempty"`
	Source       *string    `json:"source,omitempty"`
}
