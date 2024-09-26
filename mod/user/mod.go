package user

import (
	"errors"
	"fmt"
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbac"
	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbacValues"
	"github.com/GoldenSheep402/Hermes/mod/grpcGateway/gateway"
	"github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/mod/user/model/bindType"
	"github.com/GoldenSheep402/Hermes/mod/user/service"
	"github.com/GoldenSheep402/Hermes/pkg/colorful"
	userV1 "github.com/GoldenSheep402/Hermes/pkg/proto/user/v1"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/GoldenSheep402/Hermes/pkg/utils/crypto"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"strconv"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule
	config                     Config
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
	err := gw.Register(userV1.RegisterUserServiceHandler)
	if err != nil {
		h.Log.Fatalw("failed to register", "error", err)
	}

	return nil
}

func (m *Mod) Start(h *kernel.Hub) error {
	err := dao.User.DB().Transaction(func(tx *gorm.DB) error {
		// Define preset groups
		presetGroups := []struct {
			Name        string
			Description string
			Level       int
		}{
			{"normalUser", "normal user", rbacValues.Level10},
			{"rbac", "rbac user", rbacValues.Level1},
			{"admin", "admin", rbacValues.Level0},
		}

		// Map to store created group IDs
		idMap := make(map[string]string)

		// Create groups if they don't exist
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
				// Create the group
				if err := tx.Create(&group).Error; err != nil {
					return err
				}

				// Create group metadata
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

				// Add group ID to map
				idMap[pg.Name] = group.ID
			} else {
				// Get the existing group ID and store it
				var group model.Group
				if err := tx.Where("name = ?", pg.Name).First(&group).Error; err != nil {
					return err
				}
				idMap[pg.Name] = group.ID
			}
		}

		// Set subgroup relationships using Casbin
		if err := rbac.CasbinManager.SetSubgroup(idMap["rbac"], idMap["admin"]); err != nil {
			return err
		}
		if err := rbac.CasbinManager.SetSubgroup(idMap["normalUser"], idMap["rbac"]); err != nil {
			return err
		}

		// Check if an admin user exists
		var adminUser model.User
		if err := tx.Where("is_admin = ?", true).First(&adminUser).Error; err != nil {
			// If no admin found, create one
			if errors.Is(err, gorm.ErrRecordNotFound) {
				salt, err := crypto.GenerateSalt(16)
				if err != nil {
					return err
				}

				// Create the admin user
				adminUser = model.User{
					Model:    stdao.Model{ID: ulid.Make().String()},
					Name:     m.config.AdminName,
					IsAdmin:  true,
					Salt:     salt,
					Password: crypto.Md5CryptoWithSalt(m.config.AdminPassword, salt),
					Key:      ulid.Make().String(),
				}
				if err := tx.Create(&adminUser).Error; err != nil {
					return err
				}

				fmt.Printf(colorful.Blue("Admin account:") + colorful.Blue(m.config.AdminAccount) + "\n")
				fmt.Printf(colorful.Blue("Admin password:") + colorful.Blue(m.config.AdminPassword))

				// Bind the admin account to email
				adminBind := model.Bind{
					OpenID:   m.config.AdminAccount,
					Platform: bindType.Email,
					UID:      adminUser.ID,
				}

				// Check if the email is already bound
				if err := tx.Where("open_id = ?", adminBind.OpenID).First(&model.Bind{}).Error; err == nil {
					return fmt.Errorf("admin account already exists")
				} else if !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}

				// Create the bind record
				if err := tx.Create(&adminBind).Error; err != nil {
					return err
				}

				var memberShip model.GroupMembership
				memberShip.UID = adminUser.ID
				memberShip.GID = idMap["admin"]
				if err := tx.Model(&model.GroupMembership{}).Create(&memberShip).Error; err != nil {
					return err
				}

				// Add admin user to the admin group
				if err := rbac.CasbinManager.SetUserToGroup(adminUser.ID, idMap["admin"]); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
