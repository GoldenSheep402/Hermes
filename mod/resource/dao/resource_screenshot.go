package dao

import (
	"context"
	"strings"

	"github.com/GoldenSheep402/Hermes/mod/resource/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/oklog/ulid/v2"
)

type resourceScreenshot struct {
	stdao.Std[*model.ResourceScreenshot]
}

func (d *resourceScreenshot) ListByResourceID(ctx context.Context, resourceID string) ([]model.ResourceScreenshot, error) {
	var items []model.ResourceScreenshot
	err := d.GetTxFromCtx(ctx).WithContext(ctx).
		Where("resource_id = ?", resourceID).
		Order("sort_order ASC, created_at ASC").
		Find(&items).Error
	return items, err
}

func (d *resourceScreenshot) ListByResourceIDs(ctx context.Context, resourceIDs []string) (map[string][]model.ResourceScreenshot, error) {
	result := make(map[string][]model.ResourceScreenshot, len(resourceIDs))
	if len(resourceIDs) == 0 {
		return result, nil
	}

	var items []model.ResourceScreenshot
	if err := d.GetTxFromCtx(ctx).WithContext(ctx).
		Where("resource_id IN ?", resourceIDs).
		Order("sort_order ASC, created_at ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}

	for _, item := range items {
		result[item.ResourceID] = append(result[item.ResourceID], item)
	}
	return result, nil
}

func (d *resourceScreenshot) ReplaceByResourceID(ctx context.Context, resourceID string, urls []string) error {
	resourceID = strings.TrimSpace(resourceID)
	if resourceID == "" {
		return nil
	}

	db := d.GetTxFromCtx(ctx).WithContext(ctx)
	if err := db.Where("resource_id = ?", resourceID).Delete(&model.ResourceScreenshot{}).Error; err != nil {
		return err
	}

	normalizedURLs := normalizeNonEmptyUnique(urls)
	if len(normalizedURLs) == 0 {
		return nil
	}

	items := make([]model.ResourceScreenshot, 0, len(normalizedURLs))
	for idx, url := range normalizedURLs {
		items = append(items, model.ResourceScreenshot{
			Model:      stdao.Model{ID: ulid.Make().String()},
			ResourceID: resourceID,
			URL:        url,
			SortOrder:  idx,
		})
	}

	return db.Create(&items).Error
}
