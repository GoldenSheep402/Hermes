// Package dao implements tracker Redis traffic aggregation and flush to PostgreSQL.
//
// Multi-instance deployment: each Hermes/tracker process should use one authoritative Redis
// for a given logical site (optionally via RedisKeyPrefix). Flushing writes absolute totals
// from that Redis into shared user/torrent rows; multiple independent Redis instances writing
// the same user_id/torrent_id will overwrite each other's totals in DB — avoid sharing one
// database across isolated tracker pools without a merge strategy.
package dao

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	torrentModel "github.com/GoldenSheep402/Hermes/mod/torrent/model"
	trackerModel "github.com/GoldenSheep402/Hermes/mod/tracker/model"
	trafficModel "github.com/GoldenSheep402/Hermes/mod/traffic/model"
	userModel "github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	trafficVersionField = "_version"

	realtimeWindowSeconds     int64 = 60
	realtimeBucketSpanSeconds int64 = 5
	realtimeBucketTTL               = 10 * time.Minute
)

func trafficDirtyUsersKey() string { return prefixed("tracker:traffic:dirty:users") }
func trafficDirtyTorrentsKey() string {
	return prefixed("tracker:traffic:dirty:torrents")
}
func trafficDirtyPairsKey() string { return prefixed("tracker:traffic:dirty:pairs") }
func siteTrafficTotalsKey() string { return prefixed("tracker:traffic:site:totals") }
func siteRealtimeUploadKey() string { return prefixed("tracker:traffic:realtime:site:upload") }
func siteRealtimeDownloadKey() string {
	return prefixed("tracker:traffic:realtime:site:download")
}

var removeDirtyIfVersionUnchanged = redis.NewScript(`
local current = redis.call('HGET', KEYS[1], ARGV[2])
if not current then
  redis.call('SREM', KEYS[2], ARGV[1])
  return 1
end
if tonumber(current) == tonumber(ARGV[3]) then
  redis.call('SREM', KEYS[2], ARGV[1])
  return 1
end
return 0
`)

type traffic struct {
	db  *gorm.DB
	rds *redis.Client
}

func (t *traffic) Init(db *gorm.DB, rds *redis.Client) error {
	t.db = db
	t.rds = rds
	return nil
}

func (t *traffic) RecordDelta(
	ctx context.Context,
	userID, torrentID string,
	uploadDelta, downloadDelta int64,
	isActive, isFinished bool,
	actionAt time.Time,
) error {
	if t.rds == nil || t.db == nil {
		return nil
	}
	if userID == "" || torrentID == "" {
		return nil
	}

	uploadDelta = nonNegative(uploadDelta)
	downloadDelta = nonNegative(downloadDelta)
	if actionAt.IsZero() {
		actionAt = time.Now()
	}

	if err := t.ensureUserTotalsInitialized(ctx, userID); err != nil {
		return err
	}
	if err := t.ensureSiteTotalsInitialized(ctx); err != nil {
		return err
	}
	if err := t.ensureTorrentTotalsInitialized(ctx, torrentID); err != nil {
		return err
	}
	if err := t.ensurePairTotalsInitialized(ctx, userID, torrentID); err != nil {
		return err
	}

	userKey := userTrafficKey(userID)
	torrentKey := torrentTrafficKey(torrentID)
	pairKey := pairTrafficKey(userID, torrentID)
	pairMember := pairTrafficMember(userID, torrentID)

	pipe := t.rds.TxPipeline()

	if uploadDelta > 0 || downloadDelta > 0 {
		realtimeBucket := alignRealtimeBucket(actionAt.Unix())
		realtimeField := strconv.FormatInt(realtimeBucket, 10)
		expiredField := strconv.FormatInt(realtimeBucket-realtimeWindowSeconds-realtimeBucketSpanSeconds, 10)
		realtimeUploadKey := userRealtimeUploadKey(userID)
		realtimeDownloadKey := userRealtimeDownloadKey(userID)

		if uploadDelta > 0 {
			pipe.HIncrBy(ctx, userKey, "real_upload", uploadDelta)
			pipe.HIncrBy(ctx, siteTrafficTotalsKey(), "real_upload", uploadDelta)
			pipe.HIncrBy(ctx, torrentKey, "total_upload", uploadDelta)
			pipe.HIncrBy(ctx, pairKey, "uploaded", uploadDelta)
			pipe.HIncrBy(ctx, realtimeUploadKey, realtimeField, uploadDelta)
			pipe.HIncrBy(ctx, siteRealtimeUploadKey(), realtimeField, uploadDelta)
			pipe.HDel(ctx, realtimeUploadKey, expiredField)
			pipe.HDel(ctx, siteRealtimeUploadKey(), expiredField)
			pipe.Expire(ctx, realtimeUploadKey, realtimeBucketTTL)
			pipe.Expire(ctx, siteRealtimeUploadKey(), realtimeBucketTTL)
		}
		if downloadDelta > 0 {
			pipe.HIncrBy(ctx, userKey, "real_download", downloadDelta)
			pipe.HIncrBy(ctx, siteTrafficTotalsKey(), "real_download", downloadDelta)
			pipe.HIncrBy(ctx, torrentKey, "total_download", downloadDelta)
			pipe.HIncrBy(ctx, pairKey, "downloaded", downloadDelta)
			pipe.HIncrBy(ctx, realtimeDownloadKey, realtimeField, downloadDelta)
			pipe.HIncrBy(ctx, siteRealtimeDownloadKey(), realtimeField, downloadDelta)
			pipe.HDel(ctx, realtimeDownloadKey, expiredField)
			pipe.HDel(ctx, siteRealtimeDownloadKey(), expiredField)
			pipe.Expire(ctx, realtimeDownloadKey, realtimeBucketTTL)
			pipe.Expire(ctx, siteRealtimeDownloadKey(), realtimeBucketTTL)
		}
		pipe.HIncrBy(ctx, userKey, trafficVersionField, 1)
		pipe.HIncrBy(ctx, torrentKey, trafficVersionField, 1)
		pipe.SAdd(ctx, trafficDirtyUsersKey(), userID)
		pipe.SAdd(ctx, trafficDirtyTorrentsKey(), torrentID)
	}

	pipe.HSet(ctx, pairKey,
		"is_active", boolToInt(isActive),
		"is_finished", boolToInt(isFinished),
		"last_action", actionAt.Unix(),
	)
	pipe.HIncrBy(ctx, pairKey, trafficVersionField, 1)
	pipe.SAdd(ctx, trafficDirtyPairsKey(), pairMember)

	_, err := pipe.Exec(ctx)
	return err
}

