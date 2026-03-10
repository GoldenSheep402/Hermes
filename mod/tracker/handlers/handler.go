package handlers

import (
	"encoding/hex"
	"hash/fnv"
	"math/rand"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/juanjiTech/jin"
	"github.com/oklog/ulid/v2"
	"github.com/zeebo/bencode"

	"github.com/GoldenSheep402/Hermes/conf"
	torrentDao "github.com/GoldenSheep402/Hermes/mod/torrent/dao"
	trackerDao "github.com/GoldenSheep402/Hermes/mod/tracker/dao"
	trackerModel "github.com/GoldenSheep402/Hermes/mod/tracker/model"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

const (
	defaultNumWant           = 50
	maxNumWant               = 200
	defaultPeerFetchCount    = 120
	maxPeerFetchCount        = 600
	maxPeersPerSubnetInPhase = 2
)

// Registry TODO: multi tracker support
func Registry(jinE *jin.Engine) {
	jinE.GET("/announce/:passkey", Announce)
	jinE.GET("/scrape/:passkey", Scrape)
	// Backward-compatible routes for deployments that proxy tracker through /api/*.
	jinE.GET("/api/announce/:passkey", Announce)
	jinE.GET("/api/scrape/:passkey", Scrape)
}

// BencodeError sends an error back to the torrent client in bencode format.
func BencodeError(c *jin.Context, msg string) {
	resp := map[string]interface{}{
		"failure reason": msg,
	}
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusOK)
	bencode.NewEncoder(c.Writer).Encode(resp)
}

