package service

import (
	"context"

	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbac"
	resourceDao "github.com/GoldenSheep402/Hermes/mod/resource/dao"
	trafficDao "github.com/GoldenSheep402/Hermes/mod/traffic/dao"
	"github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	userV1 "github.com/GoldenSheep402/Hermes/pkg/proto/user/v1"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/GoldenSheep402/Hermes/pkg/utils/crypto"
)

var _ userV1.UserServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	userV1.UnimplementedUserServiceServer
}

func (s *S) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	user, err := dao.User.GetByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &userV1.GetUserResponse{
		User: convertUserModelToProto(user),
	}, nil
}

func (s *S) GetUserProfile(ctx context.Context, req *userV1.GetUserProfileRequest) (*userV1.GetUserProfileResponse, error) {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	id := req.Id
	if req.Id == "" {
		id = userID
	}

	// Permission logic: Can only view other profiles if you are an Admin, or if the system allows high-level users.
	// We will query the current user's DB record to check their capabilities.
	if id != userID {
		currentUser, err := dao.User.GetByID(ctx, userID)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "user not found")
		}

		// Get currentUser group to check level
		if !currentUser.IsAdmin {
			// TODO: group manager
			// For PT sites, normally only members above a certain level or Admins can see peers.
			// group, err := dao.UserGroup.Get(ctx, currentUser.GroupID)
			// if err != nil || group.Level < 1 {
			return nil, status.Error(codes.PermissionDenied, "you need a higher user group level (>= 1) to view other user profiles")
			// }
		}
	}

	user, err := dao.User.GetByID(ctx, id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Fetch real traffic data
	traffic, err := trafficDao.UserTraffic.GetByUserID(ctx, id)
	var realUpload, realDownload int64
	var ratio float64
	if err == nil && traffic != nil {
		realUpload = traffic.RealUpload
		realDownload = traffic.RealDownload
		if realDownload > 0 {
			ratio = float64(realUpload) / float64(realDownload)
		} else if realUpload > 0 {
			ratio = -1 // infinity
		}
	} else {
		// Fallback to basic user struct if traffic record not initialized yet
		realUpload = user.Uploaded
		realDownload = user.Downloaded
		if realDownload > 0 {
			ratio = float64(realUpload) / float64(realDownload)
		} else if realUpload > 0 {
			ratio = -1 // infinity
		}
	}

	// Fetch real activity counts
	publishedCount, _ := resourceDao.Resource.CountPublishedByUser(ctx, id)
	seedingCount, _ := trafficDao.TransferHistory.CountActiveSeeding(ctx, id)
	downloadCount, _ := trafficDao.TransferHistory.CountActiveDownloading(ctx, id)

	return &userV1.GetUserProfileResponse{
		User:           convertUserModelToProto(user),
		RealUpload:     realUpload,
		RealDownload:   realDownload,
		Ratio:          ratio,
		PublishedCount: int32(publishedCount),
		SeedingCount:   int32(seedingCount),
		DownloadCount:  int32(downloadCount),
	}, nil
}

func (s *S) UpdateUser(ctx context.Context, req *userV1.UpdateUserRequest) (*userV1.UpdateUserResponse, error) {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

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
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	user, err := dao.User.GetByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if crypto.Md5CryptoWithSalt(req.OldPassword, user.Salt) != user.Password {
		return nil, status.Error(codes.InvalidArgument, "incorrect old password")
	}

	salt, err := crypto.GenerateSalt(16)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate salt")
	}

	user.Salt = salt
	user.Password = crypto.Md5CryptoWithSalt(req.NewPassword, salt)

	if err := dao.User.UpdateInfo(ctx, user); err != nil {
		s.Log.Errorw("failed to update password", "error", err)
		return nil, status.Error(codes.Internal, "failed to update password")
	}

	return &userV1.UpdatePasswordResponse{}, nil
}

func (s *S) ResetPasskey(ctx context.Context, req *userV1.ResetPasskeyRequest) (*userV1.ResetPasskeyResponse, error) {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	user, err := dao.User.GetByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	user.Passkey = ulid.Make().String()

	if err := dao.User.UpdateInfo(ctx, user); err != nil {
		s.Log.Errorw("failed to reset passkey", "error", err)
		return nil, status.Error(codes.Internal, "failed to reset passkey")
	}

	return &userV1.ResetPasskeyResponse{
		Passkey: user.Passkey,
	}, nil
}

func (s *S) GetUserPasskey(ctx context.Context, req *userV1.GetUserPasskeyRequest) (*userV1.GetUserPasskeyResponse, error) {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	user, err := dao.User.GetByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &userV1.GetUserPasskeyResponse{
		Passkey: user.Passkey,
	}, nil
}

