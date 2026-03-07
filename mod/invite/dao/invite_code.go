package dao

import (
	"context"
	"time"

	"github.com/GoldenSheep402/Hermes/mod/invite/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type inviteCode struct {
	stdao.Std[*model.InviteCode]
}

func (d *inviteCode) GetByCode(ctx context.Context, code string) (*model.InviteCode, error) {
	var val model.InviteCode
	if err := d.GetTxFromCtx(ctx).WithContext(ctx).Where("code = ?", code).First(&val).Error; err != nil {
		return nil, err
	}
	return &val, nil
}

func (d *inviteCode) ListBySender(ctx context.Context, senderID string, page, pageSize int32) ([]model.InviteCode, int64, error) {
	db := d.GetTxFromCtx(ctx).WithContext(ctx).Where("sender_id = ?", senderID)

	var count int64
	if err := db.Model(&model.InviteCode{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var list []model.InviteCode
	if err := db.Order("created_at DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

// Consume marks an invite code as used within a transaction.
func (d *inviteCode) Consume(ctx context.Context, code string, receiverID string) error {
	now := time.Now()
	return d.GetTxFromCtx(ctx).WithContext(ctx).Model(&model.InviteCode{}).
		Where("code = ? AND is_used = ?", code, false).
		Updates(map[string]interface{}{
			"is_used":     true,
			"receiver_id": receiverID,
			"used_at":     now,
		}).Error
}