func Announce(c *jin.Context) {
	ctx := c.Request.Context()
	passkey := c.Params.ByName("passkey")
	if passkey == "" {
		BencodeError(c, "Missing passkey")
		return
	}

	// 1. Authenticate user by passkey
	user, err := userDao.User.GetByPasskey(ctx, passkey)
	if err != nil || user == nil {
		BencodeError(c, "Invalid passkey")
		return
	}

	// 2. Parse Tracker Request Parameters (info_hash is raw 20 bytes)
	q := c.Request.URL.Query()
	infoHashRaw := q.Get("info_hash")
	if len(infoHashRaw) != 20 {
		BencodeError(c, "Invalid info_hash length")
		return
	}
	infoHashHex := hex.EncodeToString([]byte(infoHashRaw))

	peerIDRaw := q.Get("peer_id")
	if len(peerIDRaw) != 20 {
		BencodeError(c, "Invalid peer_id length")
		return
	}

	// 3. Find Torrent by InfoHash
	torrent, err := torrentDao.Torrent.GetByHash(ctx, infoHashHex)
	if err != nil {
		BencodeError(c, "Torrent not registered")
		return
	}

	// Manual parsing instead of ShouldBindQuery due to jin limitation and binary query
	port, _ := strconv.Atoi(q.Get("port"))
	uploaded, _ := strconv.ParseInt(q.Get("uploaded"), 10, 64)
	downloaded, _ := strconv.ParseInt(q.Get("downloaded"), 10, 64)
	left, _ := strconv.ParseInt(q.Get("left"), 10, 64)
	if uploaded < 0 {
		uploaded = 0
	}
	if downloaded < 0 {
		downloaded = 0
	}
	if left < 0 {
		left = 0
	}
	compact, _ := strconv.Atoi(q.Get("compact"))
	numWant, err := strconv.Atoi(q.Get("numwant"))
	if err != nil {
		numWant = 50
	}
	event := normalizeAnnounceEvent(q.Get("event"))

	// Extract Real IP and LanIP (if any)
	// We dynamically load AllowedSubnets from the global config so it responds to hot-reloads
	var allowedSubnets []string
	if c := conf.Get(); c != nil {
		allowedSubnets = c.TrackerV1.AllowedSubnets
	}

	realIP, lanIP := ExtractIPs(c.Request, allowedSubnets)

	// 4. Update Peer info in DB (or Redis eventually)
	isSeeder := left == 0

	peerIDHex := hex.EncodeToString([]byte(peerIDRaw))
	lastPeer, _ := trackerDao.Peer.Get(ctx, torrent.ID, peerIDHex)
	now := time.Now()
	startedAt := now
	if lastPeer != nil && !lastPeer.StartedAt.IsZero() {
		startedAt = lastPeer.StartedAt
	}
	if event == "started" || lastPeer == nil {
		startedAt = now
	}

	peer := &trackerModel.Peer{
		Model:      stdao.Model{ID: ulid.Make().String()},
		TorrentID:  torrent.ID,
		UserID:     user.ID,
		PeerID:     peerIDHex,
		IP:         realIP,
		LanIP:      lanIP,
		Port:       port,
		Uploaded:   uploaded,
		Downloaded: downloaded,
		Left:       left,
		Agent:      c.Request.Header.Get("User-Agent"),
		IsSeeder:   isSeeder,
		StartedAt:  startedAt,
		LastAction: now,
	}
	uploadDelta := computeCounterDelta(lastPeer, uploaded, event, func(p *trackerModel.Peer) int64 {
		return p.Uploaded
	})
	downloadDelta := computeCounterDelta(lastPeer, downloaded, event, func(p *trackerModel.Peer) int64 {
		return p.Downloaded
	})
	isActive := event != "stopped"

	if event == "stopped" {
		// Remove peer
		_ = trackerDao.Peer.DeletePeer(ctx, torrent.ID, peer.PeerID)
	} else {
		// Upsert peer
		_ = trackerDao.Peer.Upsert(ctx, peer)

		// Create snatch record if event == completed
		if event == "completed" {
			now := time.Now()
			snatch := &trackerModel.Snatch{
				Model:      stdao.Model{ID: ulid.Make().String()},
				TorrentID:  torrent.ID,
				UserID:     user.ID,
				Uploaded:   uploaded,
				Downloaded: downloaded,
				IsActive:   true,
				FinishedAt: &now,
				LastAction: now,
			}
			_ = trackerDao.Snatch.UpdateOrCreate(ctx, snatch)
		}
	}
	_ = trackerDao.Traffic.RecordDelta(ctx, user.ID, torrent.ID, uploadDelta, downloadDelta, isActive, isSeeder, peer.LastAction)

	// 5. Build list of peers to return
	limit := sanitizeNumWant(numWant)
	dbPeers, _ := trackerDao.Peer.GetPeersForTorrent(ctx, torrent.ID, calcPeerFetchCount(limit))
	selectedPeers := SelectPeersForResponse(peer, dbPeers, limit)

	// Build bencode response
	resp := map[string]interface{}{
		"interval":     1800, // 30 minutes
		"min interval": 600,  // 10 minutes
		"complete":     torrent.SeedCount,
		"incomplete":   torrent.LeechCount,
	}

	if compact != 0 {
		resp["peers"] = BuildCompactPeerList(realIP, lanIP, peer.PeerID, selectedPeers)
	} else {
		// Dictionary format
		resp["peers"] = BuildPeerList(realIP, lanIP, peer.PeerID, selectedPeers)
	}

	bencode.NewEncoder(c.Writer).Encode(resp)
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusOK)
}

func Scrape(c *jin.Context) {
	passkey := c.Params.ByName("passkey")
	if passkey == "" {
		BencodeError(c, "Missing passkey")
		return
	}

	// Authenticate
	user, err := userDao.User.GetByPasskey(c.Request.Context(), passkey)
	if err != nil || user == nil {
		BencodeError(c, "Invalid passkey")
		return
	}

	q := c.Request.URL.Query()
	infoHashesRaw := q["info_hash"]

	filesMap := map[string]interface{}{}

	for _, rawHash := range infoHashesRaw {
		if len(rawHash) != 20 {
			continue
		}
		hashHex := hex.EncodeToString([]byte(rawHash))

		torrent, err := torrentDao.Torrent.GetByHash(c.Request.Context(), hashHex)
		if err != nil || torrent == nil {
			continue
		}

		filesMap[rawHash] = map[string]interface{}{
			"complete":   torrent.SeedCount,
			"downloaded": torrent.SnatchCount,
			"incomplete": torrent.LeechCount,
		}
	}

	resp := map[string]interface{}{
		"files": filesMap,
	}

	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusOK)
	bencode.NewEncoder(c.Writer).Encode(resp)
}

