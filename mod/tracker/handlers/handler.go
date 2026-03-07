package handlers

import (
	"encoding/hex"
	"net"
	"net/http"
	"strconv"
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

// Registry TODO: multi tracker support
func Registry(jinE *jin.Engine) {
	jinE.GET("/announce/:passkey", Announce)
	jinE.GET("/scrape/:passkey", Scrape)
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
	passkey := c.Params.ByName("passkey")
	if passkey == "" {
		BencodeError(c, "Missing passkey")
		return
	}

	// 1. Authenticate user by passkey
	user, err := userDao.User.GetByPasskey(c.Request.Context(), passkey)
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
	torrent, err := torrentDao.Torrent.GetByHash(c.Request.Context(), infoHashHex)
	if err != nil {
		BencodeError(c, "Torrent not registered")
		return
	}

	// Manual parsing instead of ShouldBindQuery due to jin limitation and binary query
	port, _ := strconv.Atoi(q.Get("port"))
	uploaded, _ := strconv.ParseInt(q.Get("uploaded"), 10, 64)
	downloaded, _ := strconv.ParseInt(q.Get("downloaded"), 10, 64)
	left, _ := strconv.ParseInt(q.Get("left"), 10, 64)
	compact, _ := strconv.Atoi(q.Get("compact"))
	numWant, err := strconv.Atoi(q.Get("numwant"))
	if err != nil {
		numWant = 50
	}
	event := q.Get("event")

	// Extract Real IP and LanIP (if any)
	// We dynamically load AllowedSubnets from the global config so it responds to hot-reloads
	var allowedSubnets []string
	if c := conf.Get(); c != nil {
		allowedSubnets = c.TrackerV1.AllowedSubnets
	}

	realIP, lanIP := ExtractIPs(c.Request, allowedSubnets)

	// 4. Update Peer info in DB (or Redis eventually)
	isSeeder := left == 0

	peer := &trackerModel.Peer{
		Model:      stdao.Model{ID: ulid.Make().String()},
		TorrentID:  torrent.ID,
		UserID:     user.ID,
		PeerID:     hex.EncodeToString([]byte(peerIDRaw)),
		IP:         realIP,
		LanIP:      lanIP,
		Port:       port,
		Uploaded:   uploaded,
		Downloaded: downloaded,
		Left:       left,
		Agent:      c.Request.Header.Get("User-Agent"),
		IsSeeder:   isSeeder,
		LastAction: time.Now(),
	}

	if event == "stopped" {
		// Remove peer
		trackerDao.Peer.DeletePeer(c.Request.Context(), torrent.ID, peer.PeerID)
	} else {
		// Upsert peer
		trackerDao.Peer.Upsert(c.Request.Context(), peer)

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
			trackerDao.Snatch.UpdateOrCreate(c.Request.Context(), snatch)
		}
	}

	// 5. Build list of peers to return
	limit := numWant
	if limit <= 0 {
		limit = 50
	}
	dbPeers, _ := trackerDao.Peer.GetPeersForTorrent(c.Request.Context(), torrent.ID, limit)

	// Build bencode response
	resp := map[string]interface{}{
		"interval":     1800, // 30 minutes
		"min interval": 600,  // 10 minutes
		"complete":     torrent.SeedCount,
		"incomplete":   torrent.LeechCount,
	}

	if compact != 0 {
		resp["peers"] = BuildCompactPeerList(realIP, peer.PeerID, dbPeers)
	} else {
		// Dictionary format
		resp["peers"] = BuildPeerList(realIP, peer.PeerID, dbPeers)
	}

	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusOK)
	bencode.NewEncoder(c.Writer).Encode(resp)
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

// BuildPeerList constructs the list of peers to be returned to the client,
// handling LAN IP substitution if peers share the same public IP.
func BuildPeerList(clientRealIP string, excludePeerID string, dbPeers []*trackerModel.Peer) []map[string]interface{} {
	var bPeers []map[string]interface{}
	for _, p := range dbPeers {
		if p.PeerID == excludePeerID {
			continue
		}
		rawPeerID, err := hex.DecodeString(p.PeerID)
		if err != nil {
			continue
		}

		ipToReturn := p.IP
		// If both peers share the exact same public (real) IP, and the target peer has a valid LAN IP registered
		if p.LanIP != "" && p.IP == clientRealIP {
			ipToReturn = p.LanIP
		}

		bPeers = append(bPeers, map[string]interface{}{
			"peer id": string(rawPeerID),
			"ip":      ipToReturn,
			"port":    p.Port,
		})
	}
	if bPeers == nil {
		bPeers = []map[string]interface{}{}
	}
	return bPeers
}

// BuildCompactPeerList constructs the binary compact representation of peers.
func BuildCompactPeerList(clientRealIP string, excludePeerID string, dbPeers []*trackerModel.Peer) string {
	var buf []byte
	for _, p := range dbPeers {
		if p.PeerID == excludePeerID {
			continue
		}

		ipToReturn := p.IP
		if p.LanIP != "" && p.IP == clientRealIP {
			ipToReturn = p.LanIP
		}

		ipBytes := net.ParseIP(ipToReturn).To4()
		if ipBytes == nil {
			continue // Skip IPv6 in IPv4 compact response
		}

		buf = append(buf, ipBytes...)
		buf = append(buf, byte(p.Port>>8), byte(p.Port&0xFF))
	}
	return string(buf)
}
