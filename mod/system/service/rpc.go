package service

import (
	systemV1 "github.com/GoldenSheep402/Hermes/pkg/proto/system/v1"
	"go.uber.org/zap"
)

// TODO: Reimplement system service RPC methods with new KV Setting model.

var _ systemV1.SystemServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	systemV1.UnimplementedSystemServiceServer
}