func computeCounterDelta(lastPeer *trackerModel.Peer, current int64, event string, selector func(*trackerModel.Peer) int64) int64 {
	if current < 0 {
		return 0
	}
	if lastPeer == nil {
		// When peer snapshot is missing (e.g. redis restart/expiry), blindly trusting
		// the cumulative counter may double count old traffic. Accept full value only
		// on explicit session boundaries.
		if event == "started" || event == "completed" {
			return current
		}
		return 0
	}

	previous := selector(lastPeer)
	if previous < 0 {
		previous = 0
	}
	if current >= previous {
		return current - previous
	}
	// Clients may reset counters on a fresh session; treat "started" as a new base.
	if event == "started" {
		return current
	}
	return 0
}

func normalizeAnnounceEvent(raw string) string {
	event := strings.ToLower(strings.TrimSpace(raw))
	switch event {
	case "start":
		return "started"
	case "stop":
		return "stopped"
	case "complete":
		return "completed"
	default:
		return event
	}
}

func sanitizeNumWant(numWant int) int {
	if numWant <= 0 {
		return defaultNumWant
	}
	if numWant > maxNumWant {
		return maxNumWant
	}
	return numWant
}

func calcPeerFetchCount(limit int) int {
	limit = sanitizeNumWant(limit)
	fetchCount := limit * 4
	if fetchCount < defaultPeerFetchCount {
		fetchCount = defaultPeerFetchCount
	}
	if fetchCount > maxPeerFetchCount {
		fetchCount = maxPeerFetchCount
	}
	return fetchCount
}

// SelectPeersForResponse applies a smarter peer selection strategy:
// - filter duplicates and same-account peers
// - prefer opposite role (leecher gets seeders, seeder gets leechers)
// - prefer LAN-affinity and active peers
// - keep subnet diversity to avoid hotspot peers
func SelectPeersForResponse(requester *trackerModel.Peer, dbPeers []*trackerModel.Peer, limit int) []*trackerModel.Peer {
	limit = sanitizeNumWant(limit)
	if requester == nil || len(dbPeers) == 0 || limit <= 0 {
		return []*trackerModel.Peer{}
	}

	candidates := filterPeerCandidates(requester, dbPeers)
	if len(candidates) == 0 {
		return []*trackerModel.Peer{}
	}

	rotatePeersForFairness(candidates, requester.PeerID)

	preferred, fallback := splitPeersByRole(requester, candidates)
	sortPeersWithStrategy(requester, preferred)
	sortPeersWithStrategy(requester, fallback)

	selected := make([]*trackerModel.Peer, 0, limit)
	selectedSet := make(map[string]struct{}, limit)
	subnetCount := make(map[string]int, limit)

	selected = appendPeersWithSubnetDiversity(selected, preferred, limit, subnetCount, selectedSet)
	if len(selected) < limit {
		selected = appendPeersWithSubnetDiversity(selected, fallback, limit, subnetCount, selectedSet)
	}
	return selected
}

func filterPeerCandidates(requester *trackerModel.Peer, peers []*trackerModel.Peer) []*trackerModel.Peer {
	result := make([]*trackerModel.Peer, 0, len(peers))
	seenPeerID := make(map[string]struct{}, len(peers))
	seenEndpoint := make(map[string]struct{}, len(peers))

	for _, p := range peers {
		if p == nil {
			continue
		}
		if p.PeerID == "" || p.PeerID == requester.PeerID {
			continue
		}
		if p.UserID != "" && requester.UserID != "" && p.UserID == requester.UserID {
			continue
		}
		if p.Port <= 0 || p.Port > 65535 {
			continue
		}
		if _, exists := seenPeerID[p.PeerID]; exists {
			continue
		}
		endpointKey := p.IP + ":" + strconv.Itoa(p.Port)
		if _, exists := seenEndpoint[endpointKey]; exists {
			continue
		}
		seenPeerID[p.PeerID] = struct{}{}
		seenEndpoint[endpointKey] = struct{}{}
		result = append(result, p)
	}
	return result
}

