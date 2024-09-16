package uptrace

import (
	"context"
	"errors"
	"github.com/juanjiTech/jframe/conf"
	"github.com/juanjiTech/jframe/core/kernel"
	"github.com/opentracing/opentracing-go"
	"github.com/uptrace/uptrace-go/uptrace"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"sync"
	"time"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule
	closer func(ctx context.Context) error
	tracer opentracing.Tracer
}

func (m *Mod) Name() string {
	return "uptrace"
}

func (m *Mod) Init(hub *kernel.Hub) error {
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(conf.Get().Uptrace.DSN),
		uptrace.WithServiceName(conf.Get().Uptrace.ServiceName),
		uptrace.WithServiceVersion(conf.Get().Uptrace.ServiceVersion),
	)

	m.tracer = opentracing.GlobalTracer()

	if m.tracer == nil {
		hub.Log.Error("failed to initialize tracer")
		return errors.New("failed to initialize tracer")
	}

	hub.Map(&m.tracer)

	m.closer = func(ctx context.Context) error {
		hub.Log.Info("Shutting down uptrace...")
		return uptrace.Shutdown(ctx)
	}

	return nil
}

func (m *Mod) Load(hub *kernel.Hub) error {
	var tracer opentracing.Tracer
	if hub.Load(&tracer) != nil {
		hub.Log.Error("can't load tracer from kernel")
		return errors.New("can't load tracer from kernel")
	}
	opentracing.SetGlobalTracer(m.tracer)

	var db *gorm.DB
	if hub.Load(&db) != nil {
		hub.Log.Info("no gorm find in kernel, skip tracing for gorm")
	} else {
		hub.Log.Info("find gorm in kernel, enable tracing for gorm ...")
		if err := db.Use(tracing.NewPlugin()); err != nil {
			hub.Log.Error("failed to enable tracing for GORM", "error", err)
			return err
		}
	}
	return nil
}

func (m *Mod) Start(hub *kernel.Hub) error {
	return nil
}

func (m *Mod) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := m.closer(ctx); err != nil {
		return err
	}

	return nil
}
