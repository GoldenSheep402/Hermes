package dao

import (
	"context"
	"errors"
	"github.com/GoldenSheep402/Hermes/mod/system/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type setting struct {
	stdao.Std[*model.Setting]
	rds *redis.Client
}

func (s *setting) Init(db *gorm.DB, rds *redis.Client) error {
	s.rds = rds
	return s.Std.Init(db)
}

func (s *setting) GetPeerExpiration(ctx context.Context) (int, error) {
	key := "PeerExpiration"
	val, err := s.rds.Get(ctx, key).Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			var setting model.Setting
			err := s.Std.DB().WithContext(ctx).Model(&model.Setting{}).First(&setting).Error
			if err != nil {
				return 0, err
			}

			err = s.rds.Set(ctx, key, setting.PeerExpireTime, 24*time.Hour).Err()
			if err != nil {
				return 0, err
			}
			return setting.PeerExpireTime, nil
		default:
			return 0, err
		}
	} else {
		return strconv.Atoi(val)
	}
}

func (s *setting) GetSettings(ctx context.Context) (*model.Setting, []model.Subnet, []model.InnetTracker, error) {
	var setting model.Setting
	err := s.Std.DB().WithContext(ctx).Model(&model.Setting{}).First(&setting).Error
	if err != nil {
		return nil, nil, nil, err
	}

	var subnets []model.Subnet
	err = s.Std.DB().WithContext(ctx).Model(&model.Subnet{}).Find(&subnets).Error
	if err != nil {
		return nil, nil, nil, err
	}

	var innetTrackers []model.InnetTracker
	err = s.Std.DB().WithContext(ctx).Model(&model.InnetTracker{}).Find(&innetTrackers).Error
	if err != nil {
		return nil, nil, nil, err
	}

	return &setting, subnets, innetTrackers, nil
}

func (s *setting) SetSettings(ctx context.Context, setting *model.Setting, subnets []model.Subnet, trackers []model.InnetTracker) error {
	tx := s.DB().WithContext(ctx).Begin()
	defer tx.Rollback()

	var existingSetting model.Setting
	if err := tx.First(&existingSetting).Error; err != nil {
		return err
	}

	if err := tx.Model(&model.Setting{}).Where("id = ?", existingSetting.ID).Select("*").Updates(&setting).Error; err != nil {
		return err
	}

	//var existingSubnets []model.Subnet
	//var ids []string
	//if err := tx.Find(&existingSubnets).Error; err != nil {
	//	return err
	//}

	//var innetMap = make(map[string]bool)

	//for _, subnet := range existingSubnets {
	//	ids = append(ids, subnet.ID)
	//	innetMap[subnet.CIDR] = true
	//}

	//for _, subnet := range subnets {
	//	if subnet.ID == "" {
	//		if _, ok := innetMap[subnet.CIDR]; ok {
	//			return errors.New("subnet already exists")
	//		} else {
	//			err := tx.Model(&model.Subnet{}).Create(&subnet).Error
	//			if err != nil {
	//				return err
	//			}
	//		}
	//	} else {
	//		err := tx.Model(&model.Subnet{}).Where("id = ?", subnet.ID).Select("*").Updates(&subnet).Error
	//		if err != nil {
	//			return err
	//		}
	//	}
	//}

	var existingTrackers []model.InnetTracker

	if err := tx.Find(&existingTrackers).Error; err != nil {
		return err
	}

	var idMap = make(map[string]bool)
	var addrMap = make(map[string]bool)

	for _, tracker := range existingTrackers {
		idMap[tracker.ID] = true
		addrMap[tracker.Address] = true
	}

	for _, tracker := range trackers {
		if tracker.ID == "" {
			if _, ok := addrMap[tracker.Address]; ok {
				return errors.New("address already exists")
			} else {
				err := tx.Model(&model.InnetTracker{}).Create(&tracker).Error
				if err != nil {
					return err
				}
			}
		} else {
			if tracker.DeletedAt.Valid {
				err := tx.Model(&model.InnetTracker{}).Where("id = ?", tracker.ID).Delete(&tracker).Error
				if err != nil {
					return err
				}
			} else {
				err := tx.Model(&model.InnetTracker{}).Where("id = ?", tracker.ID).Select("*").Updates(&tracker).Error
				if err != nil {
					return err
				}
			}
		}
	}

	err := InnetTracker.ClearTrackers(ctx)
	if err != nil {
		return err
	}

	return tx.Commit().Error
}
