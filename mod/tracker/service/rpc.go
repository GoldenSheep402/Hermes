package service

import (
	trackerV1 "github.com/GoldenSheep402/Hermes/pkg/proto/tracker/v1"

	"context"

	"go.uber.org/zap"
)

var _ trackerV1.TrackerServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	trackerV1.UnimplementedTrackerServiceServer
}

func (s S) GetTorrentPeers(ctx context.Context, request *trackerV1.GetTorrentPeersRequest) (*trackerV1.GetTorrentPeersResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) ListSnatches(ctx context.Context, request *trackerV1.ListSnatchesRequest) (*trackerV1.ListSnatchesResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) GetUserSnatches(ctx context.Context, request *trackerV1.GetUserSnatchesRequest) (*trackerV1.GetUserSnatchesResponse, error) {
	// TODO implement me
	panic("implement me")
}
