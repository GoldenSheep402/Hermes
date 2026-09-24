package dao

import (
	"context"
	"fmt"
	"hash/fnv"
	"strconv"
	"time"

	"github.com/GoldenSheep402/Hermes/mod/tracker/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type peer struct {
	rds *redis.Client
}

const peerTTL = 40 * time.Minute

// maxPeersRPCList caps peer enumeration from internal RPC when limit is unset or excessive.
const maxPeersRPCList = 5000

const (
	peerSampleRecentNumerator   = 3
	peerSampleRecentDenominator = 5
	peerSampleWindowMultiplier  = 2
	peerSampleRotationWindow    = 5 * time.Minute
)

func (p *peer) Init(db *gorm.DB, rds *redis.Client) error {
	p.rds = rds
	return nil
}

// Peer key: tracker:peer:{torrent_id}:{peer_id}
func peerKey(torrentID, peerID string) string {
	return prefixed(fmt.Sprintf("tracker:peer:v2:%s:%s", torrentID, peerID))
}

// Torrent peers set: tracker:torrent_peers:{torrent_id}
func torrentPeersKey(torrentID string) string {
	return prefixed(fmt.Sprintf("tracker:torrent_peers:%s", torrentID))
}

func announceDedupeKey(torrentID, peerID string, uploaded, downloaded, left int64) string {
	return prefixed(fmt.Sprintf(
		"tracker:announce:dedupe:%s:%s:%s:%s:%s",
		torrentID,
		peerID,
		strconv.FormatInt(uploaded, 10),
		strconv.FormatInt(downloaded, 10),
		strconv.FormatInt(left, 10),
	))
}

func (p *peer) MarkAnnounceUnique(
	ctx context.Context,
	torrentID, peerID string,
	uploaded, downloaded, left int64,
	ttl time.Duration,
) (bool, error) {
	if p.rds == nil {
		return true, nil
	}
	if ttl <= 0 {
		ttl = 3 * time.Second
	}
	ok, err := p.rds.SetNX(
		ctx,
		announceDedupeKey(torrentID, peerID, uploaded, downloaded, left),
		1,
		ttl,
	).Result()
	if err != nil {
		return true, err
	}
	return ok, nil
}

func (p *peer) Upsert(ctx context.Context, peerData *model.Peer) error {
	if p.rds == nil {
		return nil
	}
	key := peerKey(peerData.TorrentID, peerData.PeerID)

	pipe := p.rds.Pipeline()
	pipe.HSet(ctx, key, peerFields(peerData))
	pipe.Expire(ctx, key, peerTTL)
	// Add to torrent's peer set, score is last action timestamp to easily remove inactive peers
	pipe.ZAdd(ctx, torrentPeersKey(peerData.TorrentID), redis.Z{
		Score:  float64(peerLastActionUnix(peerData)),
		Member: peerData.PeerID,
	})
	pipe.Expire(ctx, torrentPeersKey(peerData.TorrentID), peerTTL)
	_, err := pipe.Exec(ctx)
	return err
}

func (p *peer) GetPeersForTorrent(ctx context.Context, torrentID string, limit int) ([]*model.Peer, error) {
	if p.rds == nil {
		return []*model.Peer{}, nil
	}
	if err := p.cleanupInactivePeers(ctx, torrentID); err != nil {
		return nil, err
	}
	var (
		peerIDs []string
		err     error
	)
	if limit <= 0 {
		limit = maxPeersRPCList
	}
	peerIDs, err = p.fetchRecentPeerIDs(ctx, torrentID, normalizePeerLimit(limit))
	if err != nil {
		return nil, err
	}
	return p.fetchPeersByIDs(ctx, torrentID, peerIDs)
}

