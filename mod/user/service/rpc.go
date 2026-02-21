package service

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/mod/user/model"
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

// helper to convert model.User to userV1.User
func convertUserModelToProto(u *model.User) *userV1.User {
	if u == nil {
		return nil
	}
	protoUser := &userV1.User{
		Id:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		Avatar:      u.Avatar,
		IsAdmin:     u.IsAdmin,
		IsEnabled:   u.IsEnabled,
		GroupId:     u.GroupID,
		BonusPoints: u.BonusPoints,
		Uploaded:    u.Uploaded,
		Downloaded:  u.Downloaded,
		SeedTime:    u.SeedTime,
		InviteCount: int32(u.InviteCount),
		Passkey:     u.Passkey,
		CreatedAt:   u.CreatedAt.String(),
	}
	if u.LastLogin != nil {
		protoUser.LastLogin = u.LastLogin.String()
	}
	return protoUser
}

func (s *S) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	// Assumes Auth interceptor has placed UserID
	userIDVal := ctx.Value("UserID")
	if userIDVal == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	userID := userIDVal.(string)

	user, err := dao.User.GetByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &userV1.GetUserResponse{
		User: convertUserModelToProto(user),
	}, nil
}

func (s *S) GetUserProfile(ctx context.Context, req *userV1.GetUserProfileRequest) (*userV1.GetUserProfileResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	user, err := dao.User.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Calculate ratio
	var ratio double = 0
	if user.Downloaded > 0 {
		ratio = float64(user.Uploaded) / float64(user.Downloaded)
	} else if user.Uploaded > 0 {
		ratio = -1 // infinity
	}

	// TODO: Fetch real upload/download records and counts from TransferHistory
	// For now, return basic user profile stats
	return &userV1.GetUserProfileResponse{
		User:           convertUserModelToProto(user),
		RealUpload:     user.Uploaded,
		RealDownload:   user.Downloaded,
		Ratio:          ratio,
		PublishedCount: 0,
		SeedingCount:   0,
		DownloadCount:  0,
	}, nil
}

func (s *S) UpdateUser(ctx context.Context, req *userV1.UpdateUserRequest) (*userV1.UpdateUserResponse, error) {
	userIDVal := ctx.Value("UserID")
	if userIDVal == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	userID := userIDVal.(string)

	user, err := dao.User.GetByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := dao.User.UpdateInfo(ctx, user); err != nil {
		s.Log.Errorw("failed to update user info", "error", err)
		return nil, status.Error(codes.Internal, "failed to update user info")
	}

	return &userV1.UpdateUserResponse{
		User: convertUserModelToProto(user),
	}, nil
}

func (s *S) UpdatePassword(ctx context.Context, req *userV1.UpdatePasswordRequest) (*userV1.UpdatePasswordResponse, error) {
	// TODO: implement password verification algorithm
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *S) ResetPasskey(ctx context.Context, req *userV1.ResetPasskeyRequest) (*userV1.ResetPasskeyResponse, error) {
	// TODO: implement passkey generation
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (s *S) GetUserPasskey(ctx context.Context, req *userV1.GetUserPasskeyRequest) (*userV1.GetUserPasskeyResponse, error) {
	userIDVal := ctx.Value("UserID")
	if userIDVal == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	userID := userIDVal.(string)

	user, err := dao.User.GetByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &userV1.GetUserPasskeyResponse{
		Passkey: user.Passkey,
	}, nil
}

func (s *S) ListUsers(ctx context.Context, req *userV1.ListUsersRequest) (*userV1.ListUsersResponse, error) {
	users, err := dao.User.GetList(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	var protoUsers []*userV1.User
	for _, u := range users {
		protoUsers = append(protoUsers, convertUserModelToProto(u))
	}

	return &userV1.ListUsersResponse{
		Users: protoUsers,
		Total: int64(len(protoUsers)),
	}, nil
}

// UserGroup methods are stubs for now
func (s *S) CreateUserGroup(ctx context.Context, req *userV1.CreateUserGroupRequest) (*userV1.CreateUserGroupResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (s *S) GetUserGroup(ctx context.Context, req *userV1.GetUserGroupRequest) (*userV1.GetUserGroupResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (s *S) ListUserGroups(ctx context.Context, req *userV1.ListUserGroupsRequest) (*userV1.ListUserGroupsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (s *S) UpdateUserGroup(ctx context.Context, req *userV1.UpdateUserGroupRequest) (*userV1.UpdateUserGroupResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (s *S) DeleteUserGroup(ctx context.Context, req *userV1.DeleteUserGroupRequest) (*userV1.DeleteUserGroupResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