func splitPeersByRole(requester *trackerModel.Peer, peers []*trackerModel.Peer) ([]*trackerModel.Peer, []*trackerModel.Peer) {
	preferred := make([]*trackerModel.Peer, 0, len(peers))
	fallback := make([]*trackerModel.Peer, 0, len(peers))

	for _, p := range peers {
		// Seeder should receive leechers first for better upload opportunities.
		// Leecher should receive seeders first for better download opportunities.
		if requester.IsSeeder {
			if !p.IsSeeder {
				preferred = append(preferred, p)
			} else {
				fallback = append(fallback, p)
			}
		} else {
			if p.IsSeeder {
				preferred = append(preferred, p)
			} else {
				fallback = append(fallback, p)
			}
		}
	}
	return preferred, fallback
}

func rotatePeersForFairness(peers []*trackerModel.Peer, requesterPeerID string) {
	if len(peers) <= 1 {
		return
	}
	h := fnv.New64a()
	_, _ = h.Write([]byte(requesterPeerID))
	_, _ = h.Write([]byte(":"))
	_, _ = h.Write([]byte(time.Now().UTC().Format("200601021504"))) // minute-level rotation
	seed := int64(h.Sum64() & 0x7fffffffffffffff)
	r := rand.New(rand.NewSource(seed))
	r.Shuffle(len(peers), func(i, j int) {
		peers[i], peers[j] = peers[j], peers[i]
	})
}

func sortPeersWithStrategy(requester *trackerModel.Peer, peers []*trackerModel.Peer) {
	if len(peers) <= 1 {
		return
	}

	sort.SliceStable(peers, func(i, j int) bool {
		a := peers[i]
		b := peers[j]

		// Prefer LAN-affinity first when available.
		aLAN := lanAffinityScore(requester, a)
		bLAN := lanAffinityScore(requester, b)
		if aLAN != bLAN {
			return aLAN > bLAN
		}

		if requester.IsSeeder {
			// For seeding peers, prioritize active leechers with more remaining data.
			if a.Left != b.Left {
				return a.Left > b.Left
			}
			if !a.LastAction.Equal(b.LastAction) {
				return a.LastAction.After(b.LastAction)
			}
			if a.Downloaded != b.Downloaded {
				return a.Downloaded > b.Downloaded
			}
		} else {
			// For leechers, prioritize active seeders with richer historical upload.
			if !a.LastAction.Equal(b.LastAction) {
				return a.LastAction.After(b.LastAction)
			}
			if a.Uploaded != b.Uploaded {
				return a.Uploaded > b.Uploaded
			}
			if a.Left != b.Left {
				return a.Left < b.Left
			}
		}

		return a.PeerID < b.PeerID
	})
}

func appendPeersWithSubnetDiversity(
	dst []*trackerModel.Peer,
	src []*trackerModel.Peer,
	limit int,
	subnetCount map[string]int,
	selectedSet map[string]struct{},
) []*trackerModel.Peer {
	if len(dst) >= limit || len(src) == 0 {
		return dst
	}

	// Phase 1: keep subnet diversity.
	for _, p := range src {
		if len(dst) >= limit {
			return dst
		}
		if _, exists := selectedSet[p.PeerID]; exists {
			continue
		}
		subnetKey := peerSubnetKey(p.IP)
		if subnetCount[subnetKey] >= maxPeersPerSubnetInPhase {
			continue
		}
		dst = append(dst, p)
		selectedSet[p.PeerID] = struct{}{}
		subnetCount[subnetKey]++
	}

	// Phase 2: fill remaining slots even if subnet quota is exceeded.
	for _, p := range src {
		if len(dst) >= limit {
			return dst
		}
		if _, exists := selectedSet[p.PeerID]; exists {
			continue
		}
		dst = append(dst, p)
		selectedSet[p.PeerID] = struct{}{}
		subnetCount[peerSubnetKey(p.IP)]++
	}
	return dst
}

