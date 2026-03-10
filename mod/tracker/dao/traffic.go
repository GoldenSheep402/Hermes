package dao

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	trafficModel "github.com/GoldenSheep402/Hermes/mod/traffic/model"
	userModel "github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	userTrafficDirtySetKey    = "tracker:traffic:dirty:users"
	torrentTrafficDirtySetKey = "tracker:traffic:dirty:torrents"
	pairTrafficDirtySetKey    = "tracker:traffic:dirty:pairs"

	trafficVersionField = "_version"

	realtimeWindowSeconds     int64 = 60
	realtimeBucketSpanSeconds int64 = 5
	realtimeBucketTTL               = 10 * time.Minute
)

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
			pipe.HIncrBy(ctx, torrentKey, "total_upload", uploadDelta)
			pipe.HIncrBy(ctx, pairKey, "uploaded", uploadDelta)
			pipe.HIncrBy(ctx, realtimeUploadKey, realtimeField, uploadDelta)
			pipe.HDel(ctx, realtimeUploadKey, expiredField)
			pipe.Expire(ctx, realtimeUploadKey, realtimeBucketTTL)
		}
		if downloadDelta > 0 {
			pipe.HIncrBy(ctx, userKey, "real_download", downloadDelta)
			pipe.HIncrBy(ctx, torrentKey, "total_download", downloadDelta)
			pipe.HIncrBy(ctx, pairKey, "downloaded", downloadDelta)
			pipe.HIncrBy(ctx, realtimeDownloadKey, realtimeField, downloadDelta)
			pipe.HDel(ctx, realtimeDownloadKey, expiredField)
			pipe.Expire(ctx, realtimeDownloadKey, realtimeBucketTTL)
		}
		pipe.HIncrBy(ctx, userKey, trafficVersionField, 1)
		pipe.HIncrBy(ctx, torrentKey, trafficVersionField, 1)
		pipe.SAdd(ctx, userTrafficDirtySetKey, userID)
		pipe.SAdd(ctx, torrentTrafficDirtySetKey, torrentID)
	}

	pipe.HSet(ctx, pairKey,
		"is_active", boolToInt(isActive),
		"is_finished", boolToInt(isFinished),
		"last_action", actionAt.Unix(),
	)
	pipe.HIncrBy(ctx, pairKey, trafficVersionField, 1)
	pipe.SAdd(ctx, pairTrafficDirtySetKey, pairMember)

	_, err := pipe.Exec(ctx)
	return err
}

func (t *traffic) GetUserRealtimeRate(ctx context.Context, userID string) (int64, int64, error) {
	if t.rds == nil || userID == "" {
		return 0, 0, nil
	}

	uploadKey := userRealtimeUploadKey(userID)
	downloadKey := userRealtimeDownloadKey(userID)

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
	if err != nil && firstErr == nil {
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
	members, err := t.rds.SRandMemberN(ctx, userTrafficDirtySetKey, int64(batchSize)).Result()
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
			_ = t.rds.SRem(ctx, userTrafficDirtySetKey, userID).Err()
			continue
		}

		realUpload := parseInt64(values["real_upload"])
		realDownload := parseInt64(values["real_download"])
		version := parseInt64(values[trafficVersionField])

		err = t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			ut := &trafficModel.UserTraffic{
				Model:         stdao.Model{ID: ulid.Make().String()},
				UserID:        userID,
				RealUpload:    realUpload,
				RealDownload:  realDownload,
				BonusUpload:   0,
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
					"uploaded":   realUpload,
					"downloaded": realDownload,
					"updated_at": time.Now(),
				}).Error
		})
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		if err := t.finalizeDirtyMember(ctx, key, userTrafficDirtySetKey, userID, version); err != nil && firstErr == nil {
			firstErr = err
		}
		flushed++
	}

	return flushed, firstErr
}

func (t *traffic) flushTorrents(ctx context.Context, batchSize int) (int, error) {
	members, err := t.rds.SRandMemberN(ctx, torrentTrafficDirtySetKey, int64(batchSize)).Result()
	if err != nil {
		return 0, err
	}

	flushed := 0
	var firstErr error
	for _, torrentID := range members {
		if torrentID == "" {
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
			_ = t.rds.SRem(ctx, torrentTrafficDirtySetKey, torrentID).Err()
			continue
		}

		totalUpload := parseInt64(values["total_upload"])
		totalDownload := parseInt64(values["total_download"])
		version := parseInt64(values[trafficVersionField])

		ts := &trafficModel.TorrentStats{
			Model:         stdao.Model{ID: ulid.Make().String()},
			TorrentID:     torrentID,
			TotalUpload:   totalUpload,
			TotalDownload: totalDownload,
		}
		err = t.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "torrent_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"total_upload":   totalUpload,
				"total_download": totalDownload,
				"updated_at":     time.Now(),
			}),
		}).Create(ts).Error
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		if err := t.finalizeDirtyMember(ctx, key, torrentTrafficDirtySetKey, torrentID, version); err != nil && firstErr == nil {
			firstErr = err
		}
		flushed++
	}

	return flushed, firstErr
}

