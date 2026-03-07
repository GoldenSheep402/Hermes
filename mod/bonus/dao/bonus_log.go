package dao

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/bonus/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type bonusLog struct {
	stdao.Std[*model.BonusLog]
}

func (d *bonusLog) ListByUserID(ctx context.Context, userID string, page, pageSize int32) ([]model.BonusLog, int64, error) {
	db := d.GetTxFromCtx(ctx).WithContext(ctx).Where("user_id = ?", userID)

	var count int64
	if err := db.Model(&model.BonusLog{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var list []model.BonusLog
	if err := db.Order("created_at DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
