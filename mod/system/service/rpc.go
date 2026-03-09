package service

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbac"
	"github.com/GoldenSheep402/Hermes/mod/system/dao"
	"github.com/GoldenSheep402/Hermes/mod/system/model"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/oklog/ulid/v2"
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
	list, err := dao.Setting.List(ctx)
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
	val, err := dao.Setting.GetByKey(ctx, req.Key)
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

	for _, setting := range req.Settings {
		if setting.Key == "" {
			continue
		}
		itemType := setting.Type
		if itemType == "" {
			itemType = "string"
		}
		if err := dao.Setting.UpdateOrCreate(ctx, &model.Setting{
			Model: stdao.Model{ID: ulid.Make().String()},
			Key:   setting.Key,
			Value: setting.Value,
			Type:  itemType,
			Desc:  setting.Desc,
		}); err != nil {
			s.Log.Errorw("failed to upsert setting", "key", setting.Key, "err", err)
			return nil, status.Error(codes.Internal, "Failed to save settings")
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

	if err := dao.Setting.DeleteByKey(ctx, req.Key); err != nil {
		return nil, status.Error(codes.Internal, "Failed to delete setting")
	}

	return &systemV1.DeleteSettingResponse{}, nil
}

func (s S) GetSiteStats(ctx context.Context, req *systemV1.GetSiteStatsRequest) (*systemV1.GetSiteStatsResponse, error) {
	// Usually public or requireAuth
	if _, err := requireAuth(ctx); err != nil {
		return nil, err
	}
	// TODO: Replace with actual DB aggregations from respective DAOs once inter-module tracking is stabilized.
	// For now, returning minimal struct to clear the interface.
	return &systemV1.GetSiteStatsResponse{
		TotalUsers:     0,
		TotalTorrents:  0,
		TotalResources: 0,
		TotalTraffic:   0,
		TotalSeeders:   0,
		TotalLeechers:  0,
	}, nil
}
