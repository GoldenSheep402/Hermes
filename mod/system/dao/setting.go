package dao

import (
	"context"
	"github.com/GoldenSheep402/Hermes/mod/system/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type setting struct {
	stdao.Std[*model.Setting]
}

func (s *setting) Init(db *gorm.DB) error {
	return s.Std.Init(db)
}

func (s *setting) GetSettings(ctx context.Context) (*model.Setting, []model.Subnet, error) {
	var setting model.Setting
	err := s.Std.DB().WithContext(ctx).Model(&model.Setting{}).First(&setting).Error
	if err != nil {
		return nil, nil, err
	}

	var subnets []model.Subnet
	err = s.Std.DB().WithContext(ctx).Model(&model.Subnet{}).Find(&subnets).Error
	if err != nil {
		return nil, nil, err
	}

	return &setting, subnets, nil
}

func (s *setting) SetSettings(ctx context.Context, setting *model.Setting, subnets []model.Subnet) error {
	tx := s.DB().WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := tx.Model(&model.Setting{}).Save(setting).Error; err != nil {
		return err
	}

	if err := tx.Delete(&model.Subnet{}).Error; err != nil {
		return err
	}

	for _, subnet := range subnets {
		if err := tx.Model(&model.Subnet{}).Save(&subnet).Error; err != nil {
			return err
		}
	}

	return tx.Commit().Error
}
