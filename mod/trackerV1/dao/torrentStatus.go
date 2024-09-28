package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/model"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/model/trackerV1Values"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type torrentStatus struct {
	stdao.Std[model.TorrentStatus]
	rds *redis.Client
}

func (tsr *torrentStatus) Init(db *gorm.DB, rds *redis.Client) error {
	err := tsr.Std.Init(db)
	if err != nil {
		return err
	}

	tsr.rds = rds
	return nil
}

func (tsr *torrentStatus) IncrementUploadAndDownload(ctx context.Context, torrentID, userID string, status int, uploadInc, downloadInc int64) error {
	key := "TorrentStatus:" + torrentID

	exists, err := tsr.rds.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to check if key exists: %v", err)
	}

	if exists == 0 {
		err := tsr.LoadInitialData(ctx, torrentID)
		//if err != nil {
		//	return fmt.Errorf("failed to load initial data: %v", err)
		//}
		if err != nil {
			switch {
			case errors.Is(err, gorm.ErrRecordNotFound):

			}
		}

	}

	err = tsr.rds.HIncrBy(ctx, key, "upload_sum", uploadInc).Err()
	if err != nil {
		return fmt.Errorf("failed to increment upload_sum: %v", err)
	}

	err = tsr.rds.HIncrBy(ctx, key, "download_sum", downloadInc).Err()
	if err != nil {
		return fmt.Errorf("failed to increment download_sum: %v", err)
	}

	err = tsr.rds.HSet(ctx, key, "last_active", strconv.FormatInt(time.Now().Unix(), 10)).Err()
	if err != nil {
		return fmt.Errorf("failed to update last_active: %v", err)
	}

	err = tsr.rds.Expire(ctx, key, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration for torrent status: %v", err)
	}

	seedingKey := "TorrentSeeding:" + torrentID
	// Handle seeding count
	switch status {
	case trackerV1Values.Downloading:
		tsr.rds.HSet(ctx, seedingKey, userID, trackerV1Values.Downloading)
	case trackerV1Values.Seeding:
		tsr.rds.HSet(ctx, seedingKey, userID, trackerV1Values.Seeding)
	case trackerV1Values.Finished:
		tsr.rds.HSet(ctx, seedingKey, userID, trackerV1Values.Finished)
	}

	return nil
}

// LoadInitialData loads initial data from db and set into redis for a torrent
func (tsr *torrentStatus) LoadInitialData(ctx context.Context, torrentID string) error {
	key := "TorrentStatus:" + torrentID

	ts, err := tsr.GetFromDbByTID(ctx, torrentID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			ts.UploadSum = 0
			ts.DownloadSum = 0
		default:
			return fmt.Errorf("unknow error: %v", err)
		}
	}

	initialData := map[string]interface{}{
		"upload_sum":     ts.UploadSum,
		"upload_count":   0,
		"download_sum":   ts.DownloadSum,
		"download_count": 0,
		"seeding_count":  0,
		"last_active":    strconv.FormatInt(time.Now().Unix(), 10),
	}

	err = tsr.rds.HSet(ctx, key, initialData).Err()
	if err != nil {
		return fmt.Errorf("failed to load initial data into redis: %v", err)
	}

	err = tsr.rds.Expire(ctx, key, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration for initial data: %v", err)
	}

	return nil
}

func (tsr *torrentStatus) GetFromDbByTID(ctx context.Context, torrentID string) (*model.TorrentStatus, error) {
	var ts model.TorrentStatus
	err := tsr.DB().WithContext(ctx).Where("torrent_id = ?", torrentID).First(&ts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get torrent status from db: %v", err)
	}

	return &ts, nil
}
