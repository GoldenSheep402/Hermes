package service

import (
	resourceV1 "github.com/GoldenSheep402/Hermes/pkg/proto/resource/v1"

	"context"
	"strings"

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

func normalizeProtoMetadata(items []*resourceV1.ResourceMeta) []model.ResourceMeta {
	result := make([]model.ResourceMeta, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		key := strings.TrimSpace(item.Key)
		value := strings.TrimSpace(item.Value)
		if key == "" || value == "" {
			continue
		}
		result = append(result, model.ResourceMeta{
			Model: stdao.Model{ID: item.Id},
			Key:   key,
			Value: value,
		})
	}
	return result
}

func toProtoMetadata(items []model.ResourceMeta) []*resourceV1.ResourceMeta {
	result := make([]*resourceV1.ResourceMeta, 0, len(items))
	for _, item := range items {
		result = append(result, &resourceV1.ResourceMeta{
			Id:    item.ID,
			Key:   item.Key,
			Value: item.Value,
		})
	}
	return result
}

func toProtoTags(items []model.Tag) []*resourceV1.Tag {
	result := make([]*resourceV1.Tag, 0, len(items))
	for _, item := range items {
		result = append(result, &resourceV1.Tag{
			Id:   item.ID,
			Name: item.Name,
		})
	}
	return result
}

func toProtoScreenshots(items []model.ResourceScreenshot) []*resourceV1.ResourceScreenshot {
	result := make([]*resourceV1.ResourceScreenshot, 0, len(items))
	for _, item := range items {
		result = append(result, &resourceV1.ResourceScreenshot{
			Id:        item.ID,
			Url:       item.URL,
			SortOrder: int32(item.SortOrder),
		})
	}
	return result
}

func saveResourceDetails(ctx context.Context, resourceID string, metadata []*resourceV1.ResourceMeta, tagNames []string, screenshotURLs []string) error {
	if err := dao.ResourceMeta.ReplaceByResourceID(ctx, resourceID, normalizeProtoMetadata(metadata)); err != nil {
		return err
	}

	tags, err := dao.Tag.EnsureByNames(ctx, tagNames)
	if err != nil {
		return err
	}
	tagIDs := make([]string, 0, len(tags))
	for _, item := range tags {
		if item.ID != "" {
			tagIDs = append(tagIDs, item.ID)
		}
	}
	if err := dao.ResourceTag.ReplaceByResourceID(ctx, resourceID, tagIDs); err != nil {
		return err
	}

	if err := dao.ResourceScreenshot.ReplaceByResourceID(ctx, resourceID, screenshotURLs); err != nil {
		return err
	}

	return nil
}

func (s *S) CreateResource(ctx context.Context, req *resourceV1.CreateResourceRequest) (*resourceV1.CreateResourceResponse, error) {
	uploaderID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || uploaderID == "" {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	if req.Title == "" || req.CategoryId == "" || len(req.TorrentData) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Title, CategoryID, and TorrentData are required")
	}

	uploader, err := userDao.User.GetByID(ctx, uploaderID)
	if err != nil || uploader == nil || uploader.Passkey == "" {
		return nil, status.Error(codes.Unauthenticated, "invalid user passkey")
	}

	announceURLs, err := torrentCommon.BuildAnnounceURLsForPasskey(ctx, uploader.Passkey)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "tracker endpoint not available")
	}

	// Rewrite for upload persistence: bind announce URLs and enforce private tracker mode.
	normalizedData, err := torrent.RewriteUploadTorrentWithTrackers(req.TorrentData, announceURLs)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid torrent data")
	}
	parsed, err := torrent.Parse(normalizedData)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid torrent data")
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

	tx := dao.Resource.Begin()
	txCtx := dao.Resource.SetTxToCtx(ctx, tx)

	if err := dao.Resource.Create(txCtx, res); err != nil {
		_ = tx.Rollback().Error
		s.Log.Errorw("failed to create resource", "err", err)
		return nil, status.Error(codes.Internal, "Failed to create resource")
	}
	if err := saveResourceDetails(txCtx, res.ID, req.Metadata, req.TagNames, req.ScreenshotUrls); err != nil {
		_ = tx.Rollback().Error
		s.Log.Errorw("failed to persist resource details", "resource_id", res.ID, "err", err)
		return nil, status.Error(codes.Internal, "Failed to save resource details")
	}
	if err := tx.Commit().Error; err != nil {
		_ = tx.Rollback().Error
		s.Log.Errorw("failed to commit resource transaction", "resource_id", res.ID, "err", err)
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

	metadata, err := dao.ResourceMeta.ListByResourceID(ctx, res.ID)
	if err != nil {
		s.Log.Errorw("failed to load resource metadata", "resource_id", res.ID, "err", err)
		return nil, status.Error(codes.Internal, "Failed to load resource details")
	}
	tagMap, err := dao.ResourceTag.ListTagMapByResourceIDs(ctx, []string{res.ID})
	if err != nil {
		s.Log.Errorw("failed to load resource tags", "resource_id", res.ID, "err", err)
		return nil, status.Error(codes.Internal, "Failed to load resource details")
	}
	screenshots, err := dao.ResourceScreenshot.ListByResourceID(ctx, res.ID)
	if err != nil {
		s.Log.Errorw("failed to load resource screenshots", "resource_id", res.ID, "err", err)
		return nil, status.Error(codes.Internal, "Failed to load resource details")
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
			Metadata:     toProtoMetadata(metadata),
			Tags:         toProtoTags(tagMap[res.ID]),
			Screenshots:  toProtoScreenshots(screenshots),
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

	resourceIDs := make([]string, 0, len(list))
	for _, item := range list {
		if item.ID != "" {
			resourceIDs = append(resourceIDs, item.ID)
		}
	}

	metadataMap, err := dao.ResourceMeta.ListByResourceIDs(ctx, resourceIDs)
	if err != nil {
		s.Log.Errorw("failed to load resource metadata batch", "err", err)
		return nil, status.Error(codes.Internal, "Failed to load resource details")
	}
	tagMap, err := dao.ResourceTag.ListTagMapByResourceIDs(ctx, resourceIDs)
	if err != nil {
		s.Log.Errorw("failed to load resource tags batch", "err", err)
		return nil, status.Error(codes.Internal, "Failed to load resource details")
	}
	screenshotMap, err := dao.ResourceScreenshot.ListByResourceIDs(ctx, resourceIDs)
	if err != nil {
		s.Log.Errorw("failed to load resource screenshots batch", "err", err)
		return nil, status.Error(codes.Internal, "Failed to load resource details")
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
			Metadata:     toProtoMetadata(metadataMap[res.ID]),
			Tags:         toProtoTags(tagMap[res.ID]),
			Screenshots:  toProtoScreenshots(screenshotMap[res.ID]),
		})
	}

	return &resourceV1.ListResourcesResponse{
		Resources: pList,
		Total:     count,
	}, nil
}

