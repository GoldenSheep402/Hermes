package settle

import (
	"context"
	"fmt"
	"math"
	"time"

	bonusDao "github.com/GoldenSheep402/Hermes/mod/bonus/dao"
	systemSetting "github.com/GoldenSheep402/Hermes/mod/system/setting"
	torrentDao "github.com/GoldenSheep402/Hermes/mod/torrent/dao"
	trackerDao "github.com/GoldenSheep402/Hermes/mod/tracker/dao"
	"go.uber.org/zap"
)

const (
	t0Weeks = 1.0
	n0      = 7.0
)

// RunOnce settles seeding bonus for all active finished pairs.
func RunOnce(ctx context.Context, log *zap.SugaredLogger) (settled int, err error) {
	if !systemSetting.BonusEnabledValue(ctx) {
		return 0, nil
	}
	multiplier := systemSetting.BonusMultiplierValue(ctx)
	now := time.Now().Unix()

	var cursor uint64
	for {
		members, next, scanErr := trackerDao.Traffic.ScanSeedingPairs(ctx, cursor, 200)
		if scanErr != nil {
			return settled, scanErr
		}
		for _, member := range members {
			userID, torrentID, ok := parsePairMember(member)
			if !ok {
				continue
			}
			if creditErr := settlePair(ctx, log, userID, torrentID, now, multiplier); creditErr != nil {
				if log != nil {
					log.Warnw("bonus settle pair failed", "user_id", userID, "torrent_id", torrentID, "err", creditErr)
				}
				continue
			}
			settled++
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return settled, nil
}

func settlePair(ctx context.Context, log *zap.SugaredLogger, userID, torrentID string, nowUnix int64, multiplier float64) error {
	seedTime, settledAt, ok, err := trackerDao.Traffic.GetPairSeedState(ctx, userID, torrentID)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	if settledAt <= 0 {
		return trackerDao.Traffic.MarkPairBonusSettled(ctx, userID, torrentID, nowUnix)
	}
	elapsed := nowUnix - settledAt
	if elapsed < 60 {
		return nil
	}
	elapsedHours := float64(elapsed) / 3600.0

	torrent, err := torrentDao.Torrent.GetBase(ctx, torrentID)
	if err != nil || torrent == nil {
		return err
	}
	sizeGiB := float64(torrent.Size) / float64(1024*1024*1024)
	if sizeGiB <= 0 {
		return trackerDao.Traffic.MarkPairBonusSettled(ctx, userID, torrentID, nowUnix)
	}

	seeders, err := trackerDao.Traffic.CountTorrentSeeders(ctx, torrentID)
	if err != nil {
		return err
	}
	if seeders < 1 {
		return trackerDao.Traffic.MarkPairBonusSettled(ctx, userID, torrentID, nowUnix)
	}

	weeks := float64(seedTime) / (3600.0 * 24.0 * 7.0)
	points := seedingBonusPoints(weeks, sizeGiB, float64(seeders), multiplier) * elapsedHours
	milli := int64(math.Floor(points * 1000))
	if milli > 0 {
		desc := fmt.Sprintf("seeding %.2f GiB for %.2fh (N=%d)", sizeGiB, elapsedHours, seeders)
		if err := bonusDao.BonusLog.CreditSeedingBonus(ctx, userID, torrentID, milli, desc); err != nil {
			return err
		}
	}
	return trackerDao.Traffic.MarkPairBonusSettled(ctx, userID, torrentID, nowUnix)
}

// seedingBonusPoints is the NexusPHP/M-Team hourly rate:
// (1 - 10^(-T/T0)) * S * (1 + sqrt(2) * 10^(-(N-1)/(N0-1))) * multiplier
func seedingBonusPoints(weeks, sizeGiB, seeders, multiplier float64) float64 {
	if sizeGiB <= 0 || seeders < 1 || multiplier <= 0 {
		return 0
	}
	timeFactor := 1 - math.Pow(10, -weeks/t0Weeks)
	seedFactor := 1 + math.Sqrt2*math.Pow(10, -(seeders-1)/(n0-1))
	return timeFactor * sizeGiB * seedFactor * multiplier
}

func parsePairMember(member string) (string, string, bool) {
	for i := 0; i < len(member); i++ {
		if member[i] == ':' {
			if i == 0 || i == len(member)-1 {
				return "", "", false
			}
			return member[:i], member[i+1:], true
		}
	}
	return "", "", false
}
