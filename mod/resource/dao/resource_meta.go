package dao

import (
	"context"
	"strings"

	"github.com/GoldenSheep402/Hermes/mod/resource/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type resourceMeta struct {
	stdao.Std[*model.ResourceMeta]
}

func (d *resourceMeta) ListByResourceID(ctx context.Context, resourceID string) ([]model.ResourceMeta, error) {
	var items []model.ResourceMeta
	err := d.GetTxFromCtx(ctx).WithContext(ctx).
		Where("resource_id = ?", resourceID).
		Order("created_at ASC").
		Find(&items).Error
	return items, err
}

func (d *resourceMeta) ListByResourceIDs(ctx context.Context, resourceIDs []string) (map[string][]model.ResourceMeta, error) {
	result := make(map[string][]model.ResourceMeta, len(resourceIDs))
	if len(resourceIDs) == 0 {
		return result, nil
	}

	var items []model.ResourceMeta
	if err := d.GetTxFromCtx(ctx).WithContext(ctx).
		Where("resource_id IN ?", resourceIDs).
		Order("created_at ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}

	for _, item := range items {
		result[item.ResourceID] = append(result[item.ResourceID], item)
	}
	return result, nil
}

func (d *resourceMeta) ReplaceByResourceID(ctx context.Context, resourceID string, items []model.ResourceMeta) error {
	resourceID = strings.TrimSpace(resourceID)
	if resourceID == "" {
		return nil
	}

	db := d.GetTxFromCtx(ctx).WithContext(ctx)
	if err := db.Where("resource_id = ?", resourceID).Delete(&model.ResourceMeta{}).Error; err != nil {
		return err
	}

	if len(items) == 0 {
		return nil
	}

	normalized := make([]model.ResourceMeta, 0, len(items))
	for _, item := range items {
		key := strings.TrimSpace(item.Key)
		value := strings.TrimSpace(item.Value)
		if key == "" || value == "" {
			continue
		}
		normalized = append(normalized, model.ResourceMeta{
			Model:      item.Model,
			ResourceID: resourceID,
			Key:        key,
			Value:      value,
		})
	}

	if len(normalized) == 0 {
		return nil
	}

	return db.Create(&normalized).Error
}
