package service

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbac"
	resourceDao "github.com/GoldenSheep402/Hermes/mod/resource/dao"
	resourceModel "github.com/GoldenSheep402/Hermes/mod/resource/model"
	systemSetting "github.com/GoldenSheep402/Hermes/mod/system/setting"
	torrentDao "github.com/GoldenSheep402/Hermes/mod/torrent/dao"
	torrentModel "github.com/GoldenSheep402/Hermes/mod/torrent/model"
	trackerDao "github.com/GoldenSheep402/Hermes/mod/tracker/dao"
	trafficDao "github.com/GoldenSheep402/Hermes/mod/traffic/dao"
	trafficModel "github.com/GoldenSheep402/Hermes/mod/traffic/model"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	userModel "github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	systemV1 "github.com/GoldenSheep402/Hermes/pkg/proto/system/v1"
	"go.uber.org/zap"
)

// TODO: Reimplement system service RPC methods with new KV Setting model.

var _ systemV1.SystemServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	systemV1.UnimplementedSystemServiceServer
}

func requireAuth(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return "", status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return userID, nil
}

func requireAdmin(ctx context.Context) error {
	userID, err := requireAuth(ctx)
	if err != nil {
		return err
	}
	isAdmin, err := rbac.CasbinManager.CheckUserIsGlobalAdmin(userID)
	if err != nil || !isAdmin {
		return status.Error(codes.PermissionDenied, "admin privileges required")
	}
	return nil
}

func (s S) GetSettings(ctx context.Context, req *systemV1.GetSettingsRequest) (*systemV1.GetSettingsResponse, error) {
	if _, err := requireAuth(ctx); err != nil {
		return nil, err
	}
	list, err := systemSetting.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to list settings")
	}

	var pbSettings []*systemV1.Setting
	for _, setting := range list {
		pbSettings = append(pbSettings, &systemV1.Setting{
			Key:   setting.Key,
			Value: setting.Value,
			Type:  setting.Type,
			Desc:  setting.Desc,
		})
	}
	return &systemV1.GetSettingsResponse{Settings: pbSettings}, nil
}

func (s S) GetSetting(ctx context.Context, req *systemV1.GetSettingRequest) (*systemV1.GetSettingResponse, error) {
	if _, err := requireAuth(ctx); err != nil {
		return nil, err
	}
	if req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "Key is required")
	}
	val, err := systemSetting.Get(ctx, req.Key)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Setting not found")
	}
	return &systemV1.GetSettingResponse{
		Setting: &systemV1.Setting{
			Key:   val.Key,
			Value: val.Value,
			Type:  val.Type,
			Desc:  val.Desc,
		},
	}, nil
}

func (s S) SetSettings(ctx context.Context, req *systemV1.SetSettingsRequest) (*systemV1.SetSettingsResponse, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	syncFreeleech := false
	for _, setting := range req.Settings {
		if setting.Key == "" {
			continue
		}
		itemType := setting.Type
		if itemType == "" {
			itemType = "string"
		}
		if err := systemSetting.Upsert(ctx, setting.Key, setting.Value, itemType, setting.Desc); err != nil {
			s.Log.Errorw("failed to upsert setting", "key", setting.Key, "err", err)
			return nil, status.Error(codes.Internal, "Failed to save settings")
		}
		if setting.Key == systemSetting.SettingKeyTrackerGlobalFreeleech ||
			setting.Key == systemSetting.SettingKeyTrackerFreeleechCountdownHrs {
			syncFreeleech = true
		}
	}

	if syncFreeleech {
		if err := trackerDao.Traffic.SyncGlobalFreeleechFromSettings(
			ctx,
			systemSetting.TrackerGlobalFreeleechValue(ctx),
			systemSetting.TrackerFreeleechCountdownHoursValue(ctx),
		); err != nil {
			s.Log.Warnw("failed to sync global freeleech to redis", "err", err)
		}
	}

	return &systemV1.SetSettingsResponse{}, nil
}

func (s S) DeleteSetting(ctx context.Context, req *systemV1.DeleteSettingRequest) (*systemV1.DeleteSettingResponse, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	if req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "Key is required")
	}

	if err := systemSetting.Delete(ctx, req.Key); err != nil {
		return nil, status.Error(codes.Internal, "Failed to delete setting")
	}

	return &systemV1.DeleteSettingResponse{}, nil
}

func (s S) GetSiteStats(ctx context.Context, req *systemV1.GetSiteStatsRequest) (*systemV1.GetSiteStatsResponse, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	if userDao.User.DB() == nil || torrentDao.Torrent.DB() == nil || resourceDao.Resource.DB() == nil || trafficDao.UserTraffic.DB() == nil {
		return nil, status.Error(codes.FailedPrecondition, "service dependencies are not initialized")
	}

	var totalUsers int64
	if err := userDao.User.DB().WithContext(ctx).Model(&userModel.User{}).Count(&totalUsers).Error; err != nil {
		s.Log.Errorw("failed to count users", "err", err)
		return nil, status.Error(codes.Internal, "Failed to aggregate site stats")
	}

	var totalTorrents int64
	if err := torrentDao.Torrent.DB().WithContext(ctx).Model(&torrentModel.Torrent{}).Count(&totalTorrents).Error; err != nil {
		s.Log.Errorw("failed to count torrents", "err", err)
		return nil, status.Error(codes.Internal, "Failed to aggregate site stats")
	}

	var totalResources int64
	if err := resourceDao.Resource.DB().WithContext(ctx).Model(&resourceModel.Resource{}).Count(&totalResources).Error; err != nil {
		s.Log.Errorw("failed to count resources", "err", err)
		return nil, status.Error(codes.Internal, "Failed to aggregate site stats")
	}

	var totalTraffic int64
	if upload, download, found, err := trackerDao.Traffic.GetSiteTotals(ctx); err == nil && found {
		totalTraffic = upload + download
	} else {
		if err != nil {
			s.Log.Warnw("failed to sample tracker site totals", "err", err)
		}
		if err := trafficDao.UserTraffic.DB().WithContext(ctx).
			Model(&trafficModel.UserTraffic{}).
			Select("COALESCE(SUM(real_upload + real_download), 0)").
			Scan(&totalTraffic).Error; err != nil {
			s.Log.Errorw("failed to aggregate total traffic", "err", err)
			return nil, status.Error(codes.Internal, "Failed to aggregate site stats")
		}
	}

	// O(1) DB aggregate: sums per-torrent seed/leech slots from last flush to DB (may lag Redis by flush interval).
	var peerSum struct {
		SeedSum  int64 `gorm:"column:seed_sum"`
		LeechSum int64 `gorm:"column:leech_sum"`
	}
	if err := torrentDao.Torrent.DB().WithContext(ctx).
		Model(&torrentModel.Torrent{}).
		Select("COALESCE(SUM(seed_count), 0) AS seed_sum, COALESCE(SUM(leech_count), 0) AS leech_sum").
		Scan(&peerSum).Error; err != nil {
		s.Log.Errorw("failed to aggregate peer slot sums", "err", err)
		return nil, status.Error(codes.Internal, "Failed to aggregate site stats")
	}

	return &systemV1.GetSiteStatsResponse{
		TotalUsers:     totalUsers,
		TotalTorrents:  totalTorrents,
		TotalResources: totalResources,
		TotalTraffic:   totalTraffic,
		TotalSeeders:   peerSum.SeedSum,
		TotalLeechers:  peerSum.LeechSum,
	}, nil
}
