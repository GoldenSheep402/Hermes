package service

import (
	commentV1 "github.com/GoldenSheep402/Hermes/pkg/proto/comment/v1"

	"context"

	"go.uber.org/zap"
)

var _ commentV1.CommentServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	commentV1.UnimplementedCommentServiceServer
}

func (s S) CreateComment(ctx context.Context, request *commentV1.CreateCommentRequest) (*commentV1.CreateCommentResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) ListComments(ctx context.Context, request *commentV1.ListCommentsRequest) (*commentV1.ListCommentsResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) UpdateComment(ctx context.Context, request *commentV1.UpdateCommentRequest) (*commentV1.UpdateCommentResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) DeleteComment(ctx context.Context, request *commentV1.DeleteCommentRequest) (*commentV1.DeleteCommentResponse, error) {
	// TODO implement me
	panic("implement me")
}
