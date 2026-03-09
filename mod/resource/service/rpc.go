package service

import (
	resourceV1 "github.com/GoldenSheep402/Hermes/pkg/proto/resource/v1"

	"context"

	"time"

	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbac"
	"github.com/GoldenSheep402/Hermes/mod/resource/dao"
	"github.com/GoldenSheep402/Hermes/mod/resource/model"
	torrentCommon "github.com/GoldenSheep402/Hermes/mod/torrent/common"
	torrentDao "github.com/GoldenSheep402/Hermes/mod/torrent/dao"
	torrentModel "github.com/GoldenSheep402/Hermes/mod/torrent/model"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	userModel "github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/GoldenSheep402/Hermes/pkg/torrent"
)

var _ resourceV1.ResourceServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	resourceV1.UnimplementedResourceServiceServer
}

// requireAdmin checks if caller is highly privileged.
func requireAdmin(ctx context.Context) error {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return status.Error(codes.Unauthenticated, "unauthenticated")
	}
	isAdmin, err := rbac.CasbinManager.CheckUserIsGlobalAdmin(userID)
	if err != nil || !isAdmin {
		return status.Error(codes.PermissionDenied, "admin privileges required")
	}
	return nil
}

// requireAuth checks if caller is authenticated.
func requireAuth(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return "", status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return userID, nil
}

func (s *S) CreateResource(ctx context.Context, req *resourceV1.CreateResourceRequest) (*resourceV1.CreateResourceResponse, error) {
	uploaderID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || uploaderID == "" {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	if req.Title == "" || req.CategoryId == "" || len(req.TorrentData) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Title, CategoryID, and TorrentData are required")
	}

	// 1. Parse original torrent first.
	// InfoHash is derived from the info dict and should not depend on announce URLs.
	parsed, err := torrent.Parse(req.TorrentData)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid torrent data")
	}

	uploader, err := userDao.User.GetByID(ctx, uploaderID)
	if err != nil || uploader == nil || uploader.Passkey == "" {
		return nil, status.Error(codes.Unauthenticated, "invalid user passkey")
	}

	announceURLs, err := torrentCommon.BuildAnnounceURLsForPasskey(ctx, uploader.Passkey)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "tracker endpoint not available")
	}

	// Rewrite tracker URLs only for persisted raw bytes.
	normalizedData, err := torrent.RewriteDownloadTorrentWithTrackers(req.TorrentData, announceURLs)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to normalize torrent data")
	}

	// 2. Process Torrent
	tID := ulid.Make().String()
	torrentBase := &torrentModel.Torrent{
		Model:        stdao.Model{ID: tID},
		InfoHash:     parsed.InfoHash,
		UploaderID:   uploaderID,
		Name:         parsed.Name,
		Size:         parsed.Size,
		PieceLength:  parsed.PieceLength,
		PieceCount:   parsed.PieceCount,
		IsSingleFile: len(parsed.Files) <= 1,
		FileCount:    len(parsed.Files),
	}
	var files []torrentModel.TorrentFile
	for _, f := range parsed.Files {
		files = append(files, torrentModel.TorrentFile{
			Model:     stdao.Model{ID: ulid.Make().String()},
			TorrentID: tID,
			Path:      f.Path,
			Size:      f.Size,
		})
	}
	if _, err := torrentDao.Torrent.Create(ctx, torrentBase, files, normalizedData, parsed.PieceHashes); err != nil {
		s.Log.Errorw("failed to save torrent", "err", err)
		if status.Code(err) == codes.AlreadyExists {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "Failed to create torrent record")
	}

	// 2. Process Resource
	rID := ulid.Make().String()
	res := &model.Resource{
		Model:       stdao.Model{ID: rID},
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Description: req.Description,
		CategoryID:  req.CategoryId,
		TorrentID:   tID,
		UploaderID:  uploaderID,
		Status:      model.ResourceStatusPending, // Depending on config, this could be approved instantly
		IsSticky:    false,
		IsFree:      false,
	}

	if err := dao.Resource.Create(ctx, res); err != nil {
		s.Log.Errorw("failed to create resource", "err", err)
		return nil, status.Error(codes.Internal, "Failed to create resource")
	}

	return &resourceV1.CreateResourceResponse{Id: res.ID}, nil
}

