package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/GoldenSheep402/Hermes/mod/tracker/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type peer struct {
	rds *redis.Client
}

func (p *peer) Init(db *gorm.DB, rds *redis.Client) error {
	p.rds = rds
	return nil
}

// Peer key: tracker:peer:{torrent_id}:{peer_id}
func peerKey(torrentID, peerID string) string {
	return fmt.Sprintf("tracker:peer:%s:%s", torrentID, peerID)
}

// Torrent peers set: tracker:torrent_peers:{torrent_id}
func torrentPeersKey(torrentID string) string {
	return fmt.Sprintf("tracker:torrent_peers:%s", torrentID)
}

func (p *peer) Upsert(ctx context.Context, peerData *model.Peer) error {
	key := peerKey(peerData.TorrentID, peerData.PeerID)
	data, err := json.Marshal(peerData)
	if err != nil {
		return err
	}

	// Save peer data with a TTL (e.g., 40 minutes, slightly longer than the announce interval)
	ttl := 40 * time.Minute

	pipe := p.rds.Pipeline()
	pipe.Set(ctx, key, data, ttl)
	// Add to torrent's peer set, score is last action timestamp to easily remove inactive peers
	pipe.ZAdd(ctx, torrentPeersKey(peerData.TorrentID), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: peerData.PeerID,
	})
	pipe.Expire(ctx, torrentPeersKey(peerData.TorrentID), ttl)
	_, err = pipe.Exec(ctx)
	return err
}

func (p *peer) GetPeersForTorrent(ctx context.Context, torrentID string, limit int) ([]*model.Peer, error) {
	// First cleanup old peers (e.g., inactive for > 40 minutes)
	cutoff := float64(time.Now().Add(-40 * time.Minute).Unix())
	p.rds.ZRemRangeByScore(ctx, torrentPeersKey(torrentID), "-inf", fmt.Sprintf("%f", cutoff))

	// Get active peer IDs
	var count int64 = int64(limit)
	if limit <= 0 {
		count = 50
	}
	peerIDs, err := p.rds.ZRevRange(ctx, torrentPeersKey(torrentID), 0, count-1).Result()
	if err != nil {
		return nil, err
	}

	var peers []*model.Peer
	if len(peerIDs) == 0 {
		return peers, nil
	}

	// Fetch peer details
	keys := make([]string, len(peerIDs))
	for i, id := range peerIDs {
		keys[i] = peerKey(torrentID, id)
	}

	res, err := p.rds.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	for _, val := range res {
		if val == nil {
			continue
		}
		strVal, ok := val.(string)
		if !ok {
			continue
		}
		var peer model.Peer
		if err := json.Unmarshal([]byte(strVal), &peer); err == nil {
			peers = append(peers, &peer)
		}
	}

	return peers, nil
}

func (p *peer) DeletePeer(ctx context.Context, torrentID, peerID string) error {
	pipe := p.rds.Pipeline()
	pipe.Del(ctx, peerKey(torrentID, peerID))
	pipe.ZRem(ctx, torrentPeersKey(torrentID), peerID)
	_, err := pipe.Exec(ctx)
	return err
}
