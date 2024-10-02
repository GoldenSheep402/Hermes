package system

import (
	"errors"
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/GoldenSheep402/Hermes/mod/grpcGateway/gateway"
	"github.com/GoldenSheep402/Hermes/mod/system/dao"
	"github.com/GoldenSheep402/Hermes/mod/system/model"
	"github.com/GoldenSheep402/Hermes/mod/system/service"
	systemV1 "github.com/GoldenSheep402/Hermes/pkg/proto/system/v1"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule
}

func (m *Mod) Name() string {
	return "system"
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

	// if the system is empty, create a default one
	err := db.Model(&model.Setting{}).First(&model.Setting{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&model.Setting{
				PeerExpireTime: 3600,
				SmtpEnable:     false,
				SmtpHost:       "",
				SmtpPort:       0,
				SmtpUser:       "",
				SmtpPass:       "",
				RegisterEnable: true,
				LoginEnable:    true,
				PublishEnable:  true,
			}).Error; err != nil {
				return err
			}
		}
	}

	var gw gateway.Gateway
	if h.Load(&gw) != nil {
		return errors.New("can't load gateway from kernel")
	}
	var GRPC grpc.Server
	if h.Load(&GRPC) != nil {
		return errors.New("can't load gRPC server from kernel")
	}
	systemV1.RegisterSystemServiceServer(&GRPC, &service.S{
		Log: h.Log.Named("system.service"),
	})

	err = gw.Register(systemV1.RegisterSystemServiceHandler)
	if err != nil {
		h.Log.Fatalw("failed to register", "error", err)
	}

	return nil
}
