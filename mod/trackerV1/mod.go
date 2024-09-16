package trackerV1

import (
	"errors"
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/GoldenSheep402/Hermes/mod/grpcGateway/gateway"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/service"
	trackerV1 "github.com/GoldenSheep402/Hermes/pkg/proto/tracker/v1"
	"google.golang.org/grpc"
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
	err := gw.Register(trackerV1.RegisterTrackerServiceHandler)
	if err != nil {
		h.Log.Fatalw("failed to register", "error", err)
	}

	return nil
}
