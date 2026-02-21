package tracker

import (
	"errors"

	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/GoldenSheep402/Hermes/mod/grpcGateway/gateway"
	"github.com/GoldenSheep402/Hermes/mod/tracker/service"
	trackerV1 "github.com/GoldenSheep402/Hermes/pkg/proto/tracker/v1"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule
	config Config
}

type Config struct {
	Endpoint       string   `yaml:"Endpoint"`
	AllowedSubnets []string `yaml:"AllowedSubnets"`
}

func (m *Mod) Config() any {
	return &m.config
}

func (m *Mod) Name() string {
	return "tracker"
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

	// TODO: dao.Init(db, rdb)
	// TODO: register announce/scrape HTTP handlers via jin

	var gw gateway.Gateway
	if h.Load(&gw) != nil {
		return errors.New("can't load gateway from kernel")
	}
	var GRPC grpc.Server
	if h.Load(&GRPC) != nil {
		return errors.New("can't load gRPC server from kernel")
	}
	trackerV1.RegisterTrackerServiceServer(&GRPC, &service.S{
		Log: h.Log.Named("tracker.service"),
	})
	if err := gw.Register(trackerV1.RegisterTrackerServiceHandler); err != nil {
		h.Log.Fatalw("failed to register", "error", err)
	}

	return nil
}