func (p *peer) GetPeersForTorrentSample(
	ctx context.Context,
	torrentID string,
	limit int,
	sampleSeed string,
	sampleAt time.Time,
) ([]*model.Peer, error) {
	if p.rds == nil {
		return []*model.Peer{}, nil
	}
	if err := p.cleanupInactivePeers(ctx, torrentID); err != nil {
		return nil, err
	}

	count := normalizePeerLimit(limit)
	total, err := p.rds.ZCard(ctx, torrentPeersKey(torrentID)).Result()
	if err != nil {
		return nil, err
	}
	if total <= count || sampleSeed == "" {
		peerIDs, err := p.fetchRecentPeerIDs(ctx, torrentID, count)
		if err != nil {
			return nil, err
		}
		return p.fetchPeersByIDs(ctx, torrentID, peerIDs)
	}

	recentCount := count * peerSampleRecentNumerator / peerSampleRecentDenominator
	if recentCount <= 0 {
		recentCount = 1
	}
	if recentCount >= count {
		recentCount = count - 1
	}
	if recentCount <= 0 {
		recentCount = 1
	}

	recentIDs, err := p.fetchRecentPeerIDs(ctx, torrentID, recentCount)
	if err != nil {
		return nil, err
	}

	remaining := count - int64(len(recentIDs))
	if remaining <= 0 {
		return p.fetchPeersByIDs(ctx, torrentID, recentIDs)
	}

	olderStart := int64(len(recentIDs))
	olderCount := total - olderStart
	if olderCount <= 0 {
		return p.fetchPeersByIDs(ctx, torrentID, recentIDs)
	}

	windowSize := remaining * peerSampleWindowMultiplier
	if windowSize < remaining {
		windowSize = remaining
	}
	if windowSize > olderCount {
		windowSize = olderCount
	}

	offset := rotatedPeerOffset(sampleSeed, sampleAt, olderCount)
	olderIDs, err := p.fetchRotatedPeerIDs(ctx, torrentID, olderStart, olderCount, windowSize, offset)
	if err != nil {
		return nil, err
	}

	peerIDs := mergePeerIDs(recentIDs, olderIDs, int(count))
	return p.fetchPeersByIDs(ctx, torrentID, peerIDs)
}

func (p *peer) cleanupInactivePeers(ctx context.Context, torrentID string) error {
	cutoff := float64(time.Now().Add(-peerTTL).Unix())
	return p.rds.ZRemRangeByScore(ctx, torrentPeersKey(torrentID), "-inf", fmt.Sprintf("%f", cutoff)).Err()
}

func normalizePeerLimit(limit int) int64 {
	if limit <= 0 {
		return 50
	}
	if limit > maxPeersRPCList {
		return maxPeersRPCList
	}
	return int64(limit)
}

func (p *peer) fetchRecentPeerIDs(ctx context.Context, torrentID string, count int64) ([]string, error) {
	if count <= 0 {
		return []string{}, nil
	}
	return p.rds.ZRevRange(ctx, torrentPeersKey(torrentID), 0, count-1).Result()
}

func (p *peer) fetchAllPeerIDs(ctx context.Context, torrentID string) ([]string, error) {
	return p.rds.ZRevRange(ctx, torrentPeersKey(torrentID), 0, -1).Result()
}

func (p *peer) fetchRotatedPeerIDs(
	ctx context.Context,
	torrentID string,
	olderStart, olderCount, windowSize, offset int64,
) ([]string, error) {
	if windowSize <= 0 || olderCount <= 0 {
		return []string{}, nil
	}

	key := torrentPeersKey(torrentID)
	start := olderStart + offset
	end := start + windowSize - 1
	olderEnd := olderStart + olderCount - 1

	if end <= olderEnd {
		return p.rds.ZRevRange(ctx, key, start, end).Result()
	}

	firstPart, err := p.rds.ZRevRange(ctx, key, start, olderEnd).Result()
	if err != nil {
		return nil, err
	}
	remaining := windowSize - int64(len(firstPart))
	if remaining <= 0 {
		return firstPart, nil
	}

	secondPart, err := p.rds.ZRevRange(ctx, key, olderStart, olderStart+remaining-1).Result()
	if err != nil {
		return nil, err
	}
	return append(firstPart, secondPart...), nil
}

func mergePeerIDs(recentIDs, olderIDs []string, limit int) []string {
	if limit <= 0 {
		return []string{}
	}
	result := make([]string, 0, limit)
	seen := make(map[string]struct{}, limit)

	appendIDs := func(ids []string) {
		for _, id := range ids {
			if len(result) >= limit {
				return
			}
			if id == "" {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			result = append(result, id)
		}
	}

	appendIDs(recentIDs)
	appendIDs(olderIDs)
	return result
}

func rotatedPeerOffset(sampleSeed string, sampleAt time.Time, olderCount int64) int64 {
	if olderCount <= 1 {
		return 0
	}
	if sampleAt.IsZero() {
		sampleAt = time.Now()
	}

	bucket := sampleAt.UTC().Unix()
	windowSeconds := int64(peerSampleRotationWindow / time.Second)
	if windowSeconds > 0 {
		bucket = bucket / windowSeconds
	}

	h := fnv.New64a()
	_, _ = h.Write([]byte(sampleSeed))
	_, _ = h.Write([]byte(":"))
	_, _ = h.Write([]byte(strconv.FormatInt(bucket, 10)))
	return int64(h.Sum64() % uint64(olderCount))
}

func (p *peer) fetchPeersByIDs(ctx context.Context, torrentID string, peerIDs []string) ([]*model.Peer, error) {
	var peers []*model.Peer
	if len(peerIDs) == 0 {
		return peers, nil
	}

	// Fetch peer details in a single pipeline to keep response cost predictable.
	keys := make([]string, len(peerIDs))
	for i, id := range peerIDs {
		keys[i] = peerKey(torrentID, id)
	}

	pipe := p.rds.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(keys))
	for i, key := range keys {
		cmds[i] = pipe.HGetAll(ctx, key)
	}
	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return nil, err
	}

	stalePeerIDs := make([]string, 0)
	for idx, cmd := range cmds {
		values := cmd.Val()
		if len(values) == 0 {
			if idx < len(peerIDs) {
				stalePeerIDs = append(stalePeerIDs, peerIDs[idx])
			}
			continue
		}
		peer, err := peerFromFields(values)
		if err != nil {
			if idx < len(peerIDs) {
				stalePeerIDs = append(stalePeerIDs, peerIDs[idx])
			}
			continue
		}
		if peer.TorrentID == "" {
			peer.TorrentID = torrentID
		}
		peers = append(peers, peer)
	}

	if len(stalePeerIDs) > 0 {
		members := make([]interface{}, 0, len(stalePeerIDs))
		for _, id := range stalePeerIDs {
			members = append(members, id)
		}
		p.rds.ZRem(ctx, torrentPeersKey(torrentID), members...)
	}

	return peers, nil
}

