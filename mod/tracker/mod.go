package tracker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/juanjiTech/jin"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/GoldenSheep402/Hermes/mod/grpcGateway/gateway"
	systemSetting "github.com/GoldenSheep402/Hermes/mod/system/setting"
	"github.com/GoldenSheep402/Hermes/mod/tracker/dao"
	"github.com/GoldenSheep402/Hermes/mod/tracker/handlers"
	"github.com/GoldenSheep402/Hermes/mod/tracker/service"
	trackerV1 "github.com/GoldenSheep402/Hermes/pkg/proto/tracker/v1"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule
	config Config

	trafficFlushCancel context.CancelFunc
	trafficFlushWG     sync.WaitGroup
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

	var jinE *jin.Engine
	if h.Load(&jinE) != nil {
		return errors.New("can't load jin engine from kernel")
	}

	if err := dao.Init(db, rdb); err != nil {
		h.Log.Fatalw("failed to init dao", "error", err)
	}

	handlers.Registry(jinE)

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

func (m *Mod) Start(h *kernel.Hub) error {
	interval, batchSize := m.resolveFlushConfig(context.Background())

	ctx, cancel := context.WithCancel(context.Background())
	m.trafficFlushCancel = cancel

	m.trafficFlushWG.Add(1)
	go func() {
		defer m.trafficFlushWG.Done()

		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				newInterval, newBatchSize := m.resolveFlushConfig(ctx)
				if newInterval != interval {
					interval = newInterval
					ticker.Reset(time.Duration(interval) * time.Second)
					h.Log.Infow("tracker flush interval changed", "seconds", interval)
				}
				batchSize = newBatchSize

				flushed, err := dao.Traffic.FlushPending(ctx, batchSize)
				if err != nil {
					h.Log.Warnw("failed to flush tracker traffic", "err", err, "batch_size", batchSize)
					continue
				}
				if flushed > 0 {
					h.Log.Debugw("tracker traffic flushed", "records", flushed)
				}
			}
		}
	}()

	return nil
}

func (m *Mod) resolveFlushConfig(ctx context.Context) (int, int) {
	return systemSetting.TrackerFlushIntervalValue(ctx), systemSetting.TrackerFlushBatchSizeValue(ctx)
}

func (m *Mod) Stop(wg *sync.WaitGroup, _ context.Context) error {
	defer wg.Done()

	if m.trafficFlushCancel != nil {
		m.trafficFlushCancel()
	}
	m.trafficFlushWG.Wait()

	flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = dao.Traffic.FlushPending(flushCtx, 5000)

	return nil
}
