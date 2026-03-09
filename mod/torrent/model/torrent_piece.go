package model

import "github.com/GoldenSheep402/Hermes/pkg/stdao"

// TorrentPiece stores piece-level SHA1 hashes for each torrent.
type TorrentPiece struct {
	stdao.Model
	TorrentID  string `gorm:"type:char(26);uniqueIndex:uidx_torrent_piece,priority:1;index;not null" json:"torrent_id"`
	PieceIndex int    `gorm:"uniqueIndex:uidx_torrent_piece,priority:2;not null" json:"piece_index"`
	PieceSHA1  string `gorm:"size:40;not null" json:"piece_sha1"`
}