func (s *S) UpdateResource(ctx context.Context, req *resourceV1.UpdateResourceRequest) (*resourceV1.UpdateResourceResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Resource ID required")
	}

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

	tx := dao.Resource.Begin()
	txCtx := dao.Resource.SetTxToCtx(ctx, tx)

	if len(updates) > 0 {
		if err := dao.Resource.UpdateFields(txCtx, req.Id, updates); err != nil {
			_ = tx.Rollback().Error
			s.Log.Errorw("failed to update resource", "resource_id", req.Id, "err", err)
			return nil, status.Error(codes.Internal, "Failed to update resource")
		}
	}

	if req.Metadata != nil {
		if err := dao.ResourceMeta.ReplaceByResourceID(txCtx, req.Id, normalizeProtoMetadata(req.Metadata)); err != nil {
			_ = tx.Rollback().Error
			s.Log.Errorw("failed to replace resource metadata", "resource_id", req.Id, "err", err)
			return nil, status.Error(codes.Internal, "Failed to update resource metadata")
		}
	}

	if req.TagNames != nil {
		tags, err := dao.Tag.EnsureByNames(txCtx, req.TagNames)
		if err != nil {
			_ = tx.Rollback().Error
			s.Log.Errorw("failed to ensure tags", "resource_id", req.Id, "err", err)
			return nil, status.Error(codes.Internal, "Failed to update resource tags")
		}
		tagIDs := make([]string, 0, len(tags))
		for _, item := range tags {
			if item.ID != "" {
				tagIDs = append(tagIDs, item.ID)
			}
		}
		if err := dao.ResourceTag.ReplaceByResourceID(txCtx, req.Id, tagIDs); err != nil {
			_ = tx.Rollback().Error
			s.Log.Errorw("failed to replace resource tags", "resource_id", req.Id, "err", err)
			return nil, status.Error(codes.Internal, "Failed to update resource tags")
		}
	}

	if req.ScreenshotUrls != nil {
		if err := dao.ResourceScreenshot.ReplaceByResourceID(txCtx, req.Id, req.ScreenshotUrls); err != nil {
			_ = tx.Rollback().Error
			s.Log.Errorw("failed to replace resource screenshots", "resource_id", req.Id, "err", err)
			return nil, status.Error(codes.Internal, "Failed to update resource screenshots")
		}
	}

	if err := tx.Commit().Error; err != nil {
		_ = tx.Rollback().Error
		s.Log.Errorw("failed to commit resource update", "resource_id", req.Id, "err", err)
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
