package auth

import (
	"errors"
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/GoldenSheep402/Hermes/mod/auth/service"
	"github.com/GoldenSheep402/Hermes/mod/casbinX/manager"
	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbacValues"
	"github.com/GoldenSheep402/Hermes/mod/grpcGateway/gateway"
	"github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	authV1 "github.com/GoldenSheep402/Hermes/pkg/proto/auth/v1"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"strconv"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule
}

func (m *Mod) Name() string {
	return "auth"
}

func (m *Mod) Load(h *kernel.Hub) error {
	var rdb *redis.Client
	if h.Load(&rdb) != nil {
		return errors.New("can't load redis client from kernel")
	}
	var db *gorm.DB
	if h.Load(&db) != nil {
		return errors.New("can't load gorm from kernel")
	}
	if err := dao.Init(db, rdb); err != nil {
		return err
	}

	var gw gateway.Gateway
	if h.Load(&gw) != nil {
		return errors.New("can't load gateway from kernel")
	}
	var GRPC grpc.Server
	if h.Load(&GRPC) != nil {
		return errors.New("can't load gRPC server from kernel")
	}

	authV1.RegisterAuthServiceServer(&GRPC, &service.S{
		Log: h.Log.Named("auth.service"),
	})
	err := gw.Register(authV1.RegisterAuthServiceHandler)
	if err != nil {
		h.Log.Fatalw("failed to register", "error", err)
	}

	return nil
}

func (m *Mod) Start(h *kernel.Hub) error {
	err := dao.User.DB().Transaction(func(tx *gorm.DB) error {
		presetGroups := []struct {
			Name        string
			Description string
			Level       int
		}{
			{
				Name:        "normalUser",
				Description: "normal user",
				Level:       rbacValues.Level10,
			},
			{
				Name:        "manager",
				Description: "manager user",
				Level:       rbacValues.Level1,
			},
			{
				Name:        "admin",
				Description: "admin",
				Level:       rbacValues.Level0,
			},
		}

		idMap := make(map[string]string)
		for _, pg := range presetGroups {
			var count int64
			if err := tx.Model(&model.Group{}).Where("name = ?", pg.Name).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				group := model.Group{
					Name:        pg.Name,
					Description: pg.Description,
				}
				if err := tx.Create(&group).Error; err != nil {
					return err
				}
				groupMetadata := model.GroupMetadata{
					GID:         group.ID,
					Key:         "Level",
					Value:       strconv.Itoa(pg.Level),
					Type:        "number",
					Description: "It is the level of the group.",
					Order:       0,
				}
				if err := tx.Create(&groupMetadata).Error; err != nil {
					return err
				}
				idMap[pg.Name] = group.ID
			}
		}

		err := manager.CasbinManager.SetSubgroup(idMap["admin"], idMap["manager"])
		if err != nil {
			return err
		}

		err = manager.CasbinManager.SetSubgroup(idMap["manager"], idMap["normalUser"])
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
