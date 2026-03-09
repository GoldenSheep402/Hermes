package dao

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/resource/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type resource struct {
	stdao.Std[*model.Resource]
}

// GetByID fetches a resource by its string ID.
func (d *resource) GetByID(ctx context.Context, id string) (*model.Resource, error) {
	var res model.Resource
	if err := d.GetTxFromCtx(ctx).WithContext(ctx).Where("id = ?", id).First(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

// IncrementViewCount increments the view stats for a resource.
func (d *resource) IncrementViewCount(ctx context.Context, id string) error {
	return d.GetTxFromCtx(ctx).WithContext(ctx).Model(&model.Resource{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (d *resource) AdjustCommentCount(ctx context.Context, id string, delta int) error {
	if delta == 0 {
		return nil
	}
	return d.GetTxFromCtx(ctx).WithContext(ctx).
		Model(&model.Resource{}).
		Where("id = ?", id).
		UpdateColumn("comment_count", gorm.Expr("CASE WHEN comment_count + ? < 0 THEN 0 ELSE comment_count + ? END", delta, delta)).
		Error
}

// ListByFilters returns a paginated list of resources based on parameters.
func (d *resource) ListByFilters(ctx context.Context, categoryID, keyword string, status int32, page, pageSize int32) ([]model.Resource, int64, error) {
	db := d.GetTxFromCtx(ctx).WithContext(ctx)
	if categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		db = db.Where("title LIKE ?", "%"+keyword+"%")
	}
	if status != 0 {
		db = db.Where("status = ?", status)
	}

	var count int64
	if err := db.Model(&model.Resource{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var list []model.Resource
	if err := db.Order("created_at DESC").Offset(int(offset)).Limit(int(pageSize)).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

// UpdateFields applies specific field updates to a resource.
func (d *resource) UpdateFields(ctx context.Context, id string, updates map[string]interface{}) error {
	return d.GetTxFromCtx(ctx).WithContext(ctx).Model(&model.Resource{}).Where("id = ?", id).Updates(updates).Error
}

func (d *resource) CountPublishedByUser(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := d.GetTxFromCtx(ctx).WithContext(ctx).Model(&model.Resource{}).
		Where("uploader_id = ? AND status = ?", userID, model.ResourceStatusApproved).
		Count(&count).Error
	return count, err
}
