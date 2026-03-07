package dao

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/comment/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type comment struct {
	stdao.Std[*model.Comment]
}

func (d *comment) GetByID(ctx context.Context, id string) (*model.Comment, error) {
	var c model.Comment
	if err := d.GetTxFromCtx(ctx).WithContext(ctx).Where("id = ?", id).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (d *comment) ListByResource(ctx context.Context, resourceID string, page, pageSize int32) ([]model.Comment, int64, error) {
	db := d.GetTxFromCtx(ctx).WithContext(ctx).Where("resource_id = ? AND parent_id IS NULL", resourceID)

	var count int64
	if err := db.Model(&model.Comment{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var list []model.Comment
	if err := db.Order("created_at ASC").Offset(int(offset)).Limit(int(pageSize)).Preload("Replies").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}
