package dao

import (
	"context"
	"errors"
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

var (
	ErrBindInfoAlreadyUsed = errors.New("bind info already used")
)

type user struct {
	stdao.Std[*model.User]
}

func (u *user) Init(db *gorm.DB) error {
	return u.Std.Init(db)
}

func (u *user) NewUserWithBind(ctx context.Context, user *model.User, bind *model.Bind) (*model.User, error) {
	db := u.DB().WithContext(ctx)
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Where("open_id = ?", bind.OpenID).First(&model.Bind{}).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, err
	} else if err == nil {
		tx.Rollback()
		return nil, ErrBindInfoAlreadyUsed
	}

	bind.UID = user.ID
	if err := tx.Create(bind).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Here regard new user as a normal
	var normalGroup model.Group
	if err := tx.Where("name = ?", "normalUser").First(&normalGroup).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	memberShip := model.GroupMembership{
		UID: user.ID,
		GID: normalGroup.ID,
	}

	if err := tx.Create(&memberShip).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var adminGroup model.Group
	if err := tx.Where("name = ?", "admin").First(&adminGroup).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// // set admin has to permission to write
	// err := rbac.CasbinManager.SetUserUnderGroup(user.ID, adminGroup.ID, rbacValues.Write)
	// if err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u *user) GetInfo(ctx context.Context, uid string) (*model.User, error) {
	var user model.User
	if err := u.GetTxFromCtx(ctx).Where("id = ?", uid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *user) IsAdmin(ctx context.Context, uid string) (bool, error) {
	var user model.User
	if err := u.DB().WithContext(ctx).Where("id = ?", uid).First(&user).Error; err != nil {
		return false, err
	}
	return user.IsAdmin, nil
}

func (u *user) UpdateInfo(ctx context.Context, _user model.User) error {
	return u.DB().Model(&model.User{}).Where("id = ?", _user.ID).Updates(&model.User{
		Name:    _user.Name,
		IsAdmin: _user.IsAdmin,
		Limit:   _user.Limit,
	}).Error
}

func (u *user) GetList(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	err := u.GetTxFromCtx(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
