package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// Torrent stores the technical metadata of a .torrent file.
// Business-level information (title, description, category) lives in the Resource model.
type Torrent struct {
	stdao.Model
	InfoHash     string `gorm:"size:64;uniqueIndex;not null" json:"info_hash"`
	UploaderID   string `gorm:"type:char(26);index;not null" json:"uploader_id"`
	Name         string `gorm:"size:512;not null" json:"name"`
	Size         int64  `gorm:"not null;default:0" json:"size"`
	PieceLength  int64  `gorm:"not null" json:"piece_length"`
	PieceCount   int    `gorm:"not null" json:"piece_count"`
	IsSingleFile bool   `gorm:"not null;default:false" json:"is_single_file"`
	FileCount    int    `gorm:"not null;default:1" json:"file_count"`
	Comment      string `gorm:"size:512" json:"comment"`
	SeedCount    int    `gorm:"not null;default:0" json:"seed_count"`
	LeechCount   int    `gorm:"not null;default:0" json:"leech_count"`
	SnatchCount  int    `gorm:"not null;default:0" json:"snatch_count"`
	IsActive     bool   `gorm:"not null;default:true" json:"is_active"`
}

// BencodeInfo is used for parsing .torrent file info dict.
type BencodeInfo struct {
	Files       *[]BencodeFile `bencode:"files,omitempty"`
	Name        string         `bencode:"name"`
	NameUTF8    *string        `bencode:"name.utf-8,omitempty"`
	Length      *int64         `bencode:"length,omitempty"`
	Md5sum      *string        `bencode:"md5sum,omitempty"`
	Pieces      string         `bencode:"pieces"`
	PieceLength int64          `bencode:"piece length"`
	Private     *int           `bencode:"private,omitempty"`
	Source      *string        `bencode:"source,omitempty"`
}

// BencodeFile represents a single file entry inside a multi-file torrent.
type BencodeFile struct {
	Length   int64    `bencode:"length"`
	Path     []string `bencode:"path"`
	PathUTF8 []string `bencode:"path.utf-8,omitempty"`
}

// BencodeTorrent is the top-level structure of a .torrent file.
type BencodeTorrent struct {
	Announce     string      `bencode:"announce"`
	AnnounceList *[][]string `bencode:"announce-list,omitempty"`
	CreatedBy    *string     `bencode:"created by,omitempty"`
	CreatedAt    *int        `bencode:"creation date,omitempty"`
	Comment      *string     `bencode:"comment,omitempty"`
	Info         BencodeInfo `bencode:"info"`
}
