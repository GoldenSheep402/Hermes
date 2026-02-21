package service

import (
	trafficV1 "github.com/GoldenSheep402/Hermes/pkg/proto/traffic/v1"

	"context"

	"go.uber.org/zap"
)

var _ trafficV1.TrafficServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	trafficV1.UnimplementedTrafficServiceServer
}

func (s S) GetUserTraffic(ctx context.Context, request *trafficV1.GetUserTrafficRequest) (*trafficV1.GetUserTrafficResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) ListTransferHistory(ctx context.Context, request *trafficV1.ListTransferHistoryRequest) (*trafficV1.ListTransferHistoryResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) GetTorrentStats(ctx context.Context, request *trafficV1.GetTorrentStatsRequest) (*trafficV1.GetTorrentStatsResponse, error) {
	// TODO implement me
	panic("implement me")
}
