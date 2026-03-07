package service

import (
	"context"
	"testing"

	commentV1 "github.com/GoldenSheep402/Hermes/pkg/proto/comment/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCommentService_Actions_NoContext(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	r1, err1 := s.CreateComment(context.Background(), &commentV1.CreateCommentRequest{})
	assert.Error(t, err1)
	assert.Nil(t, r1)
	assert.Contains(t, err1.Error(), "unauthenticated")

	r2, err2 := s.ListComments(context.Background(), &commentV1.ListCommentsRequest{})
	assert.Error(t, err2)
	assert.Nil(t, r2)
	assert.Contains(t, err2.Error(), "unauthenticated")

	r3, err3 := s.UpdateComment(context.Background(), &commentV1.UpdateCommentRequest{})
	assert.Error(t, err3)
	assert.Nil(t, r3)
	assert.Contains(t, err3.Error(), "unauthenticated")

	r4, err4 := s.DeleteComment(context.Background(), &commentV1.DeleteCommentRequest{})
	assert.Error(t, err4)
	assert.Nil(t, r4)
	assert.Contains(t, err4.Error(), "unauthenticated")
}
