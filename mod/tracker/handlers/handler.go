package handlers

import (
	"encoding/hex"
	"hash/fnv"
	"log/slog"
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
	systemSetting "github.com/GoldenSheep402/Hermes/mod/system/setting"
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
	// maxScrapeInfoHashes limits scrape request size (BEP common practice).
	maxScrapeInfoHashes = 64
)

// Registry registers BitTorrent tracker HTTP routes. Multi-site isolation uses RedisKeyPrefix
// in global config (TrackerV1.RedisKeyPrefix); passkey segments should be redacted in access logs
// (see RedactPasskeyInPath and deployment docs).
func Registry(jinE *jin.Engine) {
	jinE.Use(trackerAccessMetaMiddleware())
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

func trackerAccessMetaMiddleware() jin.HandlerFunc {
	return func(c *jin.Context) {
		c.Next()
		// Structured log (passkey-redacted path) for observability; uses default slog logger.
		path := RedactPasskeyInPath(c.Request.URL.Path)
		status := c.Writer.Status()
		if status >= 400 {
			slog.Warn("tracker_http", "path", path, "status", status)
		} else {
			slog.Debug("tracker_http", "path", path, "status", status)
		}
	}
}

func Announce(c *jin.Context) {
	ctx := c.Request.Context()
	passkey := c.Params.ByName("passkey")
	if passkey == "" {
		BencodeError(c, "Missing passkey")
		return
	}

	// 1. Authenticate user by passkey (Redis-backed cache)
	user, err := userDao.User.GetByPasskeyCached(ctx, passkey)
	if err != nil || user == nil {
		BencodeError(c, "Invalid passkey")
		return
	}
	if !user.IsEnabled {
		BencodeError(c, "User disabled")
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

	// 3. Find Torrent by InfoHash (Redis-backed cache)
	torrent, err := torrentDao.Torrent.GetByHashCached(ctx, infoHashHex)
	if err != nil {
		BencodeError(c, "Torrent not registered")
		return
	}

	// Manual parsing instead of ShouldBindQuery due to jin limitation and binary query
	port, _ := strconv.Atoi(q.Get("port"))
	if port < 1 || port > 65535 {
		BencodeError(c, "Invalid port")
		return
	}
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

	// Extract client public IP (LAN hairpin removed for single-site public tracker).
	var trustedProxies []string
	if cfg := conf.Get(); cfg != nil {
		trustedProxies = cfg.TrackerV1.TrustedProxyCIDRs
	}

	realIP := GetClientIP(c.Request, trustedProxies)

	// 4. Update Peer info in Redis
	isSeeder := left == 0

	peerIDHex := hex.EncodeToString([]byte(peerIDRaw))
	now := time.Now()

	peer := &trackerModel.Peer{
		Model:      stdao.Model{ID: ulid.Make().String()},
		TorrentID:  torrent.ID,
		UserID:     user.ID,
		PeerID:     peerIDHex,
		IP:         realIP,
		LanIP:      "",
		Port:       port,
		Uploaded:   uploaded,
		Downloaded: downloaded,
		Left:       left,
		Agent:      c.Request.Header.Get("User-Agent"),
		IsSeeder:   isSeeder,
		StartedAt:  now,
		LastAction: now,
	}
	announceResult, err := trackerDao.Traffic.ApplyAnnounce(ctx, peer, event)
	if err != nil {
		BencodeError(c, "Temporary tracker failure")
		return
	}
	if !announceResult.StartedAt.IsZero() {
		peer.StartedAt = announceResult.StartedAt
	}
	peer.Uploaded = announceResult.StoredUpload
	peer.Downloaded = announceResult.StoredDownload
	peer.Left = announceResult.StoredLeft
	peer.IsSeeder = peer.Left == 0

	if event == "completed" {
		finishedAt := announceResult.FinishedAt
		if finishedAt.IsZero() {
			finishedAt = time.Now()
		}
		snatch := &trackerModel.Snatch{
			Model:      stdao.Model{ID: ulid.Make().String()},
			TorrentID:  torrent.ID,
			UserID:     user.ID,
			Uploaded:   peer.Uploaded,
			SeedTime:   0,
			Downloaded: peer.Downloaded,
			IsActive:   true,
			FinishedAt: &finishedAt,
			LastAction: finishedAt,
		}
		if err := trackerDao.Snatch.UpdateOrCreate(ctx, snatch); err != nil {
			BencodeError(c, "Temporary tracker failure")
			return
		}
	}

	// 5. Build list of peers to return
	limit := sanitizeNumWant(numWant)
	var selectedPeers []*trackerModel.Peer
	if event != "stopped" {
		dbPeers, sampleErr := trackerDao.Peer.GetPeersForTorrentSample(ctx, torrent.ID, calcPeerFetchCount(limit), peer.PeerID, peer.LastAction)
		if sampleErr != nil {
			BencodeError(c, "Temporary tracker failure")
			return
		}
		selectedPeers = SelectPeersForResponse(peer, dbPeers, limit)
	} else {
		selectedPeers = []*trackerModel.Peer{}
	}
	stats, statsErr := trackerDao.Traffic.GetTorrentSnapshot(ctx, torrent.ID)
	if statsErr != nil {
		BencodeError(c, "Temporary tracker failure")
		return
	}

	// Build bencode response
	announceInterval := systemSetting.TrackerAnnounceIntervalValue(ctx)
	minInterval := announceInterval / 3
	if minInterval < 60 {
		minInterval = 60
	}
	resp := map[string]interface{}{
		"interval":     announceInterval,
		"min interval": minInterval,
		"complete":     stats.SeedCount,
		"incomplete":   stats.LeechCount,
	}

	if compact != 0 {
		resp["peers"] = BuildCompactPeerList(peer.PeerID, selectedPeers)
	} else {
		// Dictionary format
		resp["peers"] = BuildPeerList(peer.PeerID, selectedPeers)
	}

	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusOK)
	_ = bencode.NewEncoder(c.Writer).Encode(resp)
}

func Scrape(c *jin.Context) {
	passkey := c.Params.ByName("passkey")
	if passkey == "" {
		BencodeError(c, "Missing passkey")
		return
	}

	// Authenticate
	user, err := userDao.User.GetByPasskeyCached(c.Request.Context(), passkey)
	if err != nil || user == nil {
		BencodeError(c, "Invalid passkey")
		return
	}
	if !user.IsEnabled {
		BencodeError(c, "User disabled")
		return
	}

	q := c.Request.URL.Query()
	infoHashesRaw := q["info_hash"]
	if len(infoHashesRaw) > maxScrapeInfoHashes {
		infoHashesRaw = infoHashesRaw[:maxScrapeInfoHashes]
	}

	filesMap := map[string]interface{}{}

	for _, rawHash := range infoHashesRaw {
		if len(rawHash) != 20 {
			continue
		}
		hashHex := hex.EncodeToString([]byte(rawHash))

		torrent, err := torrentDao.Torrent.GetByHashCached(c.Request.Context(), hashHex)
		if err != nil || torrent == nil {
			continue
		}
		stats, err := trackerDao.Traffic.GetTorrentSnapshot(c.Request.Context(), torrent.ID)
		if err != nil {
			continue
		}

		filesMap[rawHash] = map[string]interface{}{
			"complete":   stats.SeedCount,
			"downloaded": stats.SnatchCount,
			"incomplete": stats.LeechCount,
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
		// the cumulative counter may double count old traffic. Accept full value on
		// started or empty event (matches announce.lua first-session behavior).
		if event == "started" || event == "" {
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
	case "":
		return ""
	case "start":
		return "started"
	case "stop":
		return "stopped"
	case "complete":
		return "completed"
	case "started", "stopped", "completed":
		return event
	default:
		// Treat unknown client events as periodic announce (empty).
		return ""
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
// - prefer active peers
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
		if net.ParseIP(p.IP) == nil {
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

func peerSubnetKey(ip string) string {
	parsed := net.ParseIP(ip).To4()
	if parsed == nil {
		return ip
	}
	return strconv.Itoa(int(parsed[0])) + "." + strconv.Itoa(int(parsed[1])) + "." + strconv.Itoa(int(parsed[2]))
}

// BuildPeerList constructs the dictionary peer list using public IPs only.
func BuildPeerList(excludePeerID string, dbPeers []*trackerModel.Peer) []map[string]interface{} {
	var bPeers []map[string]interface{}
	for _, p := range dbPeers {
		if p.PeerID == excludePeerID {
			continue
		}
		rawPeerID, err := hex.DecodeString(p.PeerID)
		if err != nil {
			continue
		}

		bPeers = append(bPeers, map[string]interface{}{
			"peer id": string(rawPeerID),
			"ip":      p.IP,
			"port":    p.Port,
		})
	}
	if bPeers == nil {
		bPeers = []map[string]interface{}{}
	}
	return bPeers
}

// BuildCompactPeerList constructs the binary compact representation of peers.
// Only IPv4 endpoints are encoded (To4); IPv6-only peers are omitted in compact mode.
func BuildCompactPeerList(excludePeerID string, dbPeers []*trackerModel.Peer) string {
	var buf []byte
	for _, p := range dbPeers {
		if p.PeerID == excludePeerID {
			continue
		}

		ipBytes := net.ParseIP(p.IP).To4()
		if ipBytes != nil {
			buf = append(buf, ipBytes...)
			buf = append(buf, byte(p.Port>>8), byte(p.Port&0xFF))
		}
	}
	return string(buf)
}
