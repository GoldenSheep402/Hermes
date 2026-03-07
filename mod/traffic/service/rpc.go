package service

import (
	trafficV1 "github.com/GoldenSheep402/Hermes/pkg/proto/traffic/v1"

	"context"

	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbac"
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

func (s S) GetUserTraffic(ctx context.Context, request *trafficV1.GetUserTrafficRequest) (*trafficV1.GetUserTrafficResponse, error) {
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
			return nil, status.Error(codes.PermissionDenied, "Not authorized to view other user's traffic")
		}
	}

	ut, err := dao.UserTraffic.GetByUserID(ctx, request.UserId)
	if err != nil {
		s.Log.Errorw("failed to get user traffic", "err", err)
		return nil, status.Error(codes.NotFound, "User traffic not found")
	}

	var ratio float64
	if ut.RealDownload+ut.BonusDownload > 0 {
		ratio = float64(ut.RealUpload+ut.BonusUpload) / float64(ut.RealDownload+ut.BonusDownload)
	}

	return &trafficV1.GetUserTrafficResponse{
		Traffic: &trafficV1.UserTrafficInfo{
			UserId:        ut.UserID,
			RealUpload:    ut.RealUpload,
			RealDownload:  ut.RealDownload,
			BonusUpload:   ut.BonusUpload,
			BonusDownload: ut.BonusDownload,
			Ratio:         ratio,
		},
	}, nil
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
