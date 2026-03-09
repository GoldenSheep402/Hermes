package service

import (
	commentV1 "github.com/GoldenSheep402/Hermes/pkg/proto/comment/v1"

	"context"

	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbac"
	"github.com/GoldenSheep402/Hermes/mod/comment/dao"
	"github.com/GoldenSheep402/Hermes/mod/comment/model"
	resourceDao "github.com/GoldenSheep402/Hermes/mod/resource/dao"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.uber.org/zap"
)

var _ commentV1.CommentServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	commentV1.UnimplementedCommentServiceServer
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

func (s S) CreateComment(ctx context.Context, req *commentV1.CreateCommentRequest) (*commentV1.CreateCommentResponse, error) {
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}
	if req.ResourceId == "" || req.Body == "" {
		return nil, status.Error(codes.InvalidArgument, "Resource ID and Body are required")
	}

	var parentID *string
	if req.ParentId != "" {
		pid := req.ParentId
		parentID = &pid
	}

	c := &model.Comment{
		Model:      stdao.Model{ID: ulid.Make().String()},
		ResourceID: req.ResourceId,
		AuthorID:   userID,
		Body:       req.Body,
		ParentID:   parentID,
	}

	tx := dao.Comment.Begin()
	txCtx := dao.Comment.SetTxToCtx(ctx, tx)

	if err := dao.Comment.Create(txCtx, c); err != nil {
		_ = tx.Rollback().Error
		s.Log.Errorw("failed to create comment", "err", err)
		return nil, status.Error(codes.Internal, "Failed to create comment")
	}
	if err := resourceDao.Resource.AdjustCommentCount(txCtx, req.ResourceId, 1); err != nil {
		_ = tx.Rollback().Error
		s.Log.Errorw("failed to increment resource comment count", "resource_id", req.ResourceId, "err", err)
		return nil, status.Error(codes.Internal, "Failed to create comment")
	}
	if err := tx.Commit().Error; err != nil {
		_ = tx.Rollback().Error
		s.Log.Errorw("failed to commit create comment", "resource_id", req.ResourceId, "err", err)
		return nil, status.Error(codes.Internal, "Failed to create comment")
	}

	return &commentV1.CreateCommentResponse{Id: c.ID}, nil
}

func (s S) ListComments(ctx context.Context, req *commentV1.ListCommentsRequest) (*commentV1.ListCommentsResponse, error) {
	if _, err := requireAuth(ctx); err != nil {
		return nil, err
	}
	if req.ResourceId == "" {
		return nil, status.Error(codes.InvalidArgument, "Resource ID required")
	}

	list, count, err := dao.Comment.ListByResource(ctx, req.ResourceId, req.Page, req.PageSize)
	if err != nil {
		s.Log.Errorw("failed to list comments", "err", err)
		return nil, status.Error(codes.Internal, "Failed to list comments")
	}

	var pList []*commentV1.CommentInfo
	// TODO: Fetch user avatars and names via User module.
	for _, c := range list {
		pList = append(pList, &commentV1.CommentInfo{
			Id:           c.ID,
			ResourceId:   c.ResourceID,
			AuthorId:     c.AuthorID,
			AuthorName:   "User " + c.AuthorID, // Placeholder
			AuthorAvatar: "",
			Body:         c.Body,
			CreatedAt:    c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			// Replies logic can be iterated here if preloaded.
		})
	}

	return &commentV1.ListCommentsResponse{
		Comments: pList,
		Total:    count,
	}, nil
}

func (s S) UpdateComment(ctx context.Context, req *commentV1.UpdateCommentRequest) (*commentV1.UpdateCommentResponse, error) {
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}
	if req.Id == "" || req.Body == "" {
		return nil, status.Error(codes.InvalidArgument, "Comment ID and Body required")
	}

	// Authorization: only author or admin can edit
	c, err := dao.Comment.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Comment not found")
	}
	if c.AuthorID != userID {
		if errAdmin := requireAdmin(ctx); errAdmin != nil {
			return nil, status.Error(codes.PermissionDenied, "Not authorized to edit this comment")
		}
	}

	if err := dao.Comment.Update(ctx, &model.Comment{Model: stdao.Model{ID: req.Id}, Body: req.Body}); err != nil {
		s.Log.Errorw("failed to update comment", "err", err)
		return nil, status.Error(codes.Internal, "Failed to update comment")
	}

	return &commentV1.UpdateCommentResponse{}, nil
}

func (s S) DeleteComment(ctx context.Context, req *commentV1.DeleteCommentRequest) (*commentV1.DeleteCommentResponse, error) {
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Comment ID required")
	}

	tx := dao.Comment.Begin()
	txCtx := dao.Comment.SetTxToCtx(ctx, tx)

	// Authorization: only author or admin can delete
	c, err := dao.Comment.GetByID(txCtx, req.Id)
	if err != nil {
		_ = tx.Rollback().Error
		return nil, status.Error(codes.NotFound, "Comment not found")
	}
	if c.AuthorID != userID {
		if errAdmin := requireAdmin(ctx); errAdmin != nil {
			_ = tx.Rollback().Error
			return nil, status.Error(codes.PermissionDenied, "Not authorized to delete this comment")
		}
	}

	if err := dao.Comment.Delete(txCtx, &model.Comment{Model: stdao.Model{ID: req.Id}}).Error; err != nil {
		_ = tx.Rollback().Error
		s.Log.Errorw("failed to delete comment", "err", err)
		return nil, status.Error(codes.Internal, "Failed to delete comment")
	}
	if err := resourceDao.Resource.AdjustCommentCount(txCtx, c.ResourceID, -1); err != nil {
		_ = tx.Rollback().Error
		s.Log.Errorw("failed to decrement resource comment count", "resource_id", c.ResourceID, "err", err)
		return nil, status.Error(codes.Internal, "Failed to delete comment")
	}
	if err := tx.Commit().Error; err != nil {
		_ = tx.Rollback().Error
		s.Log.Errorw("failed to commit delete comment", "comment_id", req.Id, "err", err)
		return nil, status.Error(codes.Internal, "Failed to delete comment")
	}

	return &commentV1.DeleteCommentResponse{}, nil
}