func (t *traffic) GetUserRealtimeRate(ctx context.Context, userID string) (int64, int64, error) {
	if t.rds == nil || userID == "" {
		return 0, 0, nil
	}
	return t.getRealtimeRateByKeys(ctx, userRealtimeUploadKey(userID), userRealtimeDownloadKey(userID))
}

func (t *traffic) GetUserTotals(ctx context.Context, userID string) (int64, int64, bool, error) {
	if t.rds == nil || userID == "" {
		return 0, 0, false, nil
	}

	key := userTrafficKey(userID)
	exists, err := t.rds.Exists(ctx, key).Result()
	if err != nil {
		return 0, 0, false, err
	}
	if exists <= 0 {
		return 0, 0, false, nil
	}

	values, err := t.rds.HMGet(ctx, key, "real_upload", "real_download").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, 0, false, err
	}
	if len(values) < 2 {
		return 0, 0, true, nil
	}

	return parseAnyInt64(values[0]), parseAnyInt64(values[1]), true, nil
}

func (t *traffic) GetUserSnapshot(ctx context.Context, userID string) (realUpload, realDownload, seedTime int64, found bool, err error) {
	if t.rds == nil || userID == "" {
		return 0, 0, 0, false, nil
	}

	if err := t.ensureUserTotalsInitialized(ctx, userID); err != nil {
		return 0, 0, 0, false, err
	}

	key := userTrafficKey(userID)
	values, err := t.rds.HMGet(ctx, key, "real_upload", "real_download", "seed_time", "credited_upload", "credited_download").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, 0, 0, false, err
	}
	if len(values) < 3 {
		return 0, 0, 0, false, nil
	}

	return parseAnyInt64(values[0]), parseAnyInt64(values[1]), parseAnyInt64(values[2]), true, nil
}

