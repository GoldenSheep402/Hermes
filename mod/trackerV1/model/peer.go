package model

import "time"

type Peer struct {
	PeerID   string
	IP       string
	Port     int
	LastSeen time.Time
	Status   int
}
