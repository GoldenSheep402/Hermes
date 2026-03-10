package service

import (
	trafficV1 "github.com/GoldenSheep402/Hermes/pkg/proto/traffic/v1"
	"math"

	"context"
	"time"

	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbac"
	trackerDao "github.com/GoldenSheep402/Hermes/mod/tracker/dao"
	"github.com/GoldenSheep402/Hermes/mod/traffic/dao"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ trafficV1.TrafficServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	trafficV1.UnimplementedTrafficServiceServer
}

// requireAuth checks if caller is authenticated.
func requireAuth(ctx context.Context) error {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return nil
}

func (s S) ListTransferHistory(ctx context.Context, request *trafficV1.ListTransferHistoryRequest) (*trafficV1.ListTransferHistoryResponse, error) {
	if err := requireAuth(ctx); err != nil {
		return nil, err
	}
	if request.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "User ID required")
	}

	callerID := ctx.Value(ctxKey.UID).(string)
	if request.UserId != callerID {
		isAdmin, err := rbac.CasbinManager.CheckUserIsGlobalAdmin(callerID)
		if err != nil || !isAdmin {
			return nil, status.Error(codes.PermissionDenied, "Not authorized to view other user's transfer history")
		}
	}

	list, count, err := dao.TransferHistory.ListByUserID(ctx, request.UserId, request.Page, request.PageSize)
	if err != nil {
		s.Log.Errorw("failed to list transfer history", "err", err)
		return nil, status.Error(codes.Internal, "Failed to fetch history")
	}

	var pList []*trafficV1.TransferHistoryItem
	for _, item := range list {
		pList = append(pList, &trafficV1.TransferHistoryItem{
			Id:          item.ID,
			TorrentId:   item.TorrentID,
			TorrentName: "Unknown", // Assuming we fetch join later if needed
			Uploaded:    item.Uploaded,
			Downloaded:  item.Downloaded,
			SeedTime:    item.SeedTime,
			IsFinished:  item.IsFinished,
			IsActive:    item.IsActive,
			LastAction:  item.LastAction.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &trafficV1.ListTransferHistoryResponse{
		Items: pList,
		Total: count,
	}, nil
}

func (s S) GetTorrentStats(ctx context.Context, request *trafficV1.GetTorrentStatsRequest) (*trafficV1.GetTorrentStatsResponse, error) {
	if err := requireAuth(ctx); err != nil {
		return nil, err
	}
	if request.TorrentId == "" {
		return nil, status.Error(codes.InvalidArgument, "Torrent ID required")
	}

	ts, err := dao.TorrentStats.GetByTorrentID(ctx, request.TorrentId)
	if err != nil {
		s.Log.Errorw("failed to get torrent stats", "err", err)
		return nil, status.Error(codes.NotFound, "Torrent stats not found")
	}

	return &trafficV1.GetTorrentStatsResponse{
		Stats: &trafficV1.TorrentStatsInfo{
			TorrentId:     ts.TorrentID,
			SeedCount:     int32(ts.SeedCount),
			LeechCount:    int32(ts.LeechCount),
			SnatchCount:   int32(ts.SnatchCount),
			TotalUpload:   ts.TotalUpload,
			TotalDownload: ts.TotalDownload,
		},
	}, nil
}

func (s S) StreamSiteTraffic(req *trafficV1.StreamSiteTrafficRequest, stream trafficV1.TrafficService_StreamSiteTrafficServer) error {
	ctx := stream.Context()
	if err := requireAuth(ctx); err != nil {
		return err
	}

	intervalSeconds := clampInt32(req.GetIntervalSeconds(), 1, 5, 1)
	smoothingFactor := clampFloat64(req.GetSmoothingFactor(), 0.05, 1, 0.35)
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()

	var smoothUpload float64
	var smoothDownload float64
	hasSample := false

	sendPoint := func() error {
		rawUploadRate, rawDownloadRate, err := trackerDao.Traffic.GetSiteRealtimeRate(ctx)
		if err != nil {
			s.Log.Warnw("failed to sample site realtime traffic", "err", err)
			rawUploadRate = 0
			rawDownloadRate = 0
		}

		if !hasSample {
			smoothUpload = float64(rawUploadRate)
			smoothDownload = float64(rawDownloadRate)
			hasSample = true
		} else {
			smoothUpload = smoothingFactor*float64(rawUploadRate) + (1-smoothingFactor)*smoothUpload
			smoothDownload = smoothingFactor*float64(rawDownloadRate) + (1-smoothingFactor)*smoothDownload
		}

		point := &trafficV1.SiteTrafficRatePoint{
			Timestamp:       time.Now().Format(time.RFC3339),
			UploadRate:      int64(math.Round(smoothUpload)),
			DownloadRate:    int64(math.Round(smoothDownload)),
			RawUploadRate:   rawUploadRate,
			RawDownloadRate: rawDownloadRate,
		}

		return stream.Send(point)
	}

	if err := sendPoint(); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := sendPoint(); err != nil {
				return err
			}
		}
	}
}

func clampInt32(value, min, max, fallback int32) int32 {
	if value == 0 {
		return fallback
	}
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func clampFloat64(value, min, max, fallback float64) float64 {
	if value == 0 {
		return fallback
	}
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