// GetUserCreditedSnapshot returns credited upload/download plus seed time.
func (t *traffic) GetUserCreditedSnapshot(ctx context.Context, userID string) (creditedUpload, creditedDownload, seedTime int64, found bool, err error) {
	if t.rds == nil || userID == "" {
		return 0, 0, 0, false, nil
	}
	if err := t.ensureUserTotalsInitialized(ctx, userID); err != nil {
		return 0, 0, 0, false, err
	}
	key := userTrafficKey(userID)
	values, err := t.rds.HMGet(ctx, key, "credited_upload", "credited_download", "seed_time", "real_upload", "real_download").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, 0, 0, false, err
	}
	if len(values) < 5 {
		return 0, 0, 0, false, nil
	}
	creditedUpload = parseAnyInt64(values[0])
	creditedDownload = parseAnyInt64(values[1])
	seedTime = parseAnyInt64(values[2])
	if values[0] == nil {
		creditedUpload = parseAnyInt64(values[3])
	}
	if values[1] == nil {
		creditedDownload = parseAnyInt64(values[4])
	}
	return creditedUpload, creditedDownload, seedTime, true, nil
}

func (t *traffic) CountUserActive(ctx context.Context, userID string) (int64, int64, error) {
	if t.rds == nil || userID == "" {
		return 0, 0, nil
	}
	if err := cleanupActivityZSet(ctx, t.rds, userSeedingKey(userID)); err != nil {
		return 0, 0, err
	}
	if err := cleanupActivityZSet(ctx, t.rds, userDownloadingKey(userID)); err != nil {
		return 0, 0, err
	}

	pipe := t.rds.Pipeline()
	seedingCmd := pipe.ZCard(ctx, userSeedingKey(userID))
	downloadingCmd := pipe.ZCard(ctx, userDownloadingKey(userID))
	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return 0, 0, err
	}
	return seedingCmd.Val(), downloadingCmd.Val(), nil
}

func (t *traffic) GetTorrentSnapshot(ctx context.Context, torrentID string) (*trafficModel.TorrentStats, error) {
	if torrentID == "" {
		return &trafficModel.TorrentStats{}, nil
	}
	if t.rds == nil {
		return t.loadTorrentSnapshotFromDB(ctx, torrentID)
	}
	if err := t.ensureTorrentTotalsInitialized(ctx, torrentID); err != nil {
		return nil, err
	}
	if err := cleanupActivityZSet(ctx, t.rds, torrentSeederKey(torrentID)); err != nil {
		return nil, err
	}
	if err := cleanupActivityZSet(ctx, t.rds, torrentLeecherKey(torrentID)); err != nil {
		return nil, err
	}
	if err := cleanupActivityZSet(ctx, t.rds, torrentPeersKey(torrentID)); err != nil {
		return nil, err
	}

	pipe := t.rds.Pipeline()
	hashCmd := pipe.HMGet(ctx, torrentTrafficKey(torrentID), "total_upload", "total_download", "snatch_count")
	seedCmd := pipe.ZCard(ctx, torrentSeederKey(torrentID))
	leechCmd := pipe.ZCard(ctx, torrentLeecherKey(torrentID))
	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	values := hashCmd.Val()
	stats := &trafficModel.TorrentStats{
		TorrentID:     torrentID,
		SeedCount:     int(seedCmd.Val()),
		LeechCount:    int(leechCmd.Val()),
		SnatchCount:   int(parseIndexedAnyInt64(values, 2)),
		TotalUpload:   parseIndexedAnyInt64(values, 0),
		TotalDownload: parseIndexedAnyInt64(values, 1),
	}
	return stats, nil
}

func (t *traffic) GetSiteRealtimeRate(ctx context.Context) (int64, int64, error) {
	if t.rds == nil {
		return 0, 0, nil
	}
	return t.getRealtimeRateByKeys(ctx, siteRealtimeUploadKey(), siteRealtimeDownloadKey())
}

func (t *traffic) GetSiteTotals(ctx context.Context) (int64, int64, bool, error) {
	if t.rds == nil {
		return 0, 0, false, nil
	}
	if err := t.ensureSiteTotalsInitialized(ctx); err != nil {
		return 0, 0, false, err
	}

	values, err := t.rds.HMGet(ctx, siteTrafficTotalsKey(), "real_upload", "real_download").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, 0, false, err
	}
	if len(values) < 2 {
		return 0, 0, false, nil
	}
	return parseAnyInt64(values[0]), parseAnyInt64(values[1]), true, nil
}

func (t *traffic) getRealtimeRateByKeys(ctx context.Context, uploadKey, downloadKey string) (int64, int64, error) {
	pipe := t.rds.Pipeline()
	uploadCmd := pipe.HGetAll(ctx, uploadKey)
	downloadCmd := pipe.HGetAll(ctx, downloadKey)
	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return 0, 0, err
	}

	nowUnix := time.Now().Unix()
	uploadRate := computeRealtimeRate(uploadCmd.Val(), nowUnix)
	downloadRate := computeRealtimeRate(downloadCmd.Val(), nowUnix)
	return uploadRate, downloadRate, nil
}

