package service

import (
	"context"
	"testing"

	resourceV1 "github.com/GoldenSheep402/Hermes/pkg/proto/resource/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestResourceService_GetResource(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	resp, err := s.GetResource(context.Background(), &resourceV1.GetResourceRequest{Id: ""})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthenticated")
}

func TestResourceService_DeleteResource_NoContext(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	// Should fail immediately because caller has no rights
	resp, err := s.DeleteResource(context.Background(), &resourceV1.DeleteResourceRequest{Id: "r123"})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthenticated")
}

func TestResourceService_Actions_NoContext(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	r1, err1 := s.ThankResource(context.Background(), &resourceV1.ThankResourceRequest{})
	assert.Error(t, err1)
	assert.Nil(t, r1)
	assert.Contains(t, err1.Error(), "unauthenticated")

	r2, err2 := s.SetPromotion(context.Background(), &resourceV1.SetPromotionRequest{})
	assert.Error(t, err2)
	assert.Nil(t, r2)
	assert.Contains(t, err2.Error(), "unauthenticated")
}
