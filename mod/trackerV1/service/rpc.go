package service

import (
	"context"
	trackerV1Dao "github.com/GoldenSheep402/Hermes/mod/trackerV1/dao"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	trackerV1 "github.com/GoldenSheep402/Hermes/pkg/proto/tracker/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ trackerV1.TrackerServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	trackerV1.UnimplementedTrackerServiceServer
}

func (s *S) GetTorrentDownloadingStatus(ctx context.Context, req *trackerV1.GetTorrentDownloadingStatusRequest) (*trackerV1.GetTorrentDownloadingStatusResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}
	// TODO: rbac

	if req.TorrentId == "" {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	}

	downloading, seeding, finished, err := trackerV1Dao.TrackerV1.GetTorrentDownloadingStatus(ctx, req.TorrentId)
	if err != nil {
		return nil, err
	}

	return &trackerV1.GetTorrentDownloadingStatusResponse{
		Downloading: downloading,
		Seeding:     seeding,
		Finished:    finished,
	}, nil

}
