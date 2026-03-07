package dao

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/category/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type category struct {
	stdao.Std[*model.Category]
}

// GetByID fetches a specific category by its ID
func (d *category) GetByID(ctx context.Context, id string) (*model.Category, error) {
	var cat model.Category
	if err := d.GetTxFromCtx(ctx).WithContext(ctx).Where("id = ?", id).First(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}
