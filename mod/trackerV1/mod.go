package trackerV1

import (
	"errors"
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/GoldenSheep402/Hermes/mod/grpcGateway/gateway"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/dao"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/handlers"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/service"
	trackerV1 "github.com/GoldenSheep402/Hermes/pkg/proto/tracker/v1"
	"github.com/juanjiTech/jin"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule
}

func (m *Mod) Name() string {
	return "trackerV1"
}

func (m *Mod) Load(h *kernel.Hub) error {
	var gw gateway.Gateway
	if h.Load(&gw) != nil {
		return errors.New("can't load gateway from kernel")
	}
	var GRPC grpc.Server
	if h.Load(&GRPC) != nil {
		return errors.New("can't load gRPC server from kernel")
	}
	trackerV1.RegisterTrackerServiceServer(&GRPC, &service.S{
		Log: h.Log.Named("trackerV1.service"),
	})
	var db *gorm.DB
	if h.Load(&db) != nil {
		return errors.New("can't load gorm from kernel")
	}

	var rdb *redis.Client
	if h.Load(&rdb) != nil {
		return errors.New("can't load redis client from kernel")
	}

	if err := dao.Init(db, rdb); err != nil {
		return err
	}

	var jinE jin.Engine
	err := h.Load(&jinE)
	if err != nil {
		return errors.New("can't load jin.Engine from kernel")
	}
	handlers.Registry(&jinE)

	err = gw.Register(trackerV1.RegisterTrackerServiceHandler)
	if err != nil {
		h.Log.Fatalw("failed to register", "error", err)
	}

	return nil
}
