package dao

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/message/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type message struct {
	stdao.Std[*model.Message]
}

func (d *message) GetByID(ctx context.Context, id string) (*model.Message, error) {
	var m model.Message
	if err := d.GetTxFromCtx(ctx).WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *message) ListInbox(ctx context.Context, receiverID string, page, pageSize int32) ([]model.Message, int64, error) {
	db := d.GetTxFromCtx(ctx).WithContext(ctx).Where("receiver_id = ?", receiverID)

	var count int64
	if err := db.Model(&model.Message{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var list []model.Message
	if err := db.Order("created_at DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (d *message) ListOutbox(ctx context.Context, senderID string, page, pageSize int32) ([]model.Message, int64, error) {
	db := d.GetTxFromCtx(ctx).WithContext(ctx).Where("sender_id = ?", senderID)

	var count int64
	if err := db.Model(&model.Message{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var list []model.Message
	if err := db.Order("created_at DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

func (d *message) GetUnreadCount(ctx context.Context, receiverID string) (int64, error) {
	var count int64
	if err := d.GetTxFromCtx(ctx).WithContext(ctx).Model(&model.Message{}).
		Where("receiver_id = ? AND is_read = ?", receiverID, false).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (d *message) MarkAsRead(ctx context.Context, receiverID string, messageIDs []string) error {
	return d.GetTxFromCtx(ctx).WithContext(ctx).Model(&model.Message{}).
		Where("receiver_id = ? AND id IN ?", receiverID, messageIDs).
		Update("is_read", true).Error
}
