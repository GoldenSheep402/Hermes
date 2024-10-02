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
│   ├── Field: "UID:peer1:192.168.1.1:6881"  → Value: "1609459200"
│   └── Field: "UID:peer2:192.168.1.2:6882"  → Value: "1609459260"
│
└── Hash Key: "torrent:def456"   (种子 B)
    └── Field: "UID:peer3:192.168.1.3:6883"  → Value: "1609459320"

Table: torrent:{ID}
-----------------------------------------------------
|   Field（PeerID）   |            Value             |
-----------------------------------------------------
| UID:peerID1:IP:PORT |           "LastSeen1"       |
|      peerID2        |           "LastSeen2"       |
|      peerID3        |           "LastSeen3"       |
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
func (p *peer) AddPeer(ctx context.Context, torrentID string, peerData *model.Peer, uid string) error {
	key := "Torrent:" + torrentID
	field := uid + ":" + peerData.PeerID + ":" + peerData.IP + ":" + strconv.Itoa(peerData.Port)
	value := strconv.FormatInt(time.Now().Unix(), 10)
	err := p.rds.HSet(ctx, key, field, value).Err()
	if err != nil {
		return err
	}
	return p.rds.Expire(ctx, key, 24*time.Hour).Err()
}

// GetPeers returns a list of peers for the torrent.
func (p *peer) GetPeers(ctx context.Context, torrentID string, numWant int) ([]*model.Peer, error) {
	key := "Torrent:" + torrentID
	peersData, err := p.rds.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	peers := make([]*model.Peer, 0, numWant)
	now := time.Now().Unix()

	for field, lastSeenStr := range peersData {
		parts := strings.Split(field, ":")
		if len(parts) != 4 {
			continue
		}
		peerID := parts[1]
		ip := parts[2]
		port, _ := strconv.Atoi(parts[3])
		lastSeen, _ := strconv.ParseInt(lastSeenStr, 10, 64)

		if now-lastSeen > 3600 {
			p.rds.HDel(ctx, key, field)
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

	err = p.rds.Expire(ctx, key, 24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	return peers, nil
}

// RemovePeer removes a peer from the torrent.
func (p *peer) RemovePeer(ctx context.Context, torrentID string, peerID string, uid string) error {
	key := "Torrent:" + torrentID

	peersData, err := p.rds.HGetAll(ctx, key).Result()
	if err != nil {
		return err
	}

	for field := range peersData {
		parts := strings.Split(field, ":")
		if len(parts) != 4 {
			continue
		}
		storedUID := parts[0]
		storedPeerID := parts[1]
		if storedUID == uid && storedPeerID == peerID {
			err := p.rds.HDel(ctx, key, field).Err()
			if err != nil {
				return err
			}
			break
		}
	}

	return p.rds.Expire(ctx, key, 24*time.Hour).Err()
}
