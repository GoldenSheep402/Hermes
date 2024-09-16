package category

import (
	"errors"
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/GoldenSheep402/Hermes/mod/category/dao"
	"github.com/GoldenSheep402/Hermes/mod/category/service"
	"github.com/GoldenSheep402/Hermes/mod/grpcGateway/gateway"
	categoryV1 "github.com/GoldenSheep402/Hermes/pkg/proto/category/v1"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule
}

func (m *Mod) Name() string {
	return "category"
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
	categoryV1.RegisterCategoryServiceServer(&GRPC, &service.S{
		Log: h.Log.Named("category.service"),
	})
	err := gw.Register(categoryV1.RegisterCategoryServiceHandler)
	if err != nil {
		h.Log.Fatalw("failed to register", "error", err)
	}

	return nil
}