func lanAffinityScore(requester, candidate *trackerModel.Peer) int {
	if requester == nil || candidate == nil {
		return 0
	}
	score := 0
	if requester.IP != "" && requester.IP == candidate.IP {
		score++
	}
	if requester.LanIP != "" && candidate.LanIP != "" && isSameSubnet24(requester.LanIP, candidate.LanIP) {
		score++
	}
	return score
}

func peerSubnetKey(ip string) string {
	parsed := net.ParseIP(ip).To4()
	if parsed == nil {
		return ip
	}
	return strconv.Itoa(int(parsed[0])) + "." + strconv.Itoa(int(parsed[1])) + "." + strconv.Itoa(int(parsed[2]))
}

// BuildPeerList constructs the list of peers to be returned to the client,
// handling LAN IP substitution only if peers share the same public IP AND
// their LAN IPs are in the same subnet (assumed /24 for typical homes).
// It returns BOTH the real IP and the LAN IP so clients can fallback to the public IP.
func BuildPeerList(clientRealIP, clientLanIP string, excludePeerID string, dbPeers []*trackerModel.Peer) []map[string]interface{} {
	var bPeers []map[string]interface{}
	for _, p := range dbPeers {
		if p.PeerID == excludePeerID {
			continue
		}
		rawPeerID, err := hex.DecodeString(p.PeerID)
		if err != nil {
			continue
		}

		// Always return the standard public IP
		bPeers = append(bPeers, map[string]interface{}{
			"peer id": string(rawPeerID),
			"ip":      p.IP,
			"port":    p.Port,
		})

		// If both peers share the exact same public (real) IP, and the target peer has a valid LAN IP registered
		// AND they are in the exact same /24 subnet for IPv4
		if p.LanIP != "" && clientLanIP != "" && p.IP == clientRealIP {
			if isSameSubnet24(p.LanIP, clientLanIP) {
				// Provide the LAN IP as an additional endpoint for fallback
				bPeers = append(bPeers, map[string]interface{}{
					"peer id": string(rawPeerID),
					"ip":      p.LanIP,
					"port":    p.Port,
				})
			}
		}
	}
	if bPeers == nil {
		bPeers = []map[string]interface{}{}
	}
	return bPeers
}

// BuildCompactPeerList constructs the binary compact representation of peers.
func BuildCompactPeerList(clientRealIP, clientLanIP string, excludePeerID string, dbPeers []*trackerModel.Peer) string {
	var buf []byte
	for _, p := range dbPeers {
		if p.PeerID == excludePeerID {
			continue
		}

		// Always add Public IP
		ipBytes := net.ParseIP(p.IP).To4()
		if ipBytes != nil {
			buf = append(buf, ipBytes...)
			buf = append(buf, byte(p.Port>>8), byte(p.Port&0xFF))
		}

		// Add LAN IP as an additional fallback endpoint
		if p.LanIP != "" && clientLanIP != "" && p.IP == clientRealIP {
			if isSameSubnet24(p.LanIP, clientLanIP) {
				lanIpBytes := net.ParseIP(p.LanIP).To4()
				if lanIpBytes != nil {
					buf = append(buf, lanIpBytes...)
					buf = append(buf, byte(p.Port>>8), byte(p.Port&0xFF))
				}
			}
		}
	}
	return string(buf)
}

func isSameSubnet24(ip1, ip2 string) bool {
	parsed1 := net.ParseIP(ip1).To4()
	parsed2 := net.ParseIP(ip2).To4()
	if parsed1 == nil || parsed2 == nil {
		return false
	}
	// Check if first 3 bytes (24 bits) match
	return parsed1[0] == parsed2[0] && parsed1[1] == parsed2[1] && parsed1[2] == parsed2[2]
}
