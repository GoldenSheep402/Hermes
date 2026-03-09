package service

import (
	"context"
	"testing"

	categoryV1 "github.com/GoldenSheep402/Hermes/pkg/proto/category/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Similar to Torrent module, we test error bounds and structure here since actual DB calls need initialized connections.
func TestCategoryService_GetCategory(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	resp, err := s.GetCategory(context.Background(), &categoryV1.GetCategoryRequest{Id: ""})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthenticated")
}

func TestCategoryService_DeleteCategory_NoContext(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	// Should fail immediately because the caller has no UID in context and thus no admin rights
	resp, err := s.DeleteCategory(context.Background(), &categoryV1.DeleteCategoryRequest{Id: "123"})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthenticated")
}

func TestCategoryService_MetaTemplates_NoContext(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	r1, err1 := s.CreateMetaTemplate(context.Background(), &categoryV1.CreateMetaTemplateRequest{})
	assert.Error(t, err1)
	assert.Nil(t, r1)
	assert.Contains(t, err1.Error(), "unauthenticated")

	r2, err2 := s.UpdateMetaTemplate(context.Background(), &categoryV1.UpdateMetaTemplateRequest{})
	assert.Error(t, err2)
	assert.Nil(t, r2)
	assert.Contains(t, err2.Error(), "unauthenticated")

	r3, err3 := s.DeleteMetaTemplate(context.Background(), &categoryV1.DeleteMetaTemplateRequest{})
	assert.Error(t, err3)
	assert.Nil(t, r3)
	assert.Contains(t, err3.Error(), "unauthenticated")

	r4, err4 := s.ListMetaTemplatePresets(context.Background(), &categoryV1.ListMetaTemplatePresetsRequest{})
	assert.Error(t, err4)
	assert.Nil(t, r4)
	assert.Contains(t, err4.Error(), "unauthenticated")
}