// FlushPending flushes at most batchSize dirty keys from Redis to DB.
// It is safe to call repeatedly; values are written as absolute totals.
func (t *traffic) FlushPending(ctx context.Context, batchSize int) (int, error) {
	if t.rds == nil || t.db == nil {
		return 0, nil
	}
	if batchSize <= 0 {
		batchSize = 200
	}

	total := 0
	var firstErr error

	n, err := t.flushUsers(ctx, batchSize)
	total += n
	if err != nil {
		firstErr = err
	}

	n, err = t.flushTorrents(ctx, batchSize)
	total += n
	if err != nil && firstErr == nil {
		firstErr = err
	}

	n, err = t.flushPairs(ctx, batchSize)
	total += n
	if err != nil && firstErr == nil {
		firstErr = err
	}

	return total, firstErr
}

func (t *traffic) flushUsers(ctx context.Context, batchSize int) (int, error) {
	members, err := t.rds.SRandMemberN(ctx, trafficDirtyUsersKey(), int64(batchSize)).Result()
	if err != nil {
		return 0, err
	}

	flushed := 0
	var firstErr error
	for _, userID := range members {
		if userID == "" {
			continue
		}
		key := userTrafficKey(userID)
		values, err := t.rds.HGetAll(ctx, key).Result()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if len(values) == 0 {
			_ = t.rds.SRem(ctx, trafficDirtyUsersKey(), userID).Err()
			continue
		}

		realUpload := parseInt64(values["real_upload"])
		realDownload := parseInt64(values["real_download"])
		creditedUpload := parseInt64(values["credited_upload"])
		creditedDownload := parseInt64(values["credited_download"])
		if _, hasCredited := values["credited_upload"]; !hasCredited {
			creditedUpload = realUpload
		}
		if _, hasCredited := values["credited_download"]; !hasCredited {
			creditedDownload = realDownload
		}
		seedTime := parseInt64(values["seed_time"])
		version := parseInt64(values[trafficVersionField])

		err = t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var existingBonusUpload int64
			var existing trafficModel.UserTraffic
			if qErr := tx.Where("user_id = ?", userID).First(&existing).Error; qErr == nil {
				existingBonusUpload = existing.BonusUpload
			} else if !errors.Is(qErr, gorm.ErrRecordNotFound) {
				return qErr
			}

			ut := &trafficModel.UserTraffic{
				Model:         stdao.Model{ID: ulid.Make().String()},
				UserID:        userID,
				RealUpload:    realUpload,
				RealDownload:  realDownload,
				BonusUpload:   existingBonusUpload,
				BonusDownload: 0,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "user_id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"real_upload":   realUpload,
					"real_download": realDownload,
					"updated_at":    time.Now(),
				}),
			}).Create(ut).Error; err != nil {
				return err
			}

			return tx.Model(&userModel.User{}).
				Where("id = ?", userID).
				Updates(map[string]interface{}{
					"uploaded":   creditedUpload + existingBonusUpload,
					"downloaded": creditedDownload,
					"seed_time":  seedTime,
					"updated_at": time.Now(),
				}).Error
		})
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		if err := t.finalizeDirtyMember(ctx, key, trafficDirtyUsersKey(), userID, version); err != nil && firstErr == nil {
			firstErr = err
		}
		flushed++
	}

	return flushed, firstErr
}

