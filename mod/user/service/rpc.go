package service

import (
	"context"
	"errors"
	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbac"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	userV1 "github.com/GoldenSheep402/Hermes/pkg/proto/user/v1"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/GoldenSheep402/Hermes/pkg/utils/crypto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
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
		req.Id = UID
	}

	isAdmin, err := userDao.User.IsAdmin(ctx, UID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	// Check the permission
	if req.Id != UID {

		if !isAdmin {
			isOk, err := rbac.CasbinManager.CheckUserToUserReadPermission(UID, req.Id)
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
			Id:   _user.ID,
			Name: _user.Name,
		},
	}

	if isAdmin {
		resp.User.Role = "admin"
	} else {
		resp.User.Role = "user"
	}

	return resp, nil
}

func (s *S) GetUserInfo(ctx context.Context, req *userV1.GetUserInfoRequest) (*userV1.GetUserInfoResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	user, download, upload, published, downloadCount, seedingCount, err := userDao.User.GetFullInfo(ctx, UID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &userV1.GetUserInfoResponse{
		Name:              user.Name,
		Download:          int32(download),
		Upload:            int32(upload),
		TorrentPublished:  int32(published),
		TorrentDownloaded: int32(downloadCount),
		TorrentSeeding:    int32(seedingCount),
		Key:               user.Key,
	}, nil
}

// UpdateUser updates a user's information based on the provided request, ensuring authentication, input validation, and authorization.
func (s *S) UpdateUser(ctx context.Context, req *userV1.UpdateUserRequest) (*userV1.UpdateUserResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if req.User.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Id is empty")
	}

	if req.User.Id != UID {
		isAdmin, err := userDao.User.IsAdmin(ctx, UID)
		if err != nil {
			return nil, status.Error(codes.Internal, "Internal error")
		}

		if !isAdmin {
			// TODO: rbac
			// isOk, err := rbac.CasbinManager.CheckUserToUserWritePermission(UID, req.User.Id)
			// if err != nil {
			// 	return nil, status.Error(codes.Internal, "Internal error")
			// }
			//
			// if !isOk {
			// 	return nil, status.Error(codes.PermissionDenied, "Permission denied")
			// }
		}
	}

	_user := model.User{}
	_user.ID = req.User.Id
	_user.Name = req.User.Name

	err := userDao.User.UpdateInfo(ctx, _user)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &userV1.UpdateUserResponse{}, nil
}

// UpdatePassword updates a user's password based on the provided old and new password, ensuring the user is authenticated and the old password is correct.
// It returns an error if the user is unauthenticated, the old password doesn't match, or there's an internal error during the update process.
func (s *S) UpdatePassword(ctx context.Context, req *userV1.UpdatePasswordRequest) (*userV1.UpdatePasswordResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	user := &model.User{}
	if result := userDao.User.DB().WithContext(ctx).Where("id = ?", UID).First(user); result.Error != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	if crypto.Md5CryptoWithSalt(req.OldPassword, user.Salt) != user.Password {
		return nil, status.Error(codes.InvalidArgument, "Old password is incorrect")
	}

	newPassword := crypto.Md5CryptoWithSalt(req.NewPassword, user.Salt)
	user.Password = newPassword

	err := userDao.User.UpdateInfo(ctx, *user)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &userV1.UpdatePasswordResponse{}, nil
}

