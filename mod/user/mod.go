package user

import (
	"context"
	"errors"

	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbac"
	"github.com/GoldenSheep402/Hermes/mod/grpcGateway/gateway"
	"github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/mod/user/service"
	userV1 "github.com/GoldenSheep402/Hermes/pkg/proto/user/v1"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/GoldenSheep402/Hermes/pkg/utils/crypto"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule
	config Config
}

type Config struct {
	AdminName     string `yaml:"AdminName"`
	AdminAccount  string `yaml:"AdminAccount"`
	AdminPassword string `yaml:"AdminPassword"`
}

func (m *Mod) Config() any {
	return &m.config
}

func (m *Mod) Name() string {
	return "user"
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
	userV1.RegisterUserServiceServer(&GRPC, &service.S{
		Log: h.Log.Named("user.service"),
	})
	if err := gw.Register(userV1.RegisterUserServiceHandler); err != nil {
		h.Log.Fatalw("failed to register", "error", err)
	}

	return nil
}

func (m *Mod) Start(h *kernel.Hub) error {
	err := dao.User.DB().Transaction(func(tx *gorm.DB) error {
		// Seed default user group if not exists
		var count int64
		if err := tx.Model(&model.UserGroup{}).Where("name = ?", "default").Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			defaultGroup := &model.UserGroup{
				Name:        "default",
				Description: "Default user group",
				Level:       10,
				CanUpload:   true,
				CanInvite:   false,
			}
			if err := tx.Create(defaultGroup).Error; err != nil {
				return err
			}
		}

		// Seed admin user if not exists
		var adminUser model.User
		if err := tx.Where("is_admin = ?", true).First(&adminUser).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				salt, err := crypto.GenerateSalt(16)
				if err != nil {
					return err
				}
				adminUser = model.User{
					Model:     stdao.Model{ID: ulid.Make().String()},
					Username:  m.config.AdminName,
					Email:     m.config.AdminAccount,
					IsAdmin:   true,
					IsEnabled: true,
					Salt:      salt,
					Password:  crypto.Md5CryptoWithSalt(m.config.AdminPassword, salt),
					Passkey:   ulid.Make().String(),
				}
				if err := tx.Create(&adminUser).Error; err != nil {
					return err
				}
				zap.S().Infow("admin user created", "username", adminUser.Username, "email", adminUser.Email)
			} else {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	if rbac.CasbinManager.Enforcer == nil {
		return errors.New("casbin enforcer is not initialized")
	}

	users, err := dao.User.GetList(context.Background())
	if err != nil {
		return err
	}

	adminCount := 0
	for _, user := range users {
		if user == nil || !user.IsAdmin {
			continue
		}
		if err = rbac.CasbinManager.SetUserGlobalAdmin(user.ID); err != nil {
			return err
		}
		adminCount++
	}

	h.Log.Infow("synced admin users to casbin", "count", adminCount)
	return nil
}