func (t *traffic) flushPairs(ctx context.Context, batchSize int) (int, error) {
	members, err := t.rds.SRandMemberN(ctx, pairTrafficDirtySetKey, int64(batchSize)).Result()
	if err != nil {
		return 0, err
	}

	flushed := 0
	var firstErr error
	for _, member := range members {
		userID, torrentID, ok := parsePairTrafficMember(member)
		if !ok {
			_ = t.rds.SRem(ctx, pairTrafficDirtySetKey, member).Err()
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
			_ = t.rds.SRem(ctx, pairTrafficDirtySetKey, member).Err()
			continue
		}

		uploaded := parseInt64(values["uploaded"])
		downloaded := parseInt64(values["downloaded"])
		version := parseInt64(values[trafficVersionField])
		isActive := parseBool(values["is_active"])
		isFinished := parseBool(values["is_finished"])
		lastActionUnix := parseInt64(values["last_action"])
		lastAction := time.Unix(lastActionUnix, 0)
		if lastActionUnix <= 0 {
			lastAction = time.Now()
		}

		err = t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var existing trafficModel.TransferHistory
			err := tx.Where("user_id = ? AND torrent_id = ?", userID, torrentID).
				Order("created_at ASC").
				First(&existing).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if errors.Is(err, gorm.ErrRecordNotFound) {
				record := &trafficModel.TransferHistory{
					Model:      stdao.Model{ID: ulid.Make().String()},
					UserID:     userID,
					TorrentID:  torrentID,
					Uploaded:   uploaded,
					Downloaded: downloaded,
					IsActive:   isActive,
					IsFinished: isFinished,
					LastAction: lastAction,
				}
				return tx.Create(record).Error
			}

			return tx.Model(&trafficModel.TransferHistory{}).
				Where("id = ?", existing.ID).
				Updates(map[string]interface{}{
					"uploaded":    uploaded,
					"downloaded":  downloaded,
					"is_active":   isActive,
					"is_finished": isFinished,
					"last_action": lastAction,
					"updated_at":  time.Now(),
				}).Error
		})
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		if err := t.finalizeDirtyMember(ctx, key, pairTrafficDirtySetKey, member, version); err != nil && firstErr == nil {
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

	var upload, download int64

	var ut trafficModel.UserTraffic
	err = t.db.WithContext(ctx).Where("user_id = ?", userID).First(&ut).Error
	if err == nil {
		upload = ut.RealUpload
		download = ut.RealDownload
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	} else {
		var user userModel.User
		uErr := t.db.WithContext(ctx).
			Select("uploaded", "downloaded").
			Where("id = ?", userID).
			First(&user).Error
		if uErr != nil && !errors.Is(uErr, gorm.ErrRecordNotFound) {
			return uErr
		}
		if uErr == nil {
			upload = user.Uploaded
			download = user.Downloaded
		}
	}

	pipe := t.rds.TxPipeline()
	pipe.HSetNX(ctx, key, "real_upload", upload)
	pipe.HSetNX(ctx, key, "real_download", download)
	pipe.HSetNX(ctx, key, trafficVersionField, 0)
	_, err = pipe.Exec(ctx)
	return err
}

func (t *traffic) ensureTorrentTotalsInitialized(ctx context.Context, torrentID string) error {
	key := torrentTrafficKey(torrentID)
	exists, err := t.rds.Exists(ctx, key).Result()
	if err != nil || exists > 0 {
		return err
	}

	var totalUpload, totalDownload int64
	var ts trafficModel.TorrentStats
	if err := t.db.WithContext(ctx).Where("torrent_id = ?", torrentID).First(&ts).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	} else {
		totalUpload = ts.TotalUpload
		totalDownload = ts.TotalDownload
	}

	pipe := t.rds.TxPipeline()
	pipe.HSetNX(ctx, key, "total_upload", totalUpload)
	pipe.HSetNX(ctx, key, "total_download", totalDownload)
	pipe.HSetNX(ctx, key, trafficVersionField, 0)
	_, err = pipe.Exec(ctx)
	return err
}

func (t *traffic) ensurePairTotalsInitialized(ctx context.Context, userID, torrentID string) error {
	key := pairTrafficKey(userID, torrentID)
	exists, err := t.rds.Exists(ctx, key).Result()
	if err != nil || exists > 0 {
		return err
	}

	var uploaded, downloaded int64
	var isActive, isFinished bool
	lastAction := time.Now()

	var th trafficModel.TransferHistory
	err = t.db.WithContext(ctx).
		Where("user_id = ? AND torrent_id = ?", userID, torrentID).
		Order("created_at ASC").
		First(&th).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err == nil {
		uploaded = th.Uploaded
		downloaded = th.Downloaded
		isActive = th.IsActive
		isFinished = th.IsFinished
		if !th.LastAction.IsZero() {
			lastAction = th.LastAction
		}
	}

	pipe := t.rds.TxPipeline()
	pipe.HSetNX(ctx, key, "uploaded", uploaded)
	pipe.HSetNX(ctx, key, "downloaded", downloaded)
	pipe.HSetNX(ctx, key, "is_active", boolToInt(isActive))
	pipe.HSetNX(ctx, key, "is_finished", boolToInt(isFinished))
	pipe.HSetNX(ctx, key, "last_action", lastAction.Unix())
	pipe.HSetNX(ctx, key, trafficVersionField, 0)
	_, err = pipe.Exec(ctx)
	return err
}

func userTrafficKey(userID string) string {
	return "tracker:traffic:user:" + userID
}

func userRealtimeUploadKey(userID string) string {
	return "tracker:traffic:realtime:user:" + userID + ":upload"
}

func userRealtimeDownloadKey(userID string) string {
	return "tracker:traffic:realtime:user:" + userID + ":download"
}

func torrentTrafficKey(torrentID string) string {
	return "tracker:traffic:torrent:" + torrentID
}

func pairTrafficKey(userID, torrentID string) string {
	return "tracker:traffic:pair:" + userID + ":" + torrentID
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