func (t *traffic) flushTorrents(ctx context.Context, batchSize int) (int, error) {
	members, err := t.rds.SRandMemberN(ctx, trafficDirtyTorrentsKey(), int64(batchSize)).Result()
	if err != nil {
		return 0, err
	}

	flushed := 0
	var firstErr error
	for _, torrentID := range members {
		if torrentID == "" {
			continue
		}
		if err := t.ensureTorrentTotalsInitialized(ctx, torrentID); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		key := torrentTrafficKey(torrentID)
		values, err := t.rds.HGetAll(ctx, key).Result()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if len(values) == 0 {
			_ = t.rds.SRem(ctx, trafficDirtyTorrentsKey(), torrentID).Err()
			continue
		}

		totalUpload := parseInt64(values["total_upload"])
		totalDownload := parseInt64(values["total_download"])
		snatchCount := int(parseInt64(values["snatch_count"]))
		version := parseInt64(values[trafficVersionField])
		if err := cleanupActivityZSet(ctx, t.rds, torrentSeederKey(torrentID)); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if err := cleanupActivityZSet(ctx, t.rds, torrentLeecherKey(torrentID)); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		pipe := t.rds.Pipeline()
		seedCmd := pipe.ZCard(ctx, torrentSeederKey(torrentID))
		leechCmd := pipe.ZCard(ctx, torrentLeecherKey(torrentID))
		if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		seedCount := int(seedCmd.Val())
		leechCount := int(leechCmd.Val())

		ts := &trafficModel.TorrentStats{
			Model:         stdao.Model{ID: ulid.Make().String()},
			TorrentID:     torrentID,
			SeedCount:     seedCount,
			LeechCount:    leechCount,
			SnatchCount:   snatchCount,
			TotalUpload:   totalUpload,
			TotalDownload: totalDownload,
		}
		err = t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "torrent_id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"seed_count":     seedCount,
					"leech_count":    leechCount,
					"snatch_count":   snatchCount,
					"total_upload":   totalUpload,
					"total_download": totalDownload,
					"updated_at":     time.Now(),
				}),
			}).Create(ts).Error; err != nil {
				return err
			}

			return tx.Model(&torrentModel.Torrent{}).
				Where("id = ?", torrentID).
				Updates(map[string]interface{}{
					"seed_count":   seedCount,
					"leech_count":  leechCount,
					"snatch_count": snatchCount,
					"updated_at":   time.Now(),
				}).Error
		})
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		if err := t.finalizeDirtyMember(ctx, key, trafficDirtyTorrentsKey(), torrentID, version); err != nil && firstErr == nil {
			firstErr = err
		}
		flushed++
	}

	return flushed, firstErr
}

func (t *traffic) flushPairs(ctx context.Context, batchSize int) (int, error) {
	members, err := t.rds.SRandMemberN(ctx, trafficDirtyPairsKey(), int64(batchSize)).Result()
	if err != nil {
		return 0, err
	}

	flushed := 0
	var firstErr error
	for _, member := range members {
		userID, torrentID, ok := parsePairTrafficMember(member)
		if !ok {
			_ = t.rds.SRem(ctx, trafficDirtyPairsKey(), member).Err()
			continue
		}

		key := pairTrafficKey(userID, torrentID)
		values, err := t.rds.HGetAll(ctx, key).Result()
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if len(values) == 0 {
			_ = t.rds.SRem(ctx, trafficDirtyPairsKey(), member).Err()
			continue
		}

		uploaded := parseInt64(values["uploaded"])
		downloaded := parseInt64(values["downloaded"])
		seedTime := parseInt64(values["seed_time"])
		version := parseInt64(values[trafficVersionField])
		isActive := parseBool(values["is_active"])
		isFinished := parseBool(values["is_finished"])
		finishedAtUnix := parseInt64(values["finished_at"])
		lastActionUnix := parseInt64(values["last_action"])
		lastAction := time.Unix(lastActionUnix, 0)
		if lastActionUnix <= 0 {
			lastAction = time.Now()
		}
		var finishedAt *time.Time
		if finishedAtUnix > 0 {
			finished := time.Unix(finishedAtUnix, 0)
			finishedAt = &finished
		}

		err = t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			record := &trafficModel.TransferHistory{
				Model:      stdao.Model{ID: ulid.Make().String()},
				UserID:     userID,
				TorrentID:  torrentID,
				Uploaded:   uploaded,
				Downloaded: downloaded,
				SeedTime:   seedTime,
				IsActive:   isActive,
				IsFinished: isFinished,
				LastAction: lastAction,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "user_id"}, {Name: "torrent_id"}},
				DoUpdates: clause.Assignments(map[string]any{
					"uploaded":    uploaded,
					"downloaded":  downloaded,
					"seed_time":   seedTime,
					"is_active":   isActive,
					"is_finished": isFinished,
					"last_action": lastAction,
					"updated_at":  time.Now(),
				}),
			}).Create(record).Error; err != nil {
				return err
			}

			if !isFinished {
				return nil
			}

			snatch := &trackerModel.Snatch{
				Model:      stdao.Model{ID: ulid.Make().String()},
				TorrentID:  torrentID,
				UserID:     userID,
				Uploaded:   uploaded,
				Downloaded: downloaded,
				SeedTime:   seedTime,
				IsActive:   isActive,
				FinishedAt: finishedAt,
				LastAction: lastAction,
			}
			return tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "torrent_id"}, {Name: "user_id"}},
				DoUpdates: clause.Assignments(map[string]any{
					"uploaded":    uploaded,
					"downloaded":  downloaded,
					"seed_time":   seedTime,
					"is_active":   isActive,
					"finished_at": finishedAt,
					"last_action": lastAction,
					"updated_at":  time.Now(),
				}),
			}).Create(snatch).Error
		})
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		if err := t.finalizeDirtyMember(ctx, key, trafficDirtyPairsKey(), member, version); err != nil && firstErr == nil {
			firstErr = err
		}
		flushed++
	}

	return flushed, firstErr
}

