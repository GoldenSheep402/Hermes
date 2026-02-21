package dao

import (
	"context"
	"errors"

	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	pageutil "github.com/GoldenSheep402/Hermes/pkg/utils/page"
)

var (
	ErrEmailAlreadyUsed    = errors.New("email already used")
	ErrUsernameAlreadyUsed = errors.New("username already used")
)

type user struct {
	stdao.Std[*model.User]
	rds *redis.Client
}

func (u *user) Init(db *gorm.DB, rds *redis.Client) error {
	u.rds = rds
	return u.Std.Init(db)
}

// CreateUser creates a new user and assigns it to the given group.
func (u *user) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	db := u.DB().WithContext(ctx)
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user.ID = ulid.Make().String()

	// Check email uniqueness
	var count int64
	if err := tx.Model(&model.User{}).Where("email = ?", user.Email).Count(&count).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if count > 0 {
		tx.Rollback()
		return nil, ErrEmailAlreadyUsed
	}

	// Check username uniqueness
	if err := tx.Model(&model.User{}).Where("username = ?", user.Username).Count(&count).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if count > 0 {
		tx.Rollback()
		return nil, ErrUsernameAlreadyUsed
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u *user) GetByID(ctx context.Context, uid string) (*model.User, error) {
	var user model.User
	if err := u.GetTxFromCtx(ctx).Where("id = ?", uid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *user) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := u.Std.GetTxFromCtx(ctx).WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *user) GetByPasskey(ctx context.Context, passkey string) (*model.User, error) {
	var user model.User
	if err := u.Std.GetTxFromCtx(ctx).WithContext(ctx).Where("passkey = ?", passkey).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *user) IsAdmin(ctx context.Context, uid string) (bool, error) {
	var user model.User
	if err := u.Std.GetTxFromCtx(ctx).WithContext(ctx).Where("id = ?", uid).First(&user).Error; err != nil {
		return false, err
	}
	return user.IsAdmin, nil
}

func (u *user) UpdateInfo(ctx context.Context, _user *model.User) error {
	return u.Std.GetTxFromCtx(ctx).WithContext(ctx).Model(&model.User{}).Where("id = ?", _user.ID).Updates(_user).Error
}

func (u *user) GetList(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	err := u.GetTxFromCtx(ctx).Find(&users).Error
	return users, err
}

func (u *user) GetListPaginated(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64
	tx := u.GetTxFromCtx(ctx).Model(&model.User{})
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	p := pageutil.Paginate{Current: page, PageSize: pageSize}
	err := tx.Scopes(p.Paginate()).Find(&users).Error
	return users, total, err
}
