package dao

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/tracker/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type snatch struct {
	stdao.Std[*model.Snatch]
}

func (s *snatch) Init(db *gorm.DB) error {
	return s.Std.Init(db)
}

func (s *snatch) UpdateOrCreate(ctx context.Context, snatchData *model.Snatch) error {
	db := s.GetTxFromCtx(ctx).WithContext(ctx)
	// Unique by TorrentID and UserID
	return db.Where("torrent_id = ? AND user_id = ?", snatchData.TorrentID, snatchData.UserID).
		Assign(model.Snatch{
			Uploaded:   snatchData.Uploaded,
			Downloaded: snatchData.Downloaded,
			SeedTime:   snatchData.SeedTime,
			IsActive:   snatchData.IsActive,
			FinishedAt: snatchData.FinishedAt,
			LastAction: snatchData.LastAction,
		}).FirstOrCreate(snatchData).Error
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