func (t *traffic) finalizeDirtyMember(ctx context.Context, hashKey, setKey, member string, version int64) error {
	if t.rds == nil {
		return nil
	}
	_, err := removeDirtyIfVersionUnchanged.Run(
		ctx,
		t.rds,
		[]string{hashKey, setKey},
		member,
		trafficVersionField,
		strconv.FormatInt(version, 10),
	).Result()
	return err
}

func (t *traffic) ensureUserTotalsInitialized(ctx context.Context, userID string) error {
	key := userTrafficKey(userID)
	exists, err := t.rds.Exists(ctx, key).Result()
	if err != nil || exists > 0 {
		return err
	}

	var upload, download, seedTime int64

	var ut trafficModel.UserTraffic
	err = t.db.WithContext(ctx).Where("user_id = ?", userID).First(&ut).Error
	if err == nil {
		upload = ut.RealUpload
		download = ut.RealDownload
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	var user userModel.User
	uErr := t.db.WithContext(ctx).
		Select("uploaded", "downloaded", "seed_time").
		Where("id = ?", userID).
		First(&user).Error
	if uErr != nil && !errors.Is(uErr, gorm.ErrRecordNotFound) {
		return uErr
	}
	if uErr == nil {
		if upload == 0 {
			upload = user.Uploaded
		}
		if download == 0 {
			download = user.Downloaded
		}
		seedTime = user.SeedTime
	}

	pipe := t.rds.TxPipeline()
	pipe.HSetNX(ctx, key, "real_upload", upload)
	pipe.HSetNX(ctx, key, "real_download", download)
	pipe.HSetNX(ctx, key, "credited_upload", upload)
	pipe.HSetNX(ctx, key, "credited_download", download)
	pipe.HSetNX(ctx, key, "seed_time", seedTime)
	pipe.HSetNX(ctx, key, trafficVersionField, 0)
	_, err = pipe.Exec(ctx)
	return err
}

func (t *traffic) ensureSiteTotalsInitialized(ctx context.Context) error {
	exists, err := t.rds.Exists(ctx, siteTrafficTotalsKey()).Result()
	if err != nil || exists > 0 {
		return err
	}

	var siteAgg struct {
		RealUpload   int64 `gorm:"column:real_upload"`
		RealDownload int64 `gorm:"column:real_download"`
	}
	if err := t.db.WithContext(ctx).
		Model(&trafficModel.UserTraffic{}).
		Select("COALESCE(SUM(real_upload), 0) AS real_upload, COALESCE(SUM(real_download), 0) AS real_download").
		Scan(&siteAgg).Error; err != nil {
		return err
	}

	if siteAgg.RealUpload == 0 && siteAgg.RealDownload == 0 {
		var userAgg struct {
			Uploaded   int64 `gorm:"column:uploaded"`
			Downloaded int64 `gorm:"column:downloaded"`
		}
		if err := t.db.WithContext(ctx).
			Model(&userModel.User{}).
			Select("COALESCE(SUM(uploaded), 0) AS uploaded, COALESCE(SUM(downloaded), 0) AS downloaded").
			Scan(&userAgg).Error; err != nil {
			return err
		}
		siteAgg.RealUpload = userAgg.Uploaded
		siteAgg.RealDownload = userAgg.Downloaded
	}

	pipe := t.rds.TxPipeline()
	pipe.HSetNX(ctx, siteTrafficTotalsKey(), "real_upload", siteAgg.RealUpload)
	pipe.HSetNX(ctx, siteTrafficTotalsKey(), "real_download", siteAgg.RealDownload)
	_, err = pipe.Exec(ctx)
	return err
}

func (t *traffic) ensureTorrentTotalsInitialized(ctx context.Context, torrentID string) error {
	key := torrentTrafficKey(torrentID)
	exists, err := t.rds.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		hasSnatchCount, err := t.rds.HExists(ctx, key, "snatch_count").Result()
		if err != nil || hasSnatchCount {
			return err
		}
	}

	var totalUpload, totalDownload int64
	var snatchCount int64
	var ts trafficModel.TorrentStats
	if err := t.db.WithContext(ctx).Where("torrent_id = ?", torrentID).First(&ts).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	} else {
		totalUpload = ts.TotalUpload
		totalDownload = ts.TotalDownload
		snatchCount = int64(ts.SnatchCount)
	}

	var torrent torrentModel.Torrent
	torrentErr := t.db.WithContext(ctx).
		Select("snatch_count").
		Where("id = ?", torrentID).
		First(&torrent).Error
	if torrentErr != nil && !errors.Is(torrentErr, gorm.ErrRecordNotFound) {
		return torrentErr
	}
	if torrentErr == nil && int64(torrent.SnatchCount) > snatchCount {
		snatchCount = int64(torrent.SnatchCount)
	}

	pipe := t.rds.TxPipeline()
	pipe.HSetNX(ctx, key, "snatch_count", snatchCount)
	if exists <= 0 {
		pipe.HSetNX(ctx, key, "total_upload", totalUpload)
		pipe.HSetNX(ctx, key, "total_download", totalDownload)
		pipe.HSetNX(ctx, key, trafficVersionField, 0)
	}
	_, err = pipe.Exec(ctx)
	return err
}

