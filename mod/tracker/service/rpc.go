package service

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbac"
	"github.com/GoldenSheep402/Hermes/mod/tracker/dao"
	"github.com/GoldenSheep402/Hermes/mod/tracker/model"
	"github.com/GoldenSheep402/Hermes/pkg/authz"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	trackerV1 "github.com/GoldenSheep402/Hermes/pkg/proto/tracker/v1"
)

var _ trackerV1.TrackerServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	trackerV1.UnimplementedTrackerServiceServer
}

// requireAuth checks if caller is authenticated.
func requireAuth(ctx context.Context) error {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return nil
}

func (s *S) GetTorrentPeers(ctx context.Context, req *trackerV1.GetTorrentPeersRequest) (*trackerV1.GetTorrentPeersResponse, error) {
	if req.TorrentId == "" {
		return nil, status.Error(codes.InvalidArgument, "Torrent ID required")
	}
	if err := authz.RequireTorrentSensitiveAccess(ctx, req.TorrentId); err != nil {
		return nil, err
	}

	peers, err := dao.Peer.GetPeersForTorrent(ctx, req.TorrentId, 0)
	if err != nil {
		s.Log.Errorw("failed to get torrent peers", "err", err)
		return nil, status.Error(codes.Internal, "Failed to fetch peers")
	}
	stats, err := dao.Traffic.GetTorrentSnapshot(ctx, req.TorrentId)
	if err != nil {
		s.Log.Errorw("failed to get torrent peer counts", "err", err)
		return nil, status.Error(codes.Internal, "Failed to fetch peer counts")
	}

	var protoSeeders, protoLeechers []*trackerV1.PeerInfo

	for _, p := range peers {
		info := &trackerV1.PeerInfo{
			PeerId: []byte(p.PeerID),
			Ip:     p.IP,
			Port:   int32(p.Port),
			LanIp:  p.LanIP,
		}
		if p.IsSeeder {
			protoSeeders = append(protoSeeders, info)
		} else {
			protoLeechers = append(protoLeechers, info)
		}
	}

	return &trackerV1.GetTorrentPeersResponse{
		SeederCount:  int32(stats.SeedCount),
		LeecherCount: int32(stats.LeechCount),
		Seeders:      protoSeeders,
		Leechers:     protoLeechers,
	}, nil
}

func (s *S) ListSnatches(ctx context.Context, req *trackerV1.ListSnatchesRequest) (*trackerV1.ListSnatchesResponse, error) {
	if req.TorrentId == "" {
		return nil, status.Error(codes.InvalidArgument, "Torrent ID required")
	}
	if err := authz.RequireTorrentSensitiveAccess(ctx, req.TorrentId); err != nil {
		return nil, err
	}

	list, count, err := dao.Snatch.ListByTorrent(ctx, req.TorrentId, req.Page, req.PageSize)
	if err != nil {
		s.Log.Errorw("failed to list snatches", "err", err)
		return nil, status.Error(codes.Internal, "Failed to list snatches")
	}

	return &trackerV1.ListSnatchesResponse{
		Snatches: convertSnatches(list),
		Total:    count,
	}, nil
}

func (s *S) GetUserSnatches(ctx context.Context, req *trackerV1.GetUserSnatchesRequest) (*trackerV1.GetUserSnatchesResponse, error) {
	if err := requireAuth(ctx); err != nil {
		return nil, err
	}

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "User ID required")
	}

	callerID := ctx.Value(ctxKey.UID).(string)
	if req.UserId != callerID {
		isAdmin, err := rbac.CasbinManager.CheckUserIsGlobalAdmin(callerID)
		if err != nil || !isAdmin {
			return nil, status.Error(codes.PermissionDenied, "Not authorized to view other user's snatches")
		}
	}

	list, count, err := dao.Snatch.ListByUser(ctx, req.UserId, req.Page, req.PageSize)
	if err != nil {
		s.Log.Errorw("failed to get user snatches", "err", err)
		return nil, status.Error(codes.Internal, "Failed to get user snatches")
	}

	return &trackerV1.GetUserSnatchesResponse{
		Snatches: convertSnatches(list),
		Total:    count,
	}, nil
}

func convertSnatches(list []model.Snatch) []*trackerV1.SnatchInfo {
	var pList []*trackerV1.SnatchInfo
	for _, sn := range list {
		var finishedAt string
		if sn.FinishedAt != nil {
			finishedAt = sn.FinishedAt.Format(time.RFC3339)
		}
		pList = append(pList, &trackerV1.SnatchInfo{
			Id:         sn.ID,
			TorrentId:  sn.TorrentID,
			UserId:     sn.UserID,
			Uploaded:   sn.Uploaded,
			Downloaded: sn.Downloaded,
			SeedTime:   sn.SeedTime,
			IsActive:   sn.IsActive,
			FinishedAt: finishedAt,
			LastAction: sn.LastAction.Format(time.RFC3339),
		})
	}
	return pList
}
