package service

import (
	trackerV1 "github.com/GoldenSheep402/Hermes/pkg/proto/tracker/v1"
	"go.uber.org/zap"
)

var _ trackerV1.TrackerServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	trackerV1.UnimplementedTrackerServiceServer
}