func (t *traffic) ensurePairTotalsInitialized(ctx context.Context, userID, torrentID string) error {
	key := pairTrafficKey(userID, torrentID)
	exists, err := t.rds.Exists(ctx, key).Result()
	if err != nil || exists > 0 {
		return err
	}

	var uploaded, downloaded, seedTime int64
	var isActive, isFinished bool
	lastAction := time.Now()
	var finishedAtUnix int64

	var th trafficModel.TransferHistory
	err = t.db.WithContext(ctx).
		Where("user_id = ? AND torrent_id = ?", userID, torrentID).
		First(&th).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil {
		uploaded = th.Uploaded
		downloaded = th.Downloaded
		seedTime = th.SeedTime
		isActive = th.IsActive
		isFinished = th.IsFinished
		if !th.LastAction.IsZero() {
			lastAction = th.LastAction
		}
	}

	var sn trackerModel.Snatch
	snErr := t.db.WithContext(ctx).
		Where("torrent_id = ? AND user_id = ?", torrentID, userID).
		First(&sn).Error
	if snErr != nil && !errors.Is(snErr, gorm.ErrRecordNotFound) {
		return snErr
	}
	if snErr == nil {
		if sn.Uploaded > uploaded {
			uploaded = sn.Uploaded
		}
		if sn.Downloaded > downloaded {
			downloaded = sn.Downloaded
		}
		if sn.SeedTime > seedTime {
			seedTime = sn.SeedTime
		}
		if sn.IsActive && sn.LastAction.After(lastAction) {
			isActive = true
			lastAction = sn.LastAction
		}
		if sn.FinishedAt != nil {
			isFinished = true
			finishedAtUnix = sn.FinishedAt.Unix()
		}
		if sn.LastAction.After(lastAction) {
			lastAction = sn.LastAction
		}
	}

	if isActive && time.Since(lastAction) > peerTTL {
		isActive = false
	}

	pipe := t.rds.TxPipeline()
	pipe.HSetNX(ctx, key, "uploaded", uploaded)
	pipe.HSetNX(ctx, key, "downloaded", downloaded)
	pipe.HSetNX(ctx, key, "seed_time", seedTime)
	pipe.HSetNX(ctx, key, "is_active", boolToInt(isActive))
	pipe.HSetNX(ctx, key, "is_finished", boolToInt(isFinished))
	pipe.HSetNX(ctx, key, "finished_at", finishedAtUnix)
	pipe.HSetNX(ctx, key, "last_action", lastAction.Unix())
	pipe.HSetNX(ctx, key, trafficVersionField, 0)
	_, err = pipe.Exec(ctx)
	return err
}

func (t *traffic) loadTorrentSnapshotFromDB(ctx context.Context, torrentID string) (*trafficModel.TorrentStats, error) {
	stats := &trafficModel.TorrentStats{TorrentID: torrentID}

	var ts trafficModel.TorrentStats
	if err := t.db.WithContext(ctx).Where("torrent_id = ?", torrentID).First(&ts).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	} else {
		stats = &ts
	}

	var torrent torrentModel.Torrent
	if err := t.db.WithContext(ctx).Where("id = ?", torrentID).First(&torrent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return stats, nil
		}
		return nil, err
	}

	if stats.SeedCount == 0 {
		stats.SeedCount = torrent.SeedCount
	}
	if stats.LeechCount == 0 {
		stats.LeechCount = torrent.LeechCount
	}
	if stats.SnatchCount == 0 {
		stats.SnatchCount = torrent.SnatchCount
	}
	return stats, nil
}

