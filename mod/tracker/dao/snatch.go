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
