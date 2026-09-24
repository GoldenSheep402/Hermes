package dao

import (
	"context"
	"errors"
	"time"

	"github.com/GoldenSheep402/Hermes/mod/tracker/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type snatch struct {
	stdao.Std[*model.Snatch]
}

func (s *snatch) Init(db *gorm.DB) error {
	if err := compactDuplicateSnatches(db); err != nil {
		return err
	}
	return s.Std.Init(db)
}

func (s *snatch) UpdateOrCreate(ctx context.Context, snatchData *model.Snatch) error {
	db := s.GetTxFromCtx(ctx).WithContext(ctx)
	if snatchData.ID == "" {
		snatchData.ID = ulid.Make().String()
	}
	if snatchData.LastAction.IsZero() {
		snatchData.LastAction = time.Now()
	}

	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "torrent_id"}, {Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"uploaded":    snatchData.Uploaded,
			"downloaded":  snatchData.Downloaded,
			"seed_time":   snatchData.SeedTime,
			"is_active":   snatchData.IsActive,
			"finished_at": snatchData.FinishedAt,
			"last_action": snatchData.LastAction,
			"updated_at":  time.Now(),
		}),
	}).Create(snatchData).Error
}

func (s *snatch) ListByTorrent(ctx context.Context, torrentID string, page, pageSize int32) ([]model.Snatch, int64, error) {
	db := s.GetTxFromCtx(ctx).WithContext(ctx).Where("torrent_id = ?", torrentID)

	var count int64
	if err := db.Model(&model.Snatch{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var list []model.Snatch
	if err := db.Order("last_action DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (s *snatch) ListByUser(ctx context.Context, userID string, page, pageSize int32) ([]model.Snatch, int64, error) {
	db := s.GetTxFromCtx(ctx).WithContext(ctx).Where("user_id = ?", userID)

	var count int64
	if err := db.Model(&model.Snatch{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var list []model.Snatch
	if err := db.Order("last_action DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

type duplicateSnatchGroup struct {
	TorrentID string
	UserID    string
	Count     int64
}

func compactDuplicateSnatches(db *gorm.DB) error {
	ctx := context.Background()
	if !db.Migrator().HasTable(&model.Snatch{}) {
		return nil
	}

	var groups []duplicateSnatchGroup
	if err := db.WithContext(ctx).
		Model(&model.Snatch{}).
		Select("torrent_id, user_id, COUNT(*) AS count").
		Group("torrent_id, user_id").
		Having("COUNT(*) > 1").
		Scan(&groups).Error; err != nil {
		return err
	}

	for _, group := range groups {
		if group.TorrentID == "" || group.UserID == "" {
			continue
		}

		if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			var rows []model.Snatch
			if err := tx.
				Where("torrent_id = ? AND user_id = ?", group.TorrentID, group.UserID).
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
				if row.FinishedAt != nil && (keeper.FinishedAt == nil || row.FinishedAt.Before(*keeper.FinishedAt)) {
					finishedAt := *row.FinishedAt
					keeper.FinishedAt = &finishedAt
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

			if err := tx.Model(&model.Snatch{}).
				Where("id = ?", keeper.ID).
				Updates(map[string]any{
					"uploaded":    keeper.Uploaded,
					"downloaded":  keeper.Downloaded,
					"seed_time":   keeper.SeedTime,
					"is_active":   keeper.IsActive,
					"finished_at": keeper.FinishedAt,
					"last_action": keeper.LastAction,
					"updated_at":  time.Now(),
				}).Error; err != nil {
				return err
			}

			if len(deleteIDs) == 0 {
				return nil
			}
			if err := tx.Unscoped().Where("id IN ?", deleteIDs).Delete(&model.Snatch{}).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}
