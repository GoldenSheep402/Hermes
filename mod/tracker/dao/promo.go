package dao

import (
	"context"
	"time"
)

func torrentPromoKey(torrentID string) string {
	return prefixed("tracker:promo:" + torrentID)
}

func globalFreeKey() string {
	return prefixed("tracker:global_free")
}

func seedingPairsKey() string {
	return prefixed("tracker:seeding_pairs")
}

// PromoFactors holds upload/download accounting multipliers for a torrent.
type PromoFactors struct {
	UploadFactor   int64 // 1 or 2
	DownloadFactor int64 // 0 (free) or 1
	UploadUntil    int64 // unix seconds, 0 = no expiry
	DownloadUntil  int64
}

// SetTorrentPromo writes per-torrent promotion multipliers into Redis.
func (t *traffic) SetTorrentPromo(ctx context.Context, torrentID string, promo PromoFactors) error {
	if t.rds == nil || torrentID == "" {
		return nil
	}
	if promo.UploadFactor <= 0 {
		promo.UploadFactor = 1
	}
	if promo.DownloadFactor < 0 {
		promo.DownloadFactor = 1
	}
	key := torrentPromoKey(torrentID)
	err := t.rds.HSet(ctx, key, map[string]interface{}{
		"up":         promo.UploadFactor,
		"down":       promo.DownloadFactor,
		"up_until":   promo.UploadUntil,
		"down_until": promo.DownloadUntil,
	}).Err()
	if err != nil {
		return err
	}
	_ = t.rds.Persist(ctx, key).Err()
	return nil
}

// ClearTorrentPromo removes per-torrent promotion multipliers.
func (t *traffic) ClearTorrentPromo(ctx context.Context, torrentID string) error {
	if t.rds == nil || torrentID == "" {
		return nil
	}
	return t.rds.Del(ctx, torrentPromoKey(torrentID)).Err()
}

// SetGlobalFreeleech enables or disables site-wide freeleech in Redis.
// durationHours > 0 sets a TTL; 0 means until explicitly cleared.
func (t *traffic) SetGlobalFreeleech(ctx context.Context, enabled bool, durationHours int) error {
	if t.rds == nil {
		return nil
	}
	key := globalFreeKey()
	if !enabled {
		return t.rds.Del(ctx, key).Err()
	}
	if err := t.rds.Set(ctx, key, "1", 0).Err(); err != nil {
		return err
	}
	if durationHours > 0 {
		return t.rds.Expire(ctx, key, time.Duration(durationHours)*time.Hour).Err()
	}
	return t.rds.Persist(ctx, key).Err()
}

// ScanSeedingPairs iterates active finished seeding pairs for bonus settlement.
func (t *traffic) ScanSeedingPairs(ctx context.Context, cursor uint64, count int64) ([]string, uint64, error) {
	if t.rds == nil {
		return nil, 0, nil
	}
	if count <= 0 {
		count = 200
	}
	return t.rds.SScan(ctx, seedingPairsKey(), cursor, "*", count).Result()
}

// GetPairSeedState returns seed_time and bonus_settled_at for a user-torrent pair.
func (t *traffic) GetPairSeedState(ctx context.Context, userID, torrentID string) (seedTime int64, bonusSettledAt int64, ok bool, err error) {
	if t.rds == nil || userID == "" || torrentID == "" {
		return 0, 0, false, nil
	}
	values, err := t.rds.HMGet(ctx, pairTrafficKey(userID, torrentID), "seed_time", "bonus_settled_at", "is_active", "is_finished").Result()
	if err != nil {
		return 0, 0, false, err
	}
	if len(values) < 4 {
		return 0, 0, false, nil
	}
	active := parseAnyInt64(values[2]) == 1
	finished := parseAnyInt64(values[3]) == 1
	if !active || !finished {
		return 0, 0, false, nil
	}
	return parseAnyInt64(values[0]), parseAnyInt64(values[1]), true, nil
}

// MarkPairBonusSettled updates bonus_settled_at for a pair.
func (t *traffic) MarkPairBonusSettled(ctx context.Context, userID, torrentID string, settledAt int64) error {
	if t.rds == nil || userID == "" || torrentID == "" {
		return nil
	}
	return t.rds.HSet(ctx, pairTrafficKey(userID, torrentID), "bonus_settled_at", settledAt).Err()
}

// CountTorrentSeeders returns the current seeder count for a torrent.
func (t *traffic) CountTorrentSeeders(ctx context.Context, torrentID string) (int64, error) {
	if t.rds == nil || torrentID == "" {
		return 0, nil
	}
	if err := cleanupActivityZSet(ctx, t.rds, torrentSeederKey(torrentID)); err != nil {
		return 0, err
	}
	return t.rds.ZCard(ctx, torrentSeederKey(torrentID)).Result()
}

// SyncGlobalFreeleechFromSettings mirrors system settings into Redis.
func (t *traffic) SyncGlobalFreeleechFromSettings(ctx context.Context, enabled bool, countdownHours int) error {
	return t.SetGlobalFreeleech(ctx, enabled, countdownHours)
}

// FormatUnix formats an optional time pointer as unix seconds (0 if nil).
func FormatUnix(t *time.Time) int64 {
	if t == nil || t.IsZero() {
		return 0
	}
	return t.Unix()
}
