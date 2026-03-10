package handlers

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	trackerModel "github.com/GoldenSheep402/Hermes/mod/tracker/model"
)

func TestComputeCounterDelta(t *testing.T) {
	getUploaded := func(p *trackerModel.Peer) int64 { return p.Uploaded }

	t.Run("first_started_uses_current", func(t *testing.T) {
		delta := computeCounterDelta(nil, 120, "started", getUploaded)
		require.Equal(t, int64(120), delta)
	})

	t.Run("first_start_alias_uses_current", func(t *testing.T) {
		delta := computeCounterDelta(nil, 120, normalizeAnnounceEvent("start"), getUploaded)
		require.Equal(t, int64(120), delta)
	})

	t.Run("first_no_event_ignored", func(t *testing.T) {
		delta := computeCounterDelta(nil, 120, "", getUploaded)
		require.Equal(t, int64(0), delta)
	})

	t.Run("normal_increment", func(t *testing.T) {
		last := &trackerModel.Peer{Uploaded: 100}
		delta := computeCounterDelta(last, 130, "", getUploaded)
		require.Equal(t, int64(30), delta)
	})

	t.Run("counter_reset_on_started", func(t *testing.T) {
		last := &trackerModel.Peer{Uploaded: 300}
		delta := computeCounterDelta(last, 50, "started", getUploaded)
		require.Equal(t, int64(50), delta)
	})

	t.Run("counter_reset_on_start_alias", func(t *testing.T) {
		last := &trackerModel.Peer{Uploaded: 300}
		delta := computeCounterDelta(last, 50, normalizeAnnounceEvent("start"), getUploaded)
		require.Equal(t, int64(50), delta)
	})

	t.Run("counter_reset_without_started", func(t *testing.T) {
		last := &trackerModel.Peer{Uploaded: 300}
		delta := computeCounterDelta(last, 50, "", getUploaded)
		require.Equal(t, int64(0), delta)
	})
}

func TestNormalizeAnnounceEvent(t *testing.T) {
	require.Equal(t, "started", normalizeAnnounceEvent("start"))
	require.Equal(t, "started", normalizeAnnounceEvent("started"))
	require.Equal(t, "stopped", normalizeAnnounceEvent("stop"))
	require.Equal(t, "completed", normalizeAnnounceEvent("complete"))
	require.Equal(t, "started", normalizeAnnounceEvent(" Started "))
	require.Equal(t, "", normalizeAnnounceEvent(""))
}

func TestSanitizeNumWant(t *testing.T) {
	require.Equal(t, defaultNumWant, sanitizeNumWant(0))
	require.Equal(t, defaultNumWant, sanitizeNumWant(-10))
	require.Equal(t, 25, sanitizeNumWant(25))
	require.Equal(t, maxNumWant, sanitizeNumWant(maxNumWant+1))
}

func TestFilterPeerCandidates(t *testing.T) {
	now := time.Now()
	requester := makePeer("req", "u1", "10.0.0.1", "", 6881, false, 100, 10, 20, now)

	peers := []*trackerModel.Peer{
		nil,
		makePeer("req", "u2", "1.1.1.1", "", 6881, true, 0, 0, 0, now),   // self by peer_id
		makePeer("p1", "u1", "1.1.1.2", "", 6881, true, 0, 0, 0, now),    // same account
		makePeer("p2", "u2", "1.1.1.3", "", 0, true, 0, 0, 0, now),       // invalid port
		makePeer("p3", "u3", "1.1.1.4", "", 6881, true, 0, 0, 0, now),    // keep
		makePeer("p3", "u4", "1.1.1.5", "", 6881, true, 0, 0, 0, now),    // duplicate peer_id
		makePeer("p4", "u5", "1.1.1.4", "", 6881, false, 0, 0, 0, now),   // duplicate endpoint
		makePeer("p5", "u6", "2.2.2.2", "", 6882, false, 0, 0, 0, now),   // keep
		makePeer("", "u7", "3.3.3.3", "", 6881, false, 0, 0, 0, now),     // empty peer_id
		makePeer("p6", "u8", "4.4.4.4", "", 65536, false, 0, 0, 0, now),  // invalid port
		makePeer("p7", "u9", "5.5.5.5", "", -1, false, 0, 0, 0, now),     // invalid port
		makePeer("p8", "u10", "6.6.6.6", "", 65535, false, 0, 0, 0, now), // keep
	}

	filtered := filterPeerCandidates(requester, peers)
	require.Len(t, filtered, 3)

	ids := make([]string, 0, len(filtered))
	for _, p := range filtered {
		ids = append(ids, p.PeerID)
	}
	require.ElementsMatch(t, []string{"p3", "p5", "p8"}, ids)
}