func (s *S) GetResource(ctx context.Context, req *resourceV1.GetResourceRequest) (*resourceV1.GetResourceResponse, error) {
	if _, err := requireAuth(ctx); err != nil {
		return nil, err
	}

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Resource ID required")
	}

	res, err := dao.Resource.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Resource not found")
	}

	// Basic view increment (should be deferred or asynchronous in prod)
	_ = dao.Resource.IncrementViewCount(ctx, res.ID)

	var freeUntil, doubleUntil string
	if res.FreeUntil != nil {
		freeUntil = res.FreeUntil.Format(time.RFC3339)
	}
	if res.DoubleUntil != nil {
		doubleUntil = res.DoubleUntil.Format(time.RFC3339)
	}
	uploaderName := ""
	if res.UploaderID != "" {
		if uploader, err := userDao.User.GetByID(ctx, res.UploaderID); err == nil && uploader != nil {
			uploaderName = uploader.Username
		}
	}

	return &resourceV1.GetResourceResponse{
		Resource: &resourceV1.Resource{
			Id:           res.ID,
			Title:        res.Title,
			Subtitle:     res.Subtitle,
			Description:  res.Description,
			CategoryId:   res.CategoryID,
			TorrentId:    res.TorrentID,
			UploaderId:   res.UploaderID,
			Status:       int32(res.Status),
			IsSticky:     res.IsSticky,
			IsFree:       res.IsFree,
			FreeUntil:    freeUntil,
			DoubleUpload: res.DoubleUpload,
			DoubleUntil:  doubleUntil,
			ViewCount:    int32(res.ViewCount),
			CommentCount: int32(res.CommentCount),
			ThankCount:   int32(res.ThankCount),
			CreatedAt:    res.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    res.UpdatedAt.Format(time.RFC3339),
			UploaderName: uploaderName,
		},
	}, nil
}

func (s *S) ListResources(ctx context.Context, req *resourceV1.ListResourcesRequest) (*resourceV1.ListResourcesResponse, error) {
	if _, err := requireAuth(ctx); err != nil {
		return nil, err
	}

	list, count, err := dao.Resource.ListByFilters(ctx, req.CategoryId, req.Keyword, req.Status, req.Page, req.PageSize)
	if err != nil {
		s.Log.Errorw("failed to list resources", "err", err)
		return nil, status.Error(codes.Internal, "Failed to list resources")
	}

	uploaderIDs := make([]string, 0, len(list))
	uploaderIDSet := make(map[string]struct{}, len(list))
	for _, item := range list {
		if item.UploaderID == "" {
			continue
		}
		if _, exists := uploaderIDSet[item.UploaderID]; exists {
			continue
		}
		uploaderIDSet[item.UploaderID] = struct{}{}
		uploaderIDs = append(uploaderIDs, item.UploaderID)
	}

	uploaderNameMap := make(map[string]string, len(uploaderIDs))
	if len(uploaderIDs) > 0 {
		var users []userModel.User
		err = userDao.User.GetTxFromCtx(ctx).
			WithContext(ctx).
			Select("id", "username").
			Where("id IN ?", uploaderIDs).
			Find(&users).Error
		if err != nil {
			s.Log.Warnw("failed to load uploader usernames", "err", err)
		} else {
			for _, user := range users {
				uploaderNameMap[user.ID] = user.Username
			}
		}
	}

	var pList []*resourceV1.Resource
	for _, res := range list {
		var freeUntil, doubleUntil string
		if res.FreeUntil != nil {
			freeUntil = res.FreeUntil.Format(time.RFC3339)
		}
		if res.DoubleUntil != nil {
			doubleUntil = res.DoubleUntil.Format(time.RFC3339)
		}
		pList = append(pList, &resourceV1.Resource{
			Id:           res.ID,
			Title:        res.Title,
			Subtitle:     res.Subtitle,
			CategoryId:   res.CategoryID,
			TorrentId:    res.TorrentID,
			UploaderId:   res.UploaderID,
			Status:       int32(res.Status),
			IsSticky:     res.IsSticky,
			IsFree:       res.IsFree,
			FreeUntil:    freeUntil,
			DoubleUpload: res.DoubleUpload,
			DoubleUntil:  doubleUntil,
			ViewCount:    int32(res.ViewCount),
			CommentCount: int32(res.CommentCount),
			ThankCount:   int32(res.ThankCount),
			CreatedAt:    res.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    res.UpdatedAt.Format(time.RFC3339),
			UploaderName: uploaderNameMap[res.UploaderID],
		})
	}

	return &resourceV1.ListResourcesResponse{
		Resources: pList,
		Total:     count,
	}, nil
}

