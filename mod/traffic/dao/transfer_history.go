package dao

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/traffic/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type transferHistory struct {
	stdao.Std[*model.TransferHistory]
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