func cleanupActivityZSet(ctx context.Context, rds *redis.Client, key string) error {
	if rds == nil || key == "" {
		return nil
	}
	cutoff := strconv.FormatInt(time.Now().Add(-peerTTL).Unix(), 10)
	return rds.ZRemRangeByScore(ctx, key, "-inf", cutoff).Err()
}

func userTrafficKey(userID string) string {
	return prefixed("tracker:traffic:user:" + userID)
}

func userRealtimeUploadKey(userID string) string {
	return prefixed("tracker:traffic:realtime:user:" + userID + ":upload")
}

func userRealtimeDownloadKey(userID string) string {
	return prefixed("tracker:traffic:realtime:user:" + userID + ":download")
}

func torrentTrafficKey(torrentID string) string {
	return prefixed("tracker:traffic:torrent:" + torrentID)
}

func torrentSeederKey(torrentID string) string {
	return prefixed("tracker:traffic:torrent:" + torrentID + ":seeders")
}

func torrentLeecherKey(torrentID string) string {
	return prefixed("tracker:traffic:torrent:" + torrentID + ":leechers")
}

func pairTrafficKey(userID, torrentID string) string {
	return prefixed("tracker:traffic:pair:" + userID + ":" + torrentID)
}

func userSeedingKey(userID string) string {
	return prefixed("tracker:traffic:user:" + userID + ":seeding")
}

func userDownloadingKey(userID string) string {
	return prefixed("tracker:traffic:user:" + userID + ":downloading")
}

func pairTrafficMember(userID, torrentID string) string {
	return userID + ":" + torrentID
}

func parsePairTrafficMember(member string) (string, string, bool) {
	parts := strings.SplitN(member, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func parseInt64(v string) int64 {
	if v == "" {
		return 0
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0
	}
	return nonNegative(n)
}

func parseAnyInt64(v interface{}) int64 {
	switch value := v.(type) {
	case nil:
		return 0
	case int64:
		return nonNegative(value)
	case int32:
		return nonNegative(int64(value))
	case int:
		return nonNegative(int64(value))
	case []byte:
		return parseInt64(string(value))
	case string:
		return parseInt64(value)
	default:
		return 0
	}
}

func parseIndexedAnyInt64(values []interface{}, idx int) int64 {
	if idx < 0 || idx >= len(values) {
		return 0
	}
	return parseAnyInt64(values[idx])
}

func parseBool(v string) bool {
	return v == "1" || strings.EqualFold(v, "true")
}

func boolToInt(v bool) int64 {
	if v {
		return 1
	}
	return 0
}

func nonNegative(v int64) int64 {
	if v < 0 {
		return 0
	}
	return v
}

func alignRealtimeBucket(unix int64) int64 {
	if unix <= 0 {
		return 0
	}
	return (unix / realtimeBucketSpanSeconds) * realtimeBucketSpanSeconds
}

func computeRealtimeRate(buckets map[string]string, nowUnix int64) int64 {
	if len(buckets) == 0 || nowUnix <= 0 {
		return 0
	}

	windowStart := nowUnix - realtimeWindowSeconds
	var sum int64
	oldest := nowUnix
	hasTraffic := false

	for ts, value := range buckets {
		bucketUnix := parseInt64(ts)
		if bucketUnix <= 0 {
			continue
		}
		if bucketUnix < windowStart || bucketUnix > nowUnix {
			continue
		}

		delta := parseInt64(value)
		if delta <= 0 {
			continue
		}

		sum += delta
		if !hasTraffic || bucketUnix < oldest {
			oldest = bucketUnix
		}
		hasTraffic = true
	}

	if !hasTraffic || sum <= 0 {
		return 0
	}

	coveredSeconds := nowUnix - oldest + realtimeBucketSpanSeconds
	if coveredSeconds < realtimeBucketSpanSeconds {
		coveredSeconds = realtimeBucketSpanSeconds
	}
	if coveredSeconds > realtimeWindowSeconds {
		coveredSeconds = realtimeWindowSeconds
	}
	if coveredSeconds <= 0 {
		coveredSeconds = realtimeWindowSeconds
	}

	return sum / coveredSeconds
}