func (s *S) UpdateResource(ctx context.Context, req *resourceV1.UpdateResourceRequest) (*resourceV1.UpdateResourceResponse, error) {
	// Access control: creator or admin
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}
	res, err := dao.Resource.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Resource not found")
	}
	if res.UploaderID != userID {
		if errAdmin := requireAdmin(ctx); errAdmin != nil {
			return nil, status.Error(codes.PermissionDenied, "Not authorized to update this resource")
		}
	}

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Resource ID required")
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Subtitle != "" {
		updates["subtitle"] = req.Subtitle
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.CategoryId != "" {
		updates["category_id"] = req.CategoryId
	}

	if err := dao.Resource.UpdateFields(ctx, req.Id, updates); err != nil {
		s.Log.Errorw("failed to update resource", "err", err)
		return nil, status.Error(codes.Internal, "Failed to update resource")
	}

	return &resourceV1.UpdateResourceResponse{}, nil
}

func (s *S) DeleteResource(ctx context.Context, req *resourceV1.DeleteResourceRequest) (*resourceV1.DeleteResourceResponse, error) {
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}
	res, err := dao.Resource.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Resource not found")
	}
	if res.UploaderID != userID {
		if errAdmin := requireAdmin(ctx); errAdmin != nil {
			return nil, status.Error(codes.PermissionDenied, "Not authorized to delete this resource")
		}
	}

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Resource ID required")
	}

	if err := dao.Resource.Delete(ctx, &model.Resource{Model: stdao.Model{ID: req.Id}}).Error; err != nil {
		s.Log.Errorw("failed to delete resource", "err", err)
		return nil, status.Error(codes.Internal, "Failed to delete resource")
	}

	return &resourceV1.DeleteResourceResponse{}, nil
}

func (s *S) ThankResource(ctx context.Context, req *resourceV1.ThankResourceRequest) (*resourceV1.ThankResourceResponse, error) {
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}

	if req.ResourceId == "" {
		return nil, status.Error(codes.InvalidArgument, "Resource ID required")
	}

	rt := &model.ResourceThank{
		Model:      stdao.Model{ID: ulid.Make().String()},
		ResourceID: req.ResourceId,
		UserID:     userID,
	}

	if err := dao.ResourceThank.AddThank(ctx, rt); err != nil {
		s.Log.Errorw("failed to thank resource", "err", err)
		return nil, status.Error(codes.Internal, "Failed to thank resource: "+err.Error())
	}

	return &resourceV1.ThankResourceResponse{}, nil
}

func (s *S) SetPromotion(ctx context.Context, req *resourceV1.SetPromotionRequest) (*resourceV1.SetPromotionResponse, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	if req.ResourceId == "" {
		return nil, status.Error(codes.InvalidArgument, "Resource ID required")
	}

	updates := map[string]interface{}{
		"is_free":       req.IsFree,
		"double_upload": req.DoubleUpload,
	}

	if req.FreeUntil != "" {
		t, err := time.Parse(time.RFC3339, req.FreeUntil)
		if err == nil {
			updates["free_until"] = t
		} else {
			return nil, status.Error(codes.InvalidArgument, "Invalid FreeUntil time format")
		}
	} else {
		updates["free_until"] = nil
	}

	if req.DoubleUntil != "" {
		t, err := time.Parse(time.RFC3339, req.DoubleUntil)
		if err == nil {
			updates["double_until"] = t
		} else {
			return nil, status.Error(codes.InvalidArgument, "Invalid DoubleUntil time format")
		}
	} else {
		updates["double_until"] = nil
	}

	if err := dao.Resource.UpdateFields(ctx, req.ResourceId, updates); err != nil {
		s.Log.Errorw("failed to set promotion", "err", err)
		return nil, status.Error(codes.Internal, "Failed to update resource promotion")
	}

	return &resourceV1.SetPromotionResponse{}, nil
}
