package gateway

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Gateway struct {
	Mux  *runtime.ServeMux
	Conn *grpc.ClientConn
}

func (g *Gateway) Register(handlers ...func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error) error {
	for _, handler := range handlers {
		if err := handler(context.Background(), g.Mux, g.Conn); err != nil {
			return err
		}
	}
	return nil
}
