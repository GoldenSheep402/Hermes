package user

import (
	"errors"
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/GoldenSheep402/Hermes/mod/grpcGateway/gateway"
	"github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/mod/user/service"
	userV1 "github.com/GoldenSheep402/Hermes/pkg/proto/user/v1"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule
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
