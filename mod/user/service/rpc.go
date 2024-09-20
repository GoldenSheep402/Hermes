package service

import (
	"context"
	"github.com/GoldenSheep402/Hermes/mod/casbinX/manager"
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

// GetUser retrieves a user's information based on the provided request, ensuring proper authentication and input validation.
func (s *S) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Id is empty")
	}

	// Check the permission
	if req.Id != UID {
		isAdmin, err := userDao.User.IsAdmin(ctx, UID)
		if err != nil {
			return nil, status.Error(codes.Internal, "Internal error")
		}

		if !isAdmin {
			isOk, err := manager.CasbinManager.CheckUserToUserReadPermission(UID, req.Id)
			if err != nil {
				return nil, status.Error(codes.Internal, "Internal error")
			}

			if !isOk {
				return nil, status.Error(codes.PermissionDenied, "Permission denied")
			}
		}
	}

	_user, err := userDao.User.GetInfo(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}
	resp := &userV1.GetUserResponse{
		User: &userV1.User{
			Id:       _user.ID,
			Nickname: _user.Name,
		},
	}

	return resp, nil
}

func (S *S) UpdateUser(ctx context.Context, req *userV1.UpdateUserRequest) (resp *userV1.UpdateUserResponse, err error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if req.User.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Id is empty")
	}

	// Check the permission
	if req.User.Id != UID {
		isAdmin, err := userDao.User.IsAdmin(ctx, UID)
		if err != nil {
			return nil, status.Error(codes.Internal, "Internal error")
		}
		// TODO: rbac
		if !isAdmin {
			return nil, status.Error(codes.PermissionDenied, "Permission denied")
		}
	}

	return resp, nil
}
