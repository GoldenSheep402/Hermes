package dao

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type userGroup struct {
	stdao.Std[*model.UserGroup]
}

func (g *userGroup) Init(db *gorm.DB) error {
	return g.Std.Init(db)
}

func (g *userGroup) Create(ctx context.Context, group *model.UserGroup) error {
	return g.DB().WithContext(ctx).Create(group).Error
}

func (g *userGroup) Get(ctx context.Context, id string) (*model.UserGroup, error) {
	var group model.UserGroup
	if err := g.GetTxFromCtx(ctx).Where("id = ?", id).First(&group).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (g *userGroup) List(ctx context.Context) ([]*model.UserGroup, error) {
	var groups []*model.UserGroup
	if err := g.GetTxFromCtx(ctx).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (g *userGroup) Update(ctx context.Context, group *model.UserGroup) error {
	return g.DB().WithContext(ctx).Model(&model.UserGroup{}).Where("id = ?", group.ID).Updates(group).Error
}

func (g *userGroup) Delete(ctx context.Context, id string) error {
	return g.DB().WithContext(ctx).Where("id = ?", id).Delete(&model.UserGroup{}).Error
}
