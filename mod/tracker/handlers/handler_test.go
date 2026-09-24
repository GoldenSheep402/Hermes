package handlers

import (
	"encoding/hex"
	"testing"

	"github.com/GoldenSheep402/Hermes/mod/tracker/model"
	"github.com/stretchr/testify/assert"
)

func TestBuildPeerList(t *testing.T) {
	excludePeerIDHex := hex.EncodeToString([]byte("01234567890123456789"))
	peerB_IDHex := hex.EncodeToString([]byte("BBBBBBBBBBBBBBBBBBBB"))
	peerC_IDHex := hex.EncodeToString([]byte("CCCCCCCCCCCCCCCCCCCC"))

	dbPeers := []*model.Peer{
		{
			PeerID: excludePeerIDHex,
			IP:     "200.200.200.1",
			Port:   1000,
		},
		{
			PeerID: peerB_IDHex,
			IP:     "200.200.200.1",
			Port:   2000,
		},
		{
			PeerID: peerC_IDHex,
			IP:     "100.100.100.1",
			Port:   3000,
		},
	}

	peers := BuildPeerList(excludePeerIDHex, dbPeers)
	assert.Len(t, peers, 2)
	assert.Equal(t, "200.200.200.1", peers[0]["ip"])
	assert.Equal(t, "100.100.100.1", peers[1]["ip"])
}

func TestBuildCompactPeerList(t *testing.T) {
	excludePeerIDHex := hex.EncodeToString([]byte("01234567890123456789"))
	peerB_IDHex := hex.EncodeToString([]byte("BBBBBBBBBBBBBBBBBBBB"))
	peerC_IDHex := hex.EncodeToString([]byte("CCCCCCCCCCCCCCCCCCCC"))

	dbPeers := []*model.Peer{
		{
			PeerID: excludePeerIDHex,
			IP:     "200.200.200.1",
			Port:   1000,
		},
		{
			PeerID: peerB_IDHex,
			IP:     "200.200.200.1",
			Port:   2000,
		},
		{
			PeerID: peerC_IDHex,
			IP:     "100.100.100.1",
			Port:   3000,
		},
	}

	compactStr := BuildCompactPeerList(excludePeerIDHex, dbPeers)
	assert.Len(t, compactStr, 12) // Peer B (6) + Peer C (6)

	assert.Equal(t, byte(200), compactStr[0])
	assert.Equal(t, byte(200), compactStr[1])
	assert.Equal(t, byte(200), compactStr[2])
	assert.Equal(t, byte(1), compactStr[3])
	assert.Equal(t, byte(2000>>8), compactStr[4])
	assert.Equal(t, byte(2000&0xFF), compactStr[5])

	assert.Equal(t, byte(100), compactStr[6])
	assert.Equal(t, byte(100), compactStr[7])
	assert.Equal(t, byte(100), compactStr[8])
	assert.Equal(t, byte(1), compactStr[9])
	assert.Equal(t, byte(3000>>8), compactStr[10])
	assert.Equal(t, byte(3000&0xFF), compactStr[11])
}