func (p *peer) Get(ctx context.Context, torrentID, peerID string) (*model.Peer, error) {
	if torrentID == "" || peerID == "" {
		return nil, nil
	}
	if p.rds == nil {
		return nil, nil
	}

	values, err := p.rds.HGetAll(ctx, peerKey(torrentID, peerID)).Result()
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, nil
	}
	peerData, err := peerFromFields(values)
	if err != nil {
		return nil, err
	}
	if peerData.TorrentID == "" {
		peerData.TorrentID = torrentID
	}
	return peerData, nil
}

func (p *peer) DeletePeer(ctx context.Context, torrentID, peerID string) error {
	if p.rds == nil {
		return nil
	}
	pipe := p.rds.Pipeline()
	pipe.Del(ctx, peerKey(torrentID, peerID))
	pipe.ZRem(ctx, torrentPeersKey(torrentID), peerID)
	_, err := pipe.Exec(ctx)
	return err
}

func peerFields(peerData *model.Peer) map[string]interface{} {
	return map[string]interface{}{
		"id":          peerData.ID,
		"torrent_id":  peerData.TorrentID,
		"user_id":     peerData.UserID,
		"peer_id":     peerData.PeerID,
		"ip":          peerData.IP,
		"lan_ip":      peerData.LanIP,
		"port":        peerData.Port,
		"uploaded":    nonNegative(peerData.Uploaded),
		"downloaded":  nonNegative(peerData.Downloaded),
		"left":        nonNegative(peerData.Left),
		"agent":       peerData.Agent,
		"is_seeder":   boolToInt(peerData.IsSeeder),
		"started_at":  unixOrZero(peerData.StartedAt),
		"last_action": unixOrZero(peerData.LastAction),
	}
}

func peerFromFields(values map[string]string) (*model.Peer, error) {
	if len(values) == 0 {
		return nil, nil
	}

	peerData := &model.Peer{}
	peerData.ID = values["id"]
	peerData.TorrentID = values["torrent_id"]
	peerData.UserID = values["user_id"]
	peerData.PeerID = values["peer_id"]
	peerData.IP = values["ip"]
	peerData.LanIP = values["lan_ip"]
	peerData.Port = int(parseInt64(values["port"]))
	peerData.Uploaded = parseInt64(values["uploaded"])
	peerData.Downloaded = parseInt64(values["downloaded"])
	peerData.Left = parseInt64(values["left"])
	peerData.Agent = values["agent"]
	peerData.IsSeeder = parseBool(values["is_seeder"])
	peerData.StartedAt = unixTime(values["started_at"])
	peerData.LastAction = unixTime(values["last_action"])

	if peerData.PeerID == "" {
		return nil, fmt.Errorf("missing peer_id")
	}
	return peerData, nil
}

func peerLastActionUnix(peerData *model.Peer) int64 {
	if peerData == nil {
		return time.Now().Unix()
	}
	return unixOrNow(peerData.LastAction)
}

func unixOrZero(ts time.Time) int64 {
	if ts.IsZero() {
		return 0
	}
	return ts.Unix()
}

func unixOrNow(ts time.Time) int64 {
	if ts.IsZero() {
		return time.Now().Unix()
	}
	return ts.Unix()
}

func unixTime(raw string) time.Time {
	sec := parseInt64(raw)
	if sec <= 0 {
		return time.Time{}
	}
	return time.Unix(sec, 0)
}