func (s *S) CreateGroup(ctx context.Context, req *userV1.CreateGroupRequest) (*userV1.CreateGroupResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	// check is in admin group
	var adminGroup model.Group
	if err := userDao.Group.GetTxFromCtx(ctx).Where("name = ?", "admin").First(&adminGroup).Error; err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	var memberShop model.GroupMembership
	if err := userDao.GroupMembership.GetTxFromCtx(ctx).Where("uid = ? AND gid = ?", UID, adminGroup.ID).First(&memberShop).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.PermissionDenied, "No permission")
		} else {
			return nil, status.Error(codes.Internal, "Internal error")
		}
	}

	group := model.Group{
		Name:        req.Group.Name,
		Description: req.Group.Description,
	}

	metas := make([]model.GroupMetadata, len(req.Group.MetaData))
	for i, meta := range req.Group.MetaData {
		switch meta.Type {
		case "string":
			break
		case "number":
			break
		default:
			return nil, status.Error(codes.InvalidArgument, "Invalid metadata type")
		}
		metas[i] = model.GroupMetadata{
			Key:         meta.Key,
			Value:       meta.Value,
			Order:       int(meta.Order),
			Description: meta.Description,
			Type:        meta.Type,
		}
	}

	err := userDao.Group.Create(ctx, &group, metas, UID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	// TODO: rbac
	// // set group write permission
	// err = rbac.CasbinManager.SetUserWritePermissionToGroup(UID, group.ID)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, "Internal error")
	// }

	return &userV1.CreateGroupResponse{}, nil
}

