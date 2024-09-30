package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type singleSum struct {
	stdao.Std[model.SingleSum]
	rds *redis.Client
}

func (ss *singleSum) Init(db *gorm.DB, rds *redis.Client) error {
	err := ss.Std.Init(db)
	if err != nil {
		return nil
	}

	ss.rds = rds

	return nil
}

func (ss *singleSum) UpdateSingleSum(ctx context.Context, torrentID, uid string, Upload int64, Download int64) error {
	db := ss.DB().WithContext(ctx)
	var singleSum model.SingleSum

	err := db.Where("uid = ? AND torrent_id = ?", uid, torrentID).First(&singleSum).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newSingleSum := model.SingleSum{
				UID:       uid,
				TorrentID: torrentID,
				Upload:    Upload,
				Download:  Download,
			}
			if err := db.Model(&model.SingleSum{}).Create(&newSingleSum).Error; err != nil {
				return fmt.Errorf("failed to create new single sum record: %v", err)
			}
		} else {
			return fmt.Errorf("failed to query single sum record: %v", err)
		}
	} else {
		updated := false
		if Upload > singleSum.Upload {
			singleSum.Upload = Upload
			updated = true
		}
		if Download > singleSum.Download {
			singleSum.Download = Download
			updated = true
		}

		if updated {
			if err := db.Model(&model.SingleSum{}).Where("uid = ? AND torrent_id = ?", uid, torrentID).Updates(&singleSum).Error; err != nil {
				return fmt.Errorf("failed to update single sum record: %v", err)
			}
		}
	}

	return nil
}

func (ss *singleSum) IncreaseSingleSum(ctx context.Context, torrentID, uid string, upload int64, download int64) error {
	db := ss.DB().WithContext(ctx)
	var singleSum model.SingleSum

	err := db.Where("uid = ? AND torrent_id = ?", uid, torrentID).First(&singleSum).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newSingleSum := model.SingleSum{
				UID:       uid,
				TorrentID: torrentID,
				Upload:    upload,
				Download:  download,
			}
			if err := db.Model(&model.SingleSum{}).Create(&newSingleSum).Error; err != nil {
				return fmt.Errorf("failed to create new single sum record: %v", err)
			}
		} else {
			return fmt.Errorf("failed to query single sum record: %v", err)
		}
	} else {
		if err := db.Model(&model.SingleSum{}).
			Where("uid = ? AND torrent_id = ?", uid, torrentID).
			Updates(map[string]interface{}{
				"upload":   gorm.Expr("upload + ?", upload),
				"download": gorm.Expr("download + ?", download),
			}).Error; err != nil {
			return fmt.Errorf("failed to update single sum record: %v", err)
		}
	}

	return nil
}

func (ss *singleSum) GetByTorrentIDandUID(ctx context.Context, torrentID string, uid string) (*model.SingleSum, error) {
	db := ss.DB().WithContext(ctx)
	var singleSums *model.SingleSum

	err := db.Where("uid = ? AND torrent_id = ?", uid, torrentID).First(&singleSums).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query single sum record: %v", err)
	}
	return singleSums, nil
}

func (ss *singleSum) GetFinishedCountByID(ctx context.Context, torrentID string) (int64, error) {
	db := ss.DB().WithContext(ctx)
	var count int64

	err := db.Model(&model.SingleSum{}).Where("torrent_id = ? AND is_finish = ?", torrentID, true).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get finished count by torrent id: %v", err)
	}

	return count, nil
}

// MarkFinished mark the single sum record as finished
// If the record is not exist in redis, update the record in db and set the record in redis and save in db
func (ss *singleSum) MarkFinished(ctx context.Context, torrentID string, uid string) error {
	ok, err := ss.CheckFinishedInRds(ctx, torrentID, uid)
	if err != nil {
		return fmt.Errorf("failed to check finished in redis: %v", err)
	}

	if !ok {
		tx := ss.DB().WithContext(ctx)
		var singleSum model.SingleSum

		err = tx.Where("uid = ? AND torrent_id = ?", uid, torrentID).First(&singleSum).Error
		if err != nil {
			return fmt.Errorf("failed to query single sum record: %v", err)
		}
		singleSum.IsFinish = true
		err = tx.Model(&model.SingleSum{}).Where("uid = ? AND torrent_id = ?", uid, torrentID).Updates(&singleSum).Error
		if err != nil {
			return fmt.Errorf("failed to update single sum record: %v", err)
		}

		err = ss.MarkFinishedInRds(ctx, torrentID, uid)
		if err != nil {
			return fmt.Errorf("failed to mark finished in redis: %v", err)
		}
	}

	return nil
}

func (ss *singleSum) CheckFinishedInRds(ctx context.Context, torrentID string, uid string) (bool, error) {
	rdsKey := fmt.Sprintf("TorrentFinished:%s:%s", torrentID, uid)
	isExist, err := ss.rds.Exists(ctx, rdsKey).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check finished in redis: %v", err)
	}

	return isExist == 1, nil
}

func (ss *singleSum) MarkFinishedInRds(ctx context.Context, torrentID string, uid string) error {
	rdsKey := fmt.Sprintf("TorrentFinished:%s:%s", torrentID, uid)
	if _, err := ss.rds.Set(ctx, rdsKey, 1, 1*time.Hour).Result(); err != nil {
		return fmt.Errorf("failed to set finished in redis: %v", err)
	}

	return nil
}
