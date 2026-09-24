package dao

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/GoldenSheep402/Hermes/mod/torrent/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

const (
	infoHashCacheTTL    = 10 * time.Minute
	infoHashNegativeTTL = 60 * time.Second
	infoHashCacheKeyPref = "tracker:infohash:"
)

func infoHashCacheKey(hash string) string {
	return infoHashCacheKeyPref + hash
}

// GetByHashCached resolves a torrent by info_hash via Redis, falling back to Postgres.
// Cached entries only carry ID and Size — enough for announce and bonus.
func (t *torrent) GetByHashCached(ctx context.Context, hash string) (*model.Torrent, error) {
	if hash == "" {
		return nil, status.Error(codes.NotFound, "Torrent not found")
	}

	if t.rds != nil {
		key := infoHashCacheKey(hash)
		values, err := t.rds.HGetAll(ctx, key).Result()
		if err == nil && len(values) > 0 {
			if values["missing"] == "1" {
				return nil, status.Error(codes.NotFound, "Torrent not found")
			}
			torrentID := values["torrent_id"]
			if torrentID == "" {
				return nil, status.Error(codes.NotFound, "Torrent not found")
			}
			size, _ := strconv.ParseInt(values["size"], 10, 64)
			return &model.Torrent{
				Model:    stdao.Model{ID: torrentID},
				InfoHash: hash,
				Size:     size,
				IsActive: true,
			}, nil
		}
	}

	torrent, err := t.GetByHash(ctx, hash)
	if err != nil {
		if status.Code(err) == codes.NotFound || errors.Is(err, gorm.ErrRecordNotFound) {
			t.cacheInfoHashMiss(ctx, hash)
		}
		return nil, err
	}
	t.cacheInfoHashHit(ctx, hash, torrent.ID, torrent.Size)
	return torrent, nil
}

func (t *torrent) cacheInfoHashHit(ctx context.Context, hash, torrentID string, size int64) {
	if t.rds == nil || hash == "" || torrentID == "" {
		return
	}
	key := infoHashCacheKey(hash)
	_ = t.rds.HSet(ctx, key, map[string]interface{}{
		"torrent_id": torrentID,
		"size":       strconv.FormatInt(size, 10),
		"missing":    "0",
	}).Err()
	_ = t.rds.Expire(ctx, key, infoHashCacheTTL).Err()
}

func (t *torrent) cacheInfoHashMiss(ctx context.Context, hash string) {
	if t.rds == nil || hash == "" {
		return
	}
	key := infoHashCacheKey(hash)
	_ = t.rds.HSet(ctx, key, map[string]interface{}{
		"missing": "1",
	}).Err()
	_ = t.rds.Expire(ctx, key, infoHashNegativeTTL).Err()
}

// InvalidateInfoHashCache removes a cached info_hash entry.
func (t *torrent) InvalidateInfoHashCache(ctx context.Context, hash string) {
	if t.rds == nil || hash == "" {
		return
	}
	_ = t.rds.Del(ctx, infoHashCacheKey(hash)).Err()
}

// CacheInfoHash writes/refreshes the info_hash cache entry.
func (t *torrent) CacheInfoHash(ctx context.Context, hash, torrentID string, size int64) {
	t.cacheInfoHashHit(ctx, hash, torrentID, size)
}

func (t *torrent) SetRedis(rds *redis.Client) {
	t.rds = rds
}
