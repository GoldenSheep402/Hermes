package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"time"
)

type Torrent struct {
	stdao.Model
	CategoryID   string `json:"category_id,omitempty"`
	InfoHash     string `json:"info_hash"`
	CreatorID    string `json:"creator_id"`
	IsSingleFile bool   `json:"is_single_file"`

	// The following fields are from the torrent file
	Announce     string     `json:"announce"`
	CreatedBy    *string    `json:"created_by,omitempty"`
	CreationDate *time.Time `json:"creation_date,omitempty"`
	Comment      *string    `json:"comment,omitempty"`
	Name         string     `json:"name"`
	NameUTF8     string     `json:"name_utf_8"`
	Length       *uint64    `json:"length,omitempty"`
	Md5sum       *string    `json:"md5sum,omitempty"`
	Pieces       []byte     `json:"pieces"`
	PieceLength  uint64     `json:"piece_length"`
	Private      *bool      `json:"private,omitempty"`
	Source       *string    `json:"source,omitempty"`
}

type FileInfo struct {
	Length   uint64   `bencode:"length"`
	Path     []string `bencode:"path"`
	PathUTF8 []string `bencode:"path.utf-8,omitempty"`
}

type BencodeInfo struct {
	Files       *[]FileInfo `bencode:"files,omitempty"`
	Name        string      `bencode:"name"`
	NameUTF8    *string     `bencode:"name.utf-8,omitempty"`
	Length      *uint64     `bencode:"length,omitempty"`
	Md5sum      *string     `bencode:"md5sum,omitempty"`
	Pieces      string      `bencode:"pieces"`
	PieceLength uint64      `bencode:"piece length"`
	Private     *int        `bencode:"private,omitempty"`
	Source      *string     `bencode:"source,omitempty"`
}

type BencodeTorrent struct {
	Announce     string      `bencode:"announce"`
	AnnounceList *[][]string `bencode:"announce-list,omitempty"`
	CreatedBy    *string     `bencode:"created by,omitempty"`
	CreatedAt    *int        `bencode:"creation date,omitempty"`
	Comment      *string     `bencode:"comment,omitempty"`
	Info         BencodeInfo `bencode:"info"`
}
