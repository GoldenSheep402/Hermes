package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoldenSheep402/Hermes/mod/category/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type category struct {
	stdao.Std[*model.Category]
}

func (c *category) Init(db *gorm.DB) error {
	return c.Std.Init(db)
}

func (c *category) Create(ctx context.Context, category *model.Category, categoryMetas []model.Metadata) error {
	_ctx := c.SetTxToCtx(ctx, c.DB())
	tx := c.GetTxFromCtx(_ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(category).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range categoryMetas {
		categoryMetas[i].CategoryID = category.ID
	}

	if err := tx.Create(categoryMetas).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (c *category) Get(ctx context.Context, categoryId string) (*model.Category, []model.Metadata, error) {
	_ctx := c.SetTxToCtx(ctx, c.DB())
	db := c.GetTxFromCtx(_ctx)
	category := &model.Category{}
	err := db.Model(&model.Category{}).Where("id = ?", categoryId).First(&category).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, nil, status.Error(codes.NotFound, "Category not found")
		default:
			return nil, nil, status.Error(codes.Internal, "Internal error")
		}
	}

	var categoryMetas []model.Metadata
	err = db.Model(&model.Metadata{}).Where("category_id = ?", categoryId).Find(&categoryMetas).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, nil, status.Error(codes.NotFound, "Category not found")
		default:
			return nil, nil, status.Error(codes.Internal, "Internal error")
		}
	}

	return category, categoryMetas, nil
}

// Update updates a category and its associated metadata in the database within a transaction.
// It ensures data consistency by rolling back the transaction if any operation fails.
// Returns an error if the category is not found, or if there's an issue with the database operation.
func (c *category) Update(ctx context.Context, category *model.Category, categoryMetas []model.Metadata) error {
	_ctx := c.SetTxToCtx(ctx, c.DB())
	tx := c.GetTxFromCtx(_ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingCategory model.Category
	if err := tx.Model(&model.Category{}).Where("id = ?", category.ID).First(&existingCategory).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return status.Error(codes.NotFound, "Category not found")
		}
		return status.Error(codes.Internal, "Internal error")
	}

	if err := tx.Model(&existingCategory).Updates(category).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, meta := range categoryMetas {
		meta.CategoryID = category.ID
		if meta.ID == "" {
			if err := tx.Create(&meta).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := tx.Model(&model.Metadata{}).
				Where("id = ? AND category_id = ?", meta.ID, existingCategory.ID).
				Select("default_value", "order", "description").
				Updates(meta).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	orderMap := make(map[int]bool)
	var metasToCheck []model.Metadata
	if err := tx.Model(&model.Metadata{}).Where("category_id = ?", existingCategory.ID).Find(&metasToCheck).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, meta := range metasToCheck {
		if orderMap[meta.Order] {
			tx.Rollback()
			return status.Error(codes.InvalidArgument, fmt.Sprintf("Duplicate order value: %d", meta.Order))
		}

		orderMap[meta.Order] = true
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (c *category) Delete(ctx context.Context, categoryId string) error {
	_ctx := c.SetTxToCtx(ctx, c.DB())
	tx := c.GetTxFromCtx(_ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var category model.Category
	if err := tx.Model(&model.Category{}).Where("id = ?", categoryId).First(&category).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return status.Error(codes.NotFound, "Category not found")
		}
		return status.Error(codes.Internal, "Internal error")
	}

	if err := tx.Model(&model.Category{}).Where("id = ?", categoryId).Delete(&model.Category{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.Metadata{}).Where("category_id = ?", categoryId).Delete(&model.Metadata{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
