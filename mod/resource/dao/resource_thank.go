package dao

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/resource/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type resourceThank struct {
	stdao.Std[*model.ResourceThank]
}

func (d *resourceThank) Exists(ctx context.Context, resourceID, userID string) (bool, error) {
	var count int64
	err := d.GetTxFromCtx(ctx).WithContext(ctx).Model(&model.ResourceThank{}).
		Where("resource_id = ? AND user_id = ?", resourceID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (d *resourceThank) AddThank(ctx context.Context, rt *model.ResourceThank) error {
	exists, err := d.Exists(ctx, rt.ResourceID, rt.UserID)
	if err != nil {
		return err
	}
	if exists {
		return nil // Alternatively, return an error, but silent ignore is fine for multiple clicks
	}

	return d.GetTxFromCtx(ctx).WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(rt).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Resource{}).Where("id = ?", rt.ResourceID).UpdateColumn("thank_count", gorm.Expr("thank_count + ?", 1)).Error; err != nil {
			return err
		}
		return nil
	})
}
