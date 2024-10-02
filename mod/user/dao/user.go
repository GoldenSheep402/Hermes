package dao

import (
	"context"
	"errors"
	torrentModel "github.com/GoldenSheep402/Hermes/mod/torrent/model"
	trackerV1Model "github.com/GoldenSheep402/Hermes/mod/trackerV1/model"
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

var (
	ErrBindInfoAlreadyUsed = errors.New("bind info already used")
)

type user struct {
	stdao.Std[*model.User]
	rds *redis.Client
}

func (u *user) Init(db *gorm.DB, rds *redis.Client) error {
	u.rds = rds
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

	user.ID = ulid.Make().String()
	user.Key = ulid.Make().String()
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

func (u *user) GetFullInfo(ctx context.Context, uid string) (user *model.User, download, upload, pubished, downloadCount, seedingCount int64, err error) {
	if err := u.DB().WithContext(ctx).Model(&model.User{}).Where("id = ?", uid).First(&user).Error; err != nil {
		return nil, 0, 0, 0, 0, 0, err
	}

	sum := trackerV1Model.Sum{}
	if err := u.DB().WithContext(ctx).Model(&trackerV1Model.Sum{}).Where("uid = ?", uid).First(&sum).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			var singleSums []trackerV1Model.SingleSum
			if err := u.DB().WithContext(ctx).Model(&trackerV1Model.SingleSum{}).Where("uid = ?", uid).Find(&singleSums).Error; err != nil {
				return nil, 0, 0, 0, 0, 0, err
			}

			for _, singleSum := range singleSums {
				download += singleSum.Download
				upload += singleSum.Upload
			}

			// create
			if err := u.DB().WithContext(ctx).Model(&trackerV1Model.Sum{}).Create(&trackerV1Model.Sum{
				UID:          uid,
				RealDownload: download,
				RealUpload:   upload,
			}).Error; err != nil {
				return nil, 0, 0, 0, 0, 0, err
			}
		default:
			return nil, 0, 0, 0, 0, 0, err
		}
	}

	timeNow := time.Now()
	download = 0
	upload = 0

	// TODO: time setting
	if timeNow.Sub(sum.UpdatedAt) > 10*time.Minute {
		var singleSums []trackerV1Model.SingleSum
		if err := u.DB().WithContext(ctx).Model(&trackerV1Model.SingleSum{}).Where("uid = ?", uid).Find(&singleSums).Error; err != nil {
			return nil, 0, 0, 0, 0, 0, err
		}

		for _, singleSum := range singleSums {
			download += singleSum.Download
			upload += singleSum.Upload
		}

		// Update
		if err := u.DB().WithContext(ctx).Model(&trackerV1Model.Sum{}).Where("id = ?", sum.ID).Updates(&trackerV1Model.Sum{
			RealDownload: download,
			RealUpload:   upload,
		}).Error; err != nil {
			return nil, 0, 0, 0, 0, 0, err
		}
	} else {
		download = sum.RealDownload
		upload = sum.RealUpload
	}

	var torrents []torrentModel.Torrent
	if err := u.DB().WithContext(ctx).Model(&torrentModel.Torrent{}).Where("creator_id = ?", uid).Find(&torrents).Error; err != nil {
		return nil, 0, 0, 0, 0, 0, err
	}

	var singleSums []trackerV1Model.SingleSum
	if err := u.DB().WithContext(ctx).Model(&trackerV1Model.SingleSum{}).Where("uid = ?", uid).Find(&singleSums).Error; err != nil {
		return nil, 0, 0, 0, 0, 0, err
	}

	for _, singleSum := range singleSums {
		if singleSum.IsFinish {
			downloadCount++
		}
	}

	keys, err := u.rds.Keys(ctx, "TorrentSeeding/*/"+uid).Result()
	if err != nil {
		return nil, 0, 0, 0, 0, 0, err
	}

	seedingCount = int64(len(keys))
	if err != nil {
		return nil, 0, 0, 0, 0, 0, err
	}

	return user, download, upload, int64(len(torrents)), downloadCount, seedingCount, nil
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

func (u *user) CheckKey(ctx context.Context, key string) (*model.User, error) {
	var _user model.User
	if err := u.DB().WithContext(ctx).Model(&model.User{}).Where("key = ?", key).First(&_user).Error; err != nil {
		return nil, err
	}

	return &_user, nil
}

func (u *user) GetPassKey(ctx context.Context, id string) (key string, err error) {
	var _user model.User
	if err := u.DB().WithContext(ctx).Model(&model.User{}).Where("id = ?", id).First(&_user).Error; err != nil {
		return "", err
	}
	return _user.Key, nil
}
