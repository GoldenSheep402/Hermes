package bonus

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/GoldenSheep402/Hermes/mod/bonus/dao"
	"github.com/GoldenSheep402/Hermes/mod/bonus/service"
	"github.com/GoldenSheep402/Hermes/mod/bonus/settle"
	"github.com/GoldenSheep402/Hermes/mod/grpcGateway/gateway"
	bonusV1 "github.com/GoldenSheep402/Hermes/pkg/proto/bonus/v1"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule

	settleCancel context.CancelFunc
	settleWG     sync.WaitGroup
}

func (m *Mod) Name() string {
	return "bonus"
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
	bonusV1.RegisterBonusServiceServer(&GRPC, &service.S{
		Log: h.Log.Named("bonus.service"),
	})
	if err := gw.Register(bonusV1.RegisterBonusServiceHandler); err != nil {
		h.Log.Fatalw("failed to register", "error", err)
	}

	return nil
}

func (m *Mod) Start(h *kernel.Hub) error {
	ctx, cancel := context.WithCancel(context.Background())
	m.settleCancel = cancel

	m.settleWG.Add(1)
	go func() {
		defer m.settleWG.Done()
		log := h.Log.Named("bonus.settle")
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()

		// Run once shortly after boot so fresh deployments settle without waiting a full hour.
		runCtx, runCancel := context.WithTimeout(ctx, 2*time.Minute)
		if n, err := settle.RunOnce(runCtx, log); err != nil {
			log.Warnw("bonus settle failed", "err", err)
		} else if n > 0 {
			log.Infow("bonus settle completed", "pairs", n)
		}
		runCancel()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				runCtx, runCancel := context.WithTimeout(ctx, 5*time.Minute)
				n, err := settle.RunOnce(runCtx, log)
				runCancel()
				if err != nil {
					log.Warnw("bonus settle failed", "err", err)
					continue
				}
				if n > 0 {
					log.Infow("bonus settle completed", "pairs", n)
				}
			}
		}
	}()

	return nil
}

func (m *Mod) Stop(wg *sync.WaitGroup, _ context.Context) error {
	defer wg.Done()
	if m.settleCancel != nil {
		m.settleCancel()
	}
	m.settleWG.Wait()
	return nil
}
