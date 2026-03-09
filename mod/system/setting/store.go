package setting

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/system/dao"
	"github.com/GoldenSheep402/Hermes/mod/system/model"
)

func List(ctx context.Context) ([]*model.Setting, error) {
	list, err := dao.Setting.List(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]*model.Setting, 0, len(list))
	for _, item := range list {
		if item == nil {
			continue
		}

		row := *item
		if value, ok := ManagedValueString(ctx, row.Key); ok {
			row.Value = value
		}
		if itemType, desc, ok := ManagedTypeAndDesc(row.Key); ok {
			row.Type = itemType
			if row.Desc == "" {
				row.Desc = desc
			}
		}
		out = append(out, &row)
	}

	return out, nil
}

func Get(ctx context.Context, key string) (*model.Setting, error) {
	row, err := dao.Setting.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	out := *row
	if value, ok := ManagedValueString(ctx, key); ok {
		out.Value = value
	}
	if itemType, desc, ok := ManagedTypeAndDesc(key); ok {
		out.Type = itemType
		if out.Desc == "" {
			out.Desc = desc
		}
	}
	return &out, nil
}

func Upsert(ctx context.Context, key, value, itemType, desc string) error {
	if managed, err := UpdateManagedValueFromString(ctx, key, value); managed {
		return err
	}

	if itemType == "" {
		itemType = "string"
	}

	err := dao.Setting.UpdateOrCreate(ctx, &model.Setting{
		Key:   key,
		Value: value,
		Type:  itemType,
		Desc:  desc,
	})
	if err == nil {
		settingCache.Delete(key)
	}
	return err
}

func Delete(ctx context.Context, key string) error {
	err := dao.Setting.DeleteByKey(ctx, key)
	if err == nil {
		settingCache.Delete(key)
	}
	return err
}
