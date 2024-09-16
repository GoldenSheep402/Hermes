package grpcGateway

import (
	"context"
	"errors"
	"fmt"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcCtxTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcOpentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/juanjiTech/jframe/conf"
	"github.com/juanjiTech/jframe/core/kernel"
	"github.com/juanjiTech/jframe/core/logx"
	"github.com/juanjiTech/jframe/mod/grpcGateway/gateway"
	"github.com/juanjiTech/jframe/mod/grpcGateway/middleware"
	"github.com/juanjiTech/jin"
	"github.com/opentracing/opentracing-go"
	"github.com/soheilhy/cmux"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"net"
	"strings"
	"sync"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule

	grpcL net.Listener
	grpc  *grpc.Server
	gw    *gateway.Gateway
}

func (m *Mod) Name() string {
	return "grpcGateway"
}

func (m *Mod) PreInit(hub *kernel.Hub) error {
	grpcZap.ReplaceGrpcLoggerV2(logx.NameSpace("grpc").Desugar())
	m.grpc = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				grpcCtxTags.UnaryServerInterceptor(),
				grpcOpentracing.UnaryServerInterceptor(),
				grpcZap.UnaryServerInterceptor(logx.NameSpace("grpc").Desugar()),
				grpcRecovery.UnaryServerInterceptor(),
				grpcAuth.UnaryServerInterceptor(middleware.AuthInterceptor),
			),
		),
	)
	reflection.Register(m.grpc)
	hub.Log.Info("init gRPC server success...")
	hub.Map(m.grpc)
	return nil
}

func (m *Mod) PostInit(h *kernel.Hub) error {
	var tcpMux cmux.CMux
	err := h.Load(&tcpMux)
	if err != nil {
		return errors.New("can't load tcpMux from kernel")
	}

	m.grpcL = tcpMux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))

	// 开始初始化grpc-gateway
	opts := []grpc.DialOption{
		//grpc.WithTimeout(10 * time.Second),
		//grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// if tracer can be found from kernel, enable tracing for gRPC gateway
	var tracer opentracing.Tracer
	if h.Load(&tracer) != nil {
		h.Log.Info("no tracer find from kernel, skip tracing for gRPC gateway")
	} else {
		h.Log.Info("tracer find from kernel, enable tracing for gRPC gateway ...")
		h.Log.Info("tracer find from kernel, set StatusHandler for gRPC server ...")
		// if tracer can be found, register tracing middleware for gRPC server
		m.grpc = grpc.NewServer(
			grpc.UnaryInterceptor(
				grpcMiddleware.ChainUnaryServer(
					grpcCtxTags.UnaryServerInterceptor(),
					grpcOpentracing.UnaryServerInterceptor(),
					grpcZap.UnaryServerInterceptor(logx.NameSpace("grpc").Desugar()),
					grpcRecovery.UnaryServerInterceptor(),
					grpcAuth.UnaryServerInterceptor(middleware.AuthInterceptor),
				),
			),
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
		)
		reflection.Register(m.grpc)
		h.Map(m.grpc)
		opts = append(opts, grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	}

	conn, err := grpc.NewClient(fmt.Sprintf("127.0.0.1:%s", conf.Get().Port), opts...)
	if err != nil {
		h.Log.Fatal("gRPC fail to dial: %v", err)
	}

	var allowedHeaders = map[string]struct{}{
		"x-request-id": {}, // 还没用到 后续做追踪
	}
	outHeaderFilter := func(s string) (string, bool) {
		if _, isAllowed := allowedHeaders[s]; isAllowed {
			return strings.ToUpper(s), true
		}
		return s, false
	}

	mux := runtime.NewServeMux(
		runtime.WithOutgoingHeaderMatcher(outHeaderFilter),
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
	)

	var http jin.Engine
	err = h.Load(&http)
	if err != nil {
		return errors.New("can't load jin from kernel")
	}

	http.Any("/gapi/*any", mux.ServeHTTP)
	m.gw = &gateway.Gateway{
		Mux:  mux,
		Conn: conn,
	}
	h.Log.Info("init gRPC gateway success...")
	h.Map(m.gw)
	return nil
}

func (m *Mod) Start(h *kernel.Hub) error {
	// 初始化grpc
	go func() {
		if err := m.grpc.Serve(m.grpcL); err != nil {
			h.Log.Infow("failed to start to listen and serve", "error", err)
		}
	}()
	return nil
}

func (m *Mod) Stop(wg *sync.WaitGroup, _ context.Context) error {
	defer wg.Done()
	m.grpc.GracefulStop()
	return nil
}
