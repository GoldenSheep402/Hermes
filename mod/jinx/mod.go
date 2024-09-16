package jinx

import (
	"context"
	"errors"
	"fmt"
	"github.com/juanjiTech/jframe/conf"
	"github.com/juanjiTech/jframe/core/kernel"
	"github.com/juanjiTech/jframe/mod/jinx/healthcheck"
	"github.com/juanjiTech/jin"
	"github.com/juanjiTech/jin/middleware/cors"
	sentryjin "github.com/juanjiTech/sentry-jin"
	"github.com/opentracing/opentracing-go"
	"github.com/soheilhy/cmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net"
	"net/http"
	"sync"
	"time"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule

	listener net.Listener
	j        *jin.Engine
	httpSrv  *http.Server
}

func (m *Mod) Name() string {
	return "jinx"
}

func (m *Mod) Init(hub *kernel.Hub) error {
	m.j = jin.New()
	corsConf := cors.DefaultConfig()
	corsConf.AllowAllOrigins = true
	corsConf.AllowCredentials = true
	corsConf.AddAllowHeaders("Authorization")
	m.j.Use(
		jin.Recovery(),
		cors.New(corsConf),
	)
	if conf.Get().SentryDsn != "" {
		m.j.Use(sentryjin.New(sentryjin.Options{Repanic: true}))
	}
	healthcheck.Register(m.j)

	hub.Map(m.j)
	return nil
}

func (m *Mod) Load(hub *kernel.Hub) error {
	var jinE jin.Engine
	if hub.Load(&jinE) != nil {
		return errors.New("can't load jin.Engine from kernel")
	}
	return nil
}

func (m *Mod) Start(hub *kernel.Hub) error {
	var tcpMux cmux.CMux
	err := hub.Load(&tcpMux)
	if err != nil {
		return errors.New("can't load tcpMux from kernel")
	}

	httpL := tcpMux.Match(cmux.HTTP1Fast())
	m.listener = httpL

	// check if tracer exist
	var tracer opentracing.Tracer
	if hub.Load(&tracer) != nil {
		m.httpSrv = &http.Server{
			Handler: m.j,
		}
	} else {
		fmt.Println("tracer find in kernel, enable http tracing ...")
		m.httpSrv = &http.Server{
			Handler: otelhttp.NewHandler(
				m.j,
				conf.Get().Uptrace.ServiceName,
			),
		}
	}

	if err := m.httpSrv.Serve(httpL); err != nil && !errors.Is(err, http.ErrServerClosed) {
		hub.Log.Infow("failed to start to listen and serve", "error", err)
	}
	return nil
}

func (m *Mod) Stop(wg *sync.WaitGroup, ctx context.Context) error {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := m.httpSrv.Shutdown(ctx); err != nil {
		fmt.Println("Server forced to shutdown: " + err.Error())
		return err
	}
	return nil
}
