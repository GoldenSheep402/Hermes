package setting

import (
	"context"
	"errors"
	"strconv"

	systemDao "github.com/GoldenSheep402/Hermes/mod/system/dao"
	systemModel "github.com/GoldenSheep402/Hermes/mod/system/model"
	"gorm.io/gorm"
)

type IntItem struct {
	Key          string
	Desc         string
	DefaultValue int
	MinValue     int
	MaxValue     int
}

var (
	TrackerFlushInterval = IntItem{
		Key:          systemDao.SettingKeyTrackerFlushInterval,
		Desc:         "Tracker traffic flush interval (seconds)",
		DefaultValue: 15,
		MinValue:     1,
		MaxValue:     600,
	}
	TrackerFlushBatchSize = IntItem{
		Key:          systemDao.SettingKeyTrackerFlushBatchSize,
		Desc:         "Tracker traffic flush batch size",
		DefaultValue: 500,
		MinValue:     50,
		MaxValue:     5000,
	}
)

var intItems = []IntItem{
	TrackerFlushInterval,
	TrackerFlushBatchSize,
}

// Init migrates default setting items into DB when they do not exist.
func Init(ctx context.Context) error {
	if systemDao.Setting.DB() == nil {
		return nil
	}

	for _, item := range intItems {
		if err := ensureIntItem(ctx, item); err != nil {
			return err
		}
	}
	return nil
}

func TrackerFlushIntervalValue(ctx context.Context) int {
	return GetIntValue(ctx, TrackerFlushInterval)
}

func TrackerFlushBatchSizeValue(ctx context.Context) int {
	return GetIntValue(ctx, TrackerFlushBatchSize)
}

func GetIntValue(ctx context.Context, item IntItem) int {
	if systemDao.Setting.DB() == nil {
		return item.DefaultValue
	}
	val, err := systemDao.Setting.GetIntByKey(ctx, item.Key)
	if err != nil {
		return item.DefaultValue
	}
	return clamp(val, item.MinValue, item.MaxValue)
}

func ensureIntItem(ctx context.Context, item IntItem) error {
	current, err := systemDao.Setting.GetByKey(ctx, item.Key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return systemDao.Setting.Create(ctx, &systemModel.Setting{
				Key:   item.Key,
				Value: strconv.Itoa(item.DefaultValue),
				Type:  "int",
				Desc:  item.Desc,
			})
		}
		return err
	}

	updates := map[string]interface{}{}
	if current.Type != "int" {
		updates["type"] = "int"
	}
	if current.Desc == "" {
		updates["desc"] = item.Desc
	}
	if current.Value == "" {
		updates["value"] = strconv.Itoa(item.DefaultValue)
	}
	if len(updates) == 0 {
		return nil
	}

	return systemDao.Setting.GetTxFromCtx(ctx).WithContext(ctx).
		Model(&systemModel.Setting{}).
		Where("id = ?", current.ID).
		Updates(updates).Error
}

func clamp(v, minV, maxV int) int {
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}
