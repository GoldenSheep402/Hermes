package service

import (
	"context"
	systemDao "github.com/GoldenSheep402/Hermes/mod/system/dao"
	systemModel "github.com/GoldenSheep402/Hermes/mod/system/model"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	systemV1 "github.com/GoldenSheep402/Hermes/pkg/proto/system/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ systemV1.SystemServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	systemV1.UnimplementedSystemServiceServer
}

func (s *S) GetSettings(ctx context.Context, req *systemV1.GetSettingsRequest) (*systemV1.GetSettingsResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	isAdmin, err := userDao.User.IsAdmin(ctx, UID)
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, "Permission Denied")
	}

	setting, subnets, err := systemDao.Setting.GetSettings(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal Error")
	}

	var resp systemV1.GetSettingsResponse
	resp.Settings = &systemV1.Settings{
		PeerExpireTime: int32(setting.PeerExpireTime),
		SmtpEnable:     setting.SmtpEnable,
		SmtpHost:       setting.SmtpHost,
		SmtpPort:       int32(setting.SmtpPort),
		SmtpUser:       setting.SmtpUser,
		SmtpPassword:   setting.SmtpPass,
		RegisterEnable: setting.RegisterEnable,
		LoginEnable:    setting.LoginEnable,
		PublishEnable:  setting.PublishEnable,
	}

	for _, subnet := range subnets {
		resp.Settings.AllowedNets = append(resp.Settings.AllowedNets, subnet.CIDR)
	}

	return &resp, nil
}

func (s *S) SetSettings(ctx context.Context, req *systemV1.SetSettingsRequest) (*systemV1.SetSettingsResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	isAdmin, err := userDao.User.IsAdmin(ctx, UID)
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, "Permission Denied")
	}

	var setting systemModel.Setting
	setting.PeerExpireTime = int(req.Settings.PeerExpireTime)
	setting.SmtpEnable = req.Settings.SmtpEnable
	setting.SmtpHost = req.Settings.SmtpHost
	setting.SmtpPort = int(req.Settings.SmtpPort)
	setting.SmtpUser = req.Settings.SmtpUser
	setting.SmtpPass = req.Settings.SmtpPassword
	setting.RegisterEnable = req.Settings.RegisterEnable
	setting.LoginEnable = req.Settings.LoginEnable
	setting.PublishEnable = req.Settings.PublishEnable

	var subnets []systemModel.Subnet
	for _, cidr := range req.Settings.AllowedNets {
		subnets = append(subnets, systemModel.Subnet{
			CIDR:    cidr,
			IsAllow: true,
		})
	}

	if err := systemDao.Setting.SetSettings(ctx, &setting, subnets); err != nil {
		return nil, status.Error(codes.Internal, "Internal Error")
	}
	return &systemV1.SetSettingsResponse{}, nil
}