func (s *S) ListUsers(ctx context.Context, req *userV1.ListUsersRequest) (*userV1.ListUsersResponse, error) {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	isAdmin, _ := rbac.CasbinManager.CheckUserIsGlobalAdmin(userID)
	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, "admin needed to list users")
	}
	users, total, err := dao.User.GetListPaginated(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	var protoUsers []*userV1.User
	for _, u := range users {
		protoUsers = append(protoUsers, convertUserModelToProto(u))
	}

	return &userV1.ListUsersResponse{
		Users: protoUsers,
		Total: total,
	}, nil
}

func convertGroupProtoToModel(p *userV1.UserGroup) *model.UserGroup {
	if p == nil {
		return nil
	}
	return &model.UserGroup{
		Model:           stdao.Model{ID: p.Id},
		Name:            p.Name,
		Description:     p.Description,
		Level:           int(p.Level),
		MinUpload:       p.MinUpload,
		MinRatio:        p.MinRatio,
		MinSeedTime:     p.MinSeedTime,
		MaxDownloads:    int(p.MaxDownloads),
		CanUpload:       p.CanUpload,
		CanInvite:       p.CanInvite,
		IsImmuneToRatio: p.IsImmuneToRatio,
		Color:           p.Color,
		Icon:            p.Icon,
	}
}

func convertGroupModelToProto(g *model.UserGroup) *userV1.UserGroup {
	if g == nil {
		return nil
	}
	return &userV1.UserGroup{
		Id:              g.ID,
		Name:            g.Name,
		Description:     g.Description,
		Level:           int32(g.Level),
		MinUpload:       g.MinUpload,
		MinRatio:        g.MinRatio,
		MinSeedTime:     g.MinSeedTime,
		MaxDownloads:    int32(g.MaxDownloads),
		CanUpload:       g.CanUpload,
		CanInvite:       g.CanInvite,
		IsImmuneToRatio: g.IsImmuneToRatio,
		Color:           g.Color,
		Icon:            g.Icon,
	}
}

// UserGroup methods
func (s *S) CreateUserGroup(ctx context.Context, req *userV1.CreateUserGroupRequest) (*userV1.CreateUserGroupResponse, error) {
	// Manual Admin Check
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	isAdmin, _ := rbac.CasbinManager.CheckUserIsGlobalAdmin(userID)
	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, "admin needed")
	}

	modelGroup := convertGroupProtoToModel(req.Group)
	modelGroup.ID = ulid.Make().String() // Ensure unique ID

	if err := dao.UserGroup.Create(ctx, modelGroup); err != nil {
		return nil, status.Error(codes.Internal, "failed to create group")
	}

	return &userV1.CreateUserGroupResponse{}, nil
}

func (s *S) GetUserGroup(ctx context.Context, req *userV1.GetUserGroupRequest) (*userV1.GetUserGroupResponse, error) {
	if _, ok := ctx.Value(ctxKey.UID).(string); !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	group, err := dao.UserGroup.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "group not found")
	}
	return &userV1.GetUserGroupResponse{
		Group: convertGroupModelToProto(group),
	}, nil
}

func (s *S) ListUserGroups(ctx context.Context, req *userV1.ListUserGroupsRequest) (*userV1.ListUserGroupsResponse, error) {
	if _, ok := ctx.Value(ctxKey.UID).(string); !ok {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	groups, err := dao.UserGroup.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list groups")
	}

	var protoGroups []*userV1.UserGroup
	for _, g := range groups {
		protoGroups = append(protoGroups, convertGroupModelToProto(g))
	}

	return &userV1.ListUserGroupsResponse{
		Groups: protoGroups,
	}, nil
}

func (s *S) UpdateUserGroup(ctx context.Context, req *userV1.UpdateUserGroupRequest) (*userV1.UpdateUserGroupResponse, error) {
	// Manual Admin Check
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	isAdmin, _ := rbac.CasbinManager.CheckUserIsGlobalAdmin(userID)
	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, "admin needed")
	}

	if req.Group == nil || req.Group.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "group id is required")
	}

	modelGroup := convertGroupProtoToModel(req.Group)
	if err := dao.UserGroup.Update(ctx, modelGroup); err != nil {
		return nil, status.Error(codes.Internal, "failed to update group")
	}
	return &userV1.UpdateUserGroupResponse{}, nil
}

func (s *S) DeleteUserGroup(ctx context.Context, req *userV1.DeleteUserGroupRequest) (*userV1.DeleteUserGroupResponse, error) {
	// Manual Admin Check
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}
	isAdmin, _ := rbac.CasbinManager.CheckUserIsGlobalAdmin(userID)
	if !isAdmin {
		return nil, status.Error(codes.PermissionDenied, "admin needed")
	}

	if err := dao.UserGroup.Delete(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete group")
	}

	return &userV1.DeleteUserGroupResponse{}, nil
}