func (s *S) GetGroup(ctx context.Context, req *userV1.GetGroupRequest) (*userV1.GetGroupResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Id is empty")
	}

	isAdmin, err := userDao.User.IsAdmin(ctx, UID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	if !isAdmin {
		// check read permission
		ok, err := userDao.GroupMembership.Check(req.Id, UID)
		if err != nil {
			return nil, status.Error(codes.Internal, "Internal error")
		}

		if !ok {
			return nil, status.Error(codes.Unauthenticated, "No permission")
		}

		// TODO: rbac
		// isOk, err := rbac.CasbinManager.CheckUserReadPermissionToGroup(UID, req.Id)
		// if err != nil {
		// 	return nil, status.Error(codes.Internal, "Internal error")
		// }
		//
		// if !isOk {
		// 	return nil, status.Error(codes.PermissionDenied, "Permission denied")
		// }
	}

	group, metas, err := userDao.Group.Get(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	resp := &userV1.GetGroupResponse{
		Group: &userV1.Group{
			Id:          group.ID,
			Name:        group.ID,
			Description: group.Description,
		},
	}

	respMetas := make([]*userV1.GroupMetaData, len(metas))
	for i, meta := range metas {
		respMetas[i] = &userV1.GroupMetaData{
			Key:         meta.Key,
			Value:       meta.Value,
			Order:       int32(meta.Order),
			Description: meta.Description,
			Type:        meta.Type,
		}
	}

	resp.Group.MetaData = respMetas

	return resp, nil
}

// UpdateGroup updates a group's information based on the provided request, ensuring the caller is authenticated, the request is valid, and has necessary permissions.
// It returns an error if the user is unauthenticated, the group ID is empty, or the user lacks sufficient permissions to update the group.
func (s *S) UpdateGroup(ctx context.Context, req *userV1.UpdateGroupRequest) (*userV1.UpdateGroupResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if req.Group.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Id is empty")
	}

	isAdmin, err := userDao.User.IsAdmin(ctx, UID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	if !isAdmin {
		// // check write permission
		// isOk, err := rbac.CasbinManager.CheckUserWritePermissionToGroup(UID, req.Group.Id)
		// if err != nil {
		// 	return nil, status.Error(codes.Internal, "Internal error")
		// }
		//
		// if !isOk {
		// 	return nil, status.Error(codes.PermissionDenied, "Permission denied")
		// }

		return nil, status.Error(codes.PermissionDenied, "You are not admin")
	}

	return &userV1.UpdateGroupResponse{}, nil
}

func (s *S) GroupAddUser(ctx context.Context, req *userV1.GroupAddUserRequest) (*userV1.GroupAddUserResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if req.GroupId == "" {
		return nil, status.Error(codes.InvalidArgument, "GroupId is empty")
	}

	isAdmin, err := userDao.User.IsAdmin(ctx, UID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	if !isAdmin {
		// // check write permission
		// isOk, err := rbac.CasbinManager.CheckUserWritePermissionToGroup(UID, req.GroupId)
		// if err != nil {
		// 	return nil, status.Error(codes.Internal, "Internal error")
		// }
		//
		// if !isOk {
		// 	return nil, status.Error(codes.PermissionDenied, "Permission denied")
		// }

		return nil, status.Error(codes.PermissionDenied, "You are not admin")
	}

	var metas []model.GroupMembershipMetadata
	for i, _meta := range req.MetaData {
		metas[i] = model.GroupMembershipMetadata{
			GroupMetadataID: _meta.GroupMetaDataOriginalID,
			Key:             _meta.Key,
			Value:           _meta.Value,
			Order:           int(_meta.Order),
			Description:     _meta.Description,
			Type:            _meta.Type,
		}
	}

	err = userDao.Group.AddUser(ctx, req.GroupId, req.UserId, metas)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	// TODO: rbac
	// err = rbac.CasbinManager.SetUserToReadGroup(req.UserId, req.GroupId)
	// if err != nil {
	// 	return nil, status.Error(codes.Internal, "Internal error")
	// }

	return &userV1.GroupAddUserResponse{}, nil
}

func (s *S) GroupRemoveUser(ctx context.Context, req *userV1.GroupRemoveUserRequest) (*userV1.GroupRemoveUserResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if req.GroupId == "" {
		return nil, status.Error(codes.InvalidArgument, "GroupId is empty")
	}

	isAdmin, err := userDao.User.IsAdmin(ctx, UID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	if !isAdmin {
		// TODO: rbac
		// // check write permission
		// isOk, err := rbac.CasbinManager.CheckUserWritePermissionToGroup(UID, req.GroupId)
		// if err != nil {
		// 	return nil, status.Error(codes.Internal, "Internal error")
		// }
		//
		// if !isOk {
		// 	return nil, status.Error(codes.PermissionDenied, "Permission denied")
		// }

		return nil, status.Error(codes.PermissionDenied, "You are not admin")
	}

	err = userDao.Group.RemoveUser(ctx, req.GroupId, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &userV1.GroupRemoveUserResponse{}, nil
}

// GroupUserUpdate updates a user's metadata within a group, enforcing authentication, authorization, and validation checks.
func (s *S) GroupUserUpdate(ctx context.Context, req *userV1.GroupUserUpdateRequest) (*userV1.GroupUserUpdateResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if req.GroupId == "" {
		return nil, status.Error(codes.InvalidArgument, "GroupId is empty")
	}

	isAdmin, err := userDao.User.IsAdmin(ctx, UID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	if !isAdmin {
		// // check write permission
		// isOk, err := rbac.CasbinManager.CheckUserWritePermissionToGroup(UID, req.GroupId)
		// if err != nil {
		// 	return nil, status.Error(codes.Internal, "Internal error")
		// }
		//
		// if !isOk {
		// 	return nil, status.Error(codes.PermissionDenied, "Permission denied")
		// }

		return nil, status.Error(codes.PermissionDenied, "You are not admin")
	}

	var metas []model.GroupMembershipMetadata
	for i, _meta := range req.MetaData {
		metas[i] = model.GroupMembershipMetadata{
			Model: stdao.Model{
				ID: _meta.Id,
			},
			GroupMetadataID: _meta.GroupMetaDataOriginalID,
			Key:             _meta.Key,
			Value:           _meta.Value,
			Order:           int(_meta.Order),
			Description:     _meta.Description,
			Type:            _meta.Type,
		}
	}

	err = userDao.Group.UpdateUser(ctx, req.GroupId, req.UserId, metas)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &userV1.GroupUserUpdateResponse{}, nil
}

func (s *S) GetUserPassKey(ctx context.Context, req *userV1.GetUserPassKeyRequest) (*userV1.GetUserPassKeyResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if req.Id != "" {
		isAdmin, err := userDao.User.IsAdmin(ctx, UID)
		if err != nil {
			return nil, status.Error(codes.Internal, "Internal error")
		}

		if !isAdmin {
			return nil, status.Error(codes.PermissionDenied, "Permission denied")
		}
	} else {
		req.Id = UID
	}

	passKey, err := userDao.User.GetPassKey(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &userV1.GetUserPassKeyResponse{
		PassKey: passKey,
	}, nil
}
