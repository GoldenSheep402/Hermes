package service

import (
	"context"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	userV1 "github.com/GoldenSheep402/Hermes/pkg/proto/user/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ userV1.UserServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	userV1.UnimplementedUserServiceServer
}

func (s *S) GetUser(ctx context.Context, req *userV1.GetUserRequest) (resp *userV1.GetUserResponse, err error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Id is empty")
	}

	if req.Id != UID {
		// 	check if the user is admin
	}

	_user, err := userDao.User.GetInfo(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	// TODO: other info
	resp.User.Id = _user.ID
	resp.User.Nickname = _user.Name

	return resp, nil
}
