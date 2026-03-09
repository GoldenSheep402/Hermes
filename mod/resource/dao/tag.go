package dao

import (
	"context"
	"strings"

	"github.com/GoldenSheep402/Hermes/mod/resource/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/oklog/ulid/v2"
)

type tag struct {
	stdao.Std[*model.Tag]
}

type resourceTag struct {
	stdao.Std[*model.ResourceTag]
}

func (d *tag) EnsureByNames(ctx context.Context, names []string) ([]model.Tag, error) {
	normalizedNames := normalizeNonEmptyUnique(names)
	if len(normalizedNames) == 0 {
		return []model.Tag{}, nil
	}

	db := d.GetTxFromCtx(ctx).WithContext(ctx)
	var existing []model.Tag
	if err := db.Where("name IN ?", normalizedNames).Find(&existing).Error; err != nil {
		return nil, err
	}

	existingMap := make(map[string]model.Tag, len(existing))
	for _, item := range existing {
		existingMap[item.Name] = item
	}

	toCreate := make([]model.Tag, 0)
	for _, name := range normalizedNames {
		if _, ok := existingMap[name]; ok {
			continue
		}
		toCreate = append(toCreate, model.Tag{
			Model: stdao.Model{ID: ulid.Make().String()},
			Name:  name,
		})
	}
	if len(toCreate) > 0 {
		if err := db.Create(&toCreate).Error; err != nil {
			return nil, err
		}
		for _, item := range toCreate {
			existingMap[item.Name] = item
		}
	}

	result := make([]model.Tag, 0, len(normalizedNames))
	for _, name := range normalizedNames {
		if item, ok := existingMap[name]; ok {
			result = append(result, item)
		}
	}
	return result, nil
}

func (d *resourceTag) ReplaceByResourceID(ctx context.Context, resourceID string, tagIDs []string) error {
	resourceID = strings.TrimSpace(resourceID)
	if resourceID == "" {
		return nil
	}

	db := d.GetTxFromCtx(ctx).WithContext(ctx)
	if err := db.Where("resource_id = ?", resourceID).Delete(&model.ResourceTag{}).Error; err != nil {
		return err
	}

	normalizedTagIDs := normalizeNonEmptyUnique(tagIDs)
	if len(normalizedTagIDs) == 0 {
		return nil
	}

	items := make([]model.ResourceTag, 0, len(normalizedTagIDs))
	for _, tagID := range normalizedTagIDs {
		items = append(items, model.ResourceTag{
			Model:      stdao.Model{ID: ulid.Make().String()},
			ResourceID: resourceID,
			TagID:      tagID,
		})
	}

	return db.Create(&items).Error
}

func (d *resourceTag) ListTagMapByResourceIDs(ctx context.Context, resourceIDs []string) (map[string][]model.Tag, error) {
	result := make(map[string][]model.Tag, len(resourceIDs))
	if len(resourceIDs) == 0 {
		return result, nil
	}

	db := d.GetTxFromCtx(ctx).WithContext(ctx)

	var links []model.ResourceTag
	if err := db.Where("resource_id IN ?", resourceIDs).Order("created_at ASC").Find(&links).Error; err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return result, nil
	}

	tagIDSet := make(map[string]struct{}, len(links))
	for _, link := range links {
		if link.TagID == "" {
			continue
		}
		tagIDSet[link.TagID] = struct{}{}
	}
	if len(tagIDSet) == 0 {
		return result, nil
	}

	tagIDs := make([]string, 0, len(tagIDSet))
	for tagID := range tagIDSet {
		tagIDs = append(tagIDs, tagID)
	}

	var tags []model.Tag
	if err := db.Model(&model.Tag{}).Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		return nil, err
	}

	tagMap := make(map[string]model.Tag, len(tags))
	for _, item := range tags {
		tagMap[item.ID] = item
	}

	seen := make(map[string]map[string]struct{}, len(resourceIDs))
	for _, link := range links {
		tagEntity, ok := tagMap[link.TagID]
		if !ok {
			continue
		}
		if _, ok := seen[link.ResourceID]; !ok {
			seen[link.ResourceID] = make(map[string]struct{})
		}
		if _, exists := seen[link.ResourceID][link.TagID]; exists {
			continue
		}
		seen[link.ResourceID][link.TagID] = struct{}{}
		result[link.ResourceID] = append(result[link.ResourceID], tagEntity)
	}

	return result, nil
}

func normalizeNonEmptyUnique(items []string) []string {
	result := make([]string, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, raw := range items {
		text := strings.TrimSpace(raw)
		if text == "" {
			continue
		}
		if _, ok := seen[text]; ok {
			continue
		}
		seen[text] = struct{}{}
		result = append(result, text)
	}
	return result
}
