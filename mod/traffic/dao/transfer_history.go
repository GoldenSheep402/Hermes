package dao

import (
	"context"
	"errors"
	"time"

	"github.com/GoldenSheep402/Hermes/mod/traffic/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type transferHistory struct {
	stdao.Std[*model.TransferHistory]
}

func (d *transferHistory) Init(db *gorm.DB) error {
	if err := compactDuplicateTransferHistory(db); err != nil {
		return err
	}
	return d.Std.Init(db)
}

func (d *transferHistory) ListByUserID(ctx context.Context, userID string, page, pageSize int32) ([]model.TransferHistory, int64, error) {
	db := d.GetTxFromCtx(ctx).WithContext(ctx).Where("user_id = ?", userID)

	var count int64
	if err := db.Model(&model.TransferHistory{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var list []model.TransferHistory
	if err := db.Order("last_action DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (d *transferHistory) CountActiveSeeding(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := d.GetTxFromCtx(ctx).WithContext(ctx).Model(&model.TransferHistory{}).
		Where("user_id = ? AND is_active = ? AND is_finished = ?", userID, true, true).
		Count(&count).Error
	return count, err
}

func (d *transferHistory) CountActiveDownloading(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := d.GetTxFromCtx(ctx).WithContext(ctx).Model(&model.TransferHistory{}).
		Where("user_id = ? AND is_active = ? AND is_finished = ?", userID, true, false).
		Count(&count).Error
	return count, err
}

type duplicateTransferHistoryGroup struct {
	UserID    string
	TorrentID string
	Count     int64
}

func compactDuplicateTransferHistory(db *gorm.DB) error {
	ctx := context.Background()
	if !db.Migrator().HasTable(&model.TransferHistory{}) {
		return nil
	}

	var groups []duplicateTransferHistoryGroup
	if err := db.WithContext(ctx).
		Model(&model.TransferHistory{}).
		Select("user_id, torrent_id, COUNT(*) AS count").
		Group("user_id, torrent_id").
		Having("COUNT(*) > 1").
		Scan(&groups).Error; err != nil {
		return err
	}

	for _, group := range groups {
		if group.UserID == "" || group.TorrentID == "" {
			continue
		}

		if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var rows []model.TransferHistory
			if err := tx.
				Where("user_id = ? AND torrent_id = ?", group.UserID, group.TorrentID).
				Order("created_at ASC, id ASC").
				Find(&rows).Error; err != nil {
				return err
			}
			if len(rows) <= 1 {
				return nil
			}

			keeper := rows[0]
			latest := rows[0]
			deleteIDs := make([]string, 0, len(rows)-1)
			for _, row := range rows[1:] {
				deleteIDs = append(deleteIDs, row.ID)
				if row.Uploaded > keeper.Uploaded {
					keeper.Uploaded = row.Uploaded
				}
				if row.Downloaded > keeper.Downloaded {
					keeper.Downloaded = row.Downloaded
				}
				if row.SeedTime > keeper.SeedTime {
					keeper.SeedTime = row.SeedTime
				}
				if row.IsFinished {
					keeper.IsFinished = true
				}
				if row.LastAction.After(latest.LastAction) {
					latest = row
				}
			}

			keeper.IsActive = latest.IsActive
			if latest.LastAction.After(keeper.LastAction) {
				keeper.LastAction = latest.LastAction
			}
			if keeper.LastAction.IsZero() {
				keeper.LastAction = time.Now()
			}

			if err := tx.Model(&model.TransferHistory{}).
				Where("id = ?", keeper.ID).
				Updates(map[string]any{
					"uploaded":    keeper.Uploaded,
					"downloaded":  keeper.Downloaded,
					"seed_time":   keeper.SeedTime,
					"is_finished": keeper.IsFinished,
					"is_active":   keeper.IsActive,
					"last_action": keeper.LastAction,
					"updated_at":  time.Now(),
				}).Error; err != nil {
				return err
			}

			if len(deleteIDs) == 0 {
				return nil
			}
			if err := tx.Unscoped().Where("id IN ?", deleteIDs).Delete(&model.TransferHistory{}).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}
