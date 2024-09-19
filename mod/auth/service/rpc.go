package service

import (
	"context"
	"errors"
	"github.com/GoldenSheep402/Hermes/mod/user/dao"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/mod/user/model/bindType"
	"github.com/GoldenSheep402/Hermes/pkg/auth"
	authV1 "github.com/GoldenSheep402/Hermes/pkg/proto/auth/v1"
	"github.com/GoldenSheep402/Hermes/pkg/utils/check"
	"github.com/GoldenSheep402/Hermes/pkg/utils/crypto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ authV1.AuthServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	authV1.UnimplementedAuthServiceServer
}

// RegisterSendEmail TODO: SMTP
func (s *S) RegisterSendEmail(ctx context.Context, req *authV1.RegisterSendEmailRequest) (*authV1.RegisterSendEmailResponse, error) {
	return nil, nil
}

// RegisterWithEmail TODO: SMTP
func (s *S) RegisterWithEmail(ctx context.Context, req *authV1.RegisterWithEmailRequest) (*authV1.RegisterWithEmailResponse, error) {
	email := req.Email
	password := req.Password

	if !check.VerifyEmailFormat(email) {
		return nil, status.Error(codes.InvalidArgument, "Email format error")
	}

	if len(password) < 6 {
		return nil, status.Error(codes.InvalidArgument, "Password too short")
	}

	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "Username is empty")
	}

	salt, err := crypto.GenerateSalt(16)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	_, err = dao.User.NewUserWithBind(ctx,
		&model.User{
			Name:     req.Username,
			IsAdmin:  false,
			Salt:     salt,
			Password: crypto.Md5CryptoWithSalt(password, salt),
		},
		&model.Bind{
			OpenID:   email,
			Platform: bindType.Email,
		},
	)

	if err != nil {
		if errors.Is(err, userDao.ErrBindInfoAlreadyUsed) {
			return nil, status.Error(codes.InvalidArgument, "Email already used")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &authV1.RegisterWithEmailResponse{}, nil
}

func (s *S) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	email := req.Email
	password := req.Password

	if !check.VerifyEmailFormat(email) {
		return nil, status.Error(codes.InvalidArgument, "Email format error")
	}

	bind := &model.Bind{
		OpenID: email,
	}

	if result := dao.Bind.DB().Where("open_id = ?", email).First(bind); result.Error != nil {
		return nil, status.Error(codes.InvalidArgument, "Email error")
	}

	user := &model.User{}

	if result := dao.User.DB().Where("id = ?", bind.UID).First(user); result.Error != nil {
		return nil, status.Error(codes.InvalidArgument, "Email error")
	}

	if crypto.Md5CryptoWithSalt(password, user.Salt) != user.Password {
		return nil, status.Error(codes.InvalidArgument, "Password error")
	}

	refreshToken, err := auth.GenToken(auth.Info{
		UID:            bind.UID,
		IsRefreshToken: true,
	}, auth.RefreshTokenExpireIn)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")

	}

	accessToken, err := auth.GenToken(auth.Info{
		UID:            bind.UID,
		IsRefreshToken: false,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &authV1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *S) RefreshToken(_ context.Context, request *authV1.RefreshTokenRequest) (*authV1.RefreshTokenResponse, error) {
	entity, err := auth.ParseToken(request.RefreshToken)
	if err != nil || !entity.Info.IsRefreshToken {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	accessToken, err := auth.GenToken(
		auth.Info{
			// following the same UID
			UID: entity.Info.UID,
		})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authV1.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: request.RefreshToken,
	}, nil
}
