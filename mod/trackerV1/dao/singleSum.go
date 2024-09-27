package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type singleSum struct {
	stdao.Std[model.SingleSum]
}

func (ss *singleSum) Init(db *gorm.DB, rds *redis.Client) error {
	err := ss.Std.Init(db)
	if err != nil {
		return nil
	}

	return nil
}

func (ss *singleSum) UpdateSingleSum(ctx context.Context, uid, torrentID string, Upload int64, Download int64) error {
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
