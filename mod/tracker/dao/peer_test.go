package dao

import (
	"context"
	"fmt"
	"testing"
	"time"

	trackerModel "github.com/GoldenSheep402/Hermes/mod/tracker/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func TestGetPeersForTorrentSampleIncludesOlderPeers(t *testing.T) {
	_, peerDAO, cleanup := newTrackerTrafficTestStore(t)
	defer cleanup()

	ctx := context.Background()
	torrentID := "torrent-sample"
	baseTime := time.Now().UTC().Truncate(time.Second)

	for i := 0; i < 20; i++ {
		peerData := &trackerModel.Peer{
			Model:      stdao.Model{ID: ulid.Make().String()},
			TorrentID:  torrentID,
			UserID:     fmt.Sprintf("user-%02d", i),
			PeerID:     fmt.Sprintf("peer-%02d", i),
			IP:         fmt.Sprintf("10.0.0.%d", i+1),
			Port:       6000 + i,
			Uploaded:   int64(i),
			Downloaded: int64(i),
			Left:       1024,
			LastAction: baseTime.Add(-time.Duration(i) * time.Minute),
		}
		require.NoError(t, peerDAO.Upsert(ctx, peerData))
	}

	peers, err := peerDAO.GetPeersForTorrentSample(ctx, torrentID, 10, "requester-peer", baseTime)
	require.NoError(t, err)
	require.Len(t, peers, 10)

	ids := make([]string, 0, len(peers))
	hasOlderPeer := false
	for _, peer := range peers {
		ids = append(ids, peer.PeerID)
		// Lexicographic order matches numeric order for peer-00 .. peer-19 with this padding.
		if peer.PeerID >= "peer-06" {
			hasOlderPeer = true
		}
	}

	require.Contains(t, ids, "peer-00")
	require.Contains(t, ids, "peer-05")
	require.True(t, hasOlderPeer, "sample should include peers outside the most recent window")
}

func TestGetPeersForTorrentReturnsAllWhenLimitNonPositive(t *testing.T) {
	_, peerDAO, cleanup := newTrackerTrafficTestStore(t)
	defer cleanup()

	ctx := context.Background()
	torrentID := "torrent-all"
	baseTime := time.Now().UTC().Truncate(time.Second)

	for i := 0; i < 25; i++ {
		peerData := &trackerModel.Peer{
			Model:      stdao.Model{ID: ulid.Make().String()},
			TorrentID:  torrentID,
			UserID:     fmt.Sprintf("user-all-%02d", i),
			PeerID:     fmt.Sprintf("peer-all-%02d", i),
			IP:         fmt.Sprintf("10.1.0.%d", i+1),
			Port:       7000 + i,
			Uploaded:   int64(i),
			Downloaded: int64(i),
			Left:       1024,
			LastAction: baseTime.Add(-time.Duration(i) * time.Minute),
		}
		require.NoError(t, peerDAO.Upsert(ctx, peerData))
	}

	peers, err := peerDAO.GetPeersForTorrent(ctx, torrentID, 0)
	require.NoError(t, err)
	require.Len(t, peers, 25)
}