func TestSelectPeersForResponseLeecherPrefersSeeders(t *testing.T) {
	now := time.Now()
	requester := makePeer("req", "u1", "10.0.0.1", "", 6881, false, 100, 0, 0, now)

	candidates := []*trackerModel.Peer{
		makePeer("s1", "u2", "1.1.1.1", "", 6001, true, 0, 1000, 200, now),
		makePeer("s2", "u3", "2.2.2.2", "", 6002, true, 0, 900, 300, now.Add(-time.Second)),
		makePeer("s3", "u4", "3.3.3.3", "", 6003, true, 0, 800, 100, now.Add(-2*time.Second)),
		makePeer("l1", "u5", "4.4.4.4", "", 6004, false, 200, 10, 20, now),
		makePeer("l2", "u6", "5.5.5.5", "", 6005, false, 300, 20, 30, now),
	}

	selected := SelectPeersForResponse(requester, candidates, 3)
	require.Len(t, selected, 3)
	for _, p := range selected {
		require.True(t, p.IsSeeder)
	}
}

func TestSelectPeersForResponseSeederPrefersLeechers(t *testing.T) {
	now := time.Now()
	requester := makePeer("req", "u1", "10.0.0.1", "", 6881, true, 0, 0, 0, now)

	candidates := []*trackerModel.Peer{
		makePeer("l1", "u2", "1.1.1.1", "", 6001, false, 900, 100, 500, now),
		makePeer("l2", "u3", "2.2.2.2", "", 6002, false, 800, 90, 400, now.Add(-time.Second)),
		makePeer("l3", "u4", "3.3.3.3", "", 6003, false, 700, 80, 300, now.Add(-2*time.Second)),
		makePeer("s1", "u5", "4.4.4.4", "", 6004, true, 0, 1000, 200, now),
		makePeer("s2", "u6", "5.5.5.5", "", 6005, true, 0, 900, 300, now),
	}

	selected := SelectPeersForResponse(requester, candidates, 3)
	require.Len(t, selected, 3)
	for _, p := range selected {
		require.False(t, p.IsSeeder)
	}
}

func TestSelectPeersForResponseSubnetDiversity(t *testing.T) {
	now := time.Now()
	requester := makePeer("req", "u1", "10.0.0.1", "", 6881, false, 100, 0, 0, now)

	candidates := []*trackerModel.Peer{
		makePeer("s1", "u2", "11.11.11.1", "", 6001, true, 0, 1000, 200, now),
		makePeer("s2", "u3", "11.11.11.2", "", 6002, true, 0, 900, 300, now),
		makePeer("s3", "u4", "11.11.11.3", "", 6003, true, 0, 800, 100, now),
		makePeer("s4", "u5", "12.12.12.1", "", 6004, true, 0, 700, 100, now),
		makePeer("s5", "u6", "13.13.13.1", "", 6005, true, 0, 600, 100, now),
	}

	selected := SelectPeersForResponse(requester, candidates, 4)
	require.Len(t, selected, 4)

	subnetCounts := make(map[string]int)
	for _, p := range selected {
		subnetCounts[subnet24KeyForTest(p.IP)]++
	}
	require.LessOrEqual(t, subnetCounts["11.11.11"], maxPeersPerSubnetInPhase)
}

func TestSelectPeersForResponseLanAffinity(t *testing.T) {
	now := time.Now()
	requester := makePeer("req", "u1", "20.20.20.1", "192.168.1.2", 6881, false, 100, 0, 0, now)

	lanPeer := makePeer("lan", "u2", "20.20.20.1", "192.168.1.3", 6001, true, 0, 500, 100, now)
	wanPeer := makePeer("wan", "u3", "30.30.30.1", "", 6002, true, 0, 500, 100, now)

	selected := SelectPeersForResponse(requester, []*trackerModel.Peer{lanPeer, wanPeer}, 1)
	require.Len(t, selected, 1)
	require.Equal(t, "lan", selected[0].PeerID)
}

func makePeer(
	peerID, userID, ip, lanIP string,
	port int,
	isSeeder bool,
	left, uploaded, downloaded int64,
	lastAction time.Time,
) *trackerModel.Peer {
	return &trackerModel.Peer{
		PeerID:     peerID,
		UserID:     userID,
		IP:         ip,
		LanIP:      lanIP,
		Port:       port,
		IsSeeder:   isSeeder,
		Left:       left,
		Uploaded:   uploaded,
		Downloaded: downloaded,
		LastAction: lastAction,
	}
}

func subnet24KeyForTest(ip string) string {
	parts := []rune(ip)
	dot := 0
	for i, ch := range parts {
		if ch == '.' {
			dot++
			if dot == 3 {
				return string(parts[:i])
			}
		}
	}
	if ip == "" {
		return ""
	}
	// keep fallback behavior readable in assertions
	return strconv.Quote(ip)
}
