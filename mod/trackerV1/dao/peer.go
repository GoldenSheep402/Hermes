package dao

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/GoldenSheep402/Hermes/mod/trackerV1/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

/*
Redis
│
├── Hash Key: "torrent:abc123"   (种子 A)
│   ├── Field: "peer1"           → Value: "192.168.1.1:6881:1609459200"
│   └── Field: "peer2"           → Value: "192.168.1.2:6882:1609459260"
│
└── Hash Key: "torrent:def456"   (种子 B)
    └── Field: "peer3"           → Value: "192.168.1.3:6883:1609459320"

Table: torrent:{HASH}
-----------------------------------------------------
|   Field（PeerID）   |            Value             |
-----------------------------------------------------
|      peerID1        | "IP1:Port1:LastSeen1"       |
|      peerID2        | "IP2:Port2:LastSeen2"       |
|      peerID3        | "IP3:Port3:LastSeen3"       |
|         ...         |            ...              |
-----------------------------------------------------
*/

type peer struct {
	stdao.Std[*model.Peer]
	rds *redis.Client
}

func (p *peer) Init(db *gorm.DB, rds *redis.Client) error {
	err := p.Std.Init(db)
	if err != nil {
		return err
	}
	p.rds = rds
	return nil
}

// AddPeer adds a peer to the torrent.
func (p *peer) AddPeer(ctx context.Context, infoHash string, peerData *model.Peer) error {
	key := "torrent:" + infoHash
	field := peerData.PeerID
	value := peerData.IP + ":" + strconv.Itoa(peerData.Port) + ":" + strconv.FormatInt(time.Now().Unix(), 10)
	err := p.rds.HSet(ctx, key, field, value).Err()
	if err != nil {
		return err
	}
	return p.rds.Expire(ctx, key, time.Hour).Err()
}

// GetPeers returns a list of peers for the torrent.
func (p *peer) GetPeers(ctx context.Context, infoHash string, numWant int) ([]*model.Peer, error) {
	key := "torrent:" + infoHash
	peersData, err := p.rds.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	peers := make([]*model.Peer, 0, len(peersData))
	now := time.Now().Unix()
	for peerID, data := range peersData {
		parts := strings.Split(data, ":")
		if len(parts) != 3 {
			continue
		}
		ip := parts[0]
		port, _ := strconv.Atoi(parts[1])
		lastSeen, _ := strconv.ParseInt(parts[2], 10, 64)
		if now-lastSeen > 3600 {
			p.rds.HDel(ctx, key, peerID)
			continue
		}
		peers = append(peers, &model.Peer{
			PeerID:   peerID,
			IP:       ip,
			Port:     port,
			LastSeen: time.Unix(lastSeen, 0),
		})
		if len(peers) >= numWant {
			break
		}
	}
	return peers, nil
}

// RemovePeer removes a peer from the torrent.
func (p *peer) RemovePeer(ctx context.Context, infoHash string, peerID string) error {
	key := "torrent:" + infoHash
	return p.rds.HDel(ctx, key, peerID).Err()
}
