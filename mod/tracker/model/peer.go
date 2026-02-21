package model

import (
	"time"

	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

// Peer represents an active peer in the swarm for a torrent.
// Hot data is stored in Redis; this table serves as persistent backup.
type Peer struct {
	stdao.Model
	TorrentID  string    `gorm:"type:char(26);index;not null" json:"torrent_id"`
	UserID     string    `gorm:"type:char(26);index;not null" json:"user_id"`
	PeerID     string    `gorm:"size:40;not null" json:"peer_id"`
	IP         string    `gorm:"size:45;not null" json:"ip"`
	Port       int       `gorm:"not null" json:"port"`
	Uploaded   int64     `gorm:"not null;default:0" json:"uploaded"`
	Downloaded int64     `gorm:"not null;default:0" json:"downloaded"`
	Left       int64     `gorm:"not null;default:0" json:"left"`
	Agent      string    `gorm:"size:128" json:"agent"`
	IsSeeder   bool      `gorm:"not null;default:false" json:"is_seeder"`
	StartedAt  time.Time `gorm:"not null" json:"started_at"`
	LastAction time.Time `gorm:"not null" json:"last_action"`
}
