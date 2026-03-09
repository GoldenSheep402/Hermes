package service

import (
	"context"

	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbac"
	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbacValues"
	"github.com/GoldenSheep402/Hermes/mod/category/dao"
	"github.com/GoldenSheep402/Hermes/mod/category/model"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	categoryV1 "github.com/GoldenSheep402/Hermes/pkg/proto/category/v1"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

var _ categoryV1.CategoryServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	categoryV1.UnimplementedCategoryServiceServer
}

// requireAdmin is a helper to check if the caller context has admin privileges.
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

// requireCategoryManage is a future-proof hook:
// allow global admin OR users granted category write permission in Casbin.
func requireCategoryManage(ctx context.Context, categoryID string) error {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return status.Error(codes.Unauthenticated, "unauthenticated")
	}

	isAdmin, err := rbac.CasbinManager.CheckUserIsGlobalAdmin(userID)
	if err == nil && isAdmin {
		return nil
	}

	resourceCandidates := []string{
		rbacValues.CategoryIDPrefix("*"),
	}
	if categoryID != "" {
		resourceCandidates = append(resourceCandidates, rbacValues.CategoryIDPrefix(categoryID))
	}

	for _, resource := range resourceCandidates {
		allowed, checkErr := rbac.CasbinManager.CheckUserPermission(userID, resource, rbacValues.Write)
		if checkErr == nil && allowed {
			return nil
		}
	}

	return status.Error(codes.PermissionDenied, "category manage permission required")
}

// requireAuth checks if caller is authenticated.
func requireAuth(ctx context.Context) error {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return nil
}

func convertModelToProto(c *model.Category) *categoryV1.Category {
	if c == nil {
		return nil
	}
	var parentID string
	if c.ParentID != nil {
		parentID = *c.ParentID
	}
	return &categoryV1.Category{
		Id:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		ParentId:    parentID,
		Icon:        c.Icon,
		SortOrder:   int32(c.SortOrder),
		IsEnabled:   c.IsEnabled,
	}
}

func (s *S) CreateCategory(ctx context.Context, req *categoryV1.CreateCategoryRequest) (*categoryV1.CreateCategoryResponse, error) {
	if err := requireCategoryManage(ctx, ""); err != nil {
		return nil, err
	}

	if req.Name == "" || req.Slug == "" {
		return nil, status.Error(codes.InvalidArgument, "Name and Slug are required")
	}

	var parentID *string
	if req.ParentId != "" {
		pid := req.ParentId
		parentID = &pid
	}

	cat := &model.Category{
		Model:       stdao.Model{ID: ulid.Make().String()},
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    parentID,
		Icon:        req.Icon,
		SortOrder:   int(req.SortOrder),
		IsEnabled:   true,
	}

	if err := dao.Category.Create(ctx, cat); err != nil {
		s.Log.Errorw("failed to create category", "err", err)
		return nil, status.Error(codes.Internal, "Failed to create category")
	}

	return &categoryV1.CreateCategoryResponse{Id: cat.ID}, nil
}

func (s *S) GetCategory(ctx context.Context, req *categoryV1.GetCategoryRequest) (*categoryV1.GetCategoryResponse, error) {
	if err := requireAuth(ctx); err != nil {
		return nil, err
	}

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Category ID required")
	}

	cat, err := dao.Category.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Category not found")
	}

	return &categoryV1.GetCategoryResponse{
		Category: convertModelToProto(cat),
	}, nil
}

func (s *S) ListCategories(ctx context.Context, req *categoryV1.ListCategoriesRequest) (*categoryV1.ListCategoriesResponse, error) {
	if err := requireAuth(ctx); err != nil {
		return nil, err
	}

	cats, err := dao.Category.List(ctx)
	if err != nil {
		s.Log.Errorw("failed to list categories", "err", err)
		return nil, status.Error(codes.Internal, "Failed to list categories")
	}

	var protoCats []*categoryV1.Category
	for _, c := range cats {
		protoCats = append(protoCats, convertModelToProto(c))
	}

	return &categoryV1.ListCategoriesResponse{Categories: protoCats}, nil
}

func (s *S) UpdateCategory(ctx context.Context, req *categoryV1.UpdateCategoryRequest) (*categoryV1.UpdateCategoryResponse, error) {
	categoryID := ""
	if req != nil && req.Category != nil {
		categoryID = req.Category.Id
	}
	if err := requireCategoryManage(ctx, categoryID); err != nil {
		return nil, err
	}
	if req.Category == nil || req.Category.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Category and ID required")
	}

	var parentID *string
	if req.Category.ParentId != "" {
		pid := req.Category.ParentId
		parentID = &pid
	}

	cat := &model.Category{
		Model:       stdao.Model{ID: req.Category.Id},
		Name:        req.Category.Name,
		Slug:        req.Category.Slug,
		Description: req.Category.Description,
		ParentID:    parentID,
		Icon:        req.Category.Icon,
		SortOrder:   int(req.Category.SortOrder),
		IsEnabled:   req.Category.IsEnabled,
	}

	if err := dao.Category.Update(ctx, cat); err != nil {
		s.Log.Errorw("failed to update category", "err", err)
		return nil, status.Error(codes.Internal, "Failed to update category")
	}

	return &categoryV1.UpdateCategoryResponse{}, nil
}

func (s *S) DeleteCategory(ctx context.Context, req *categoryV1.DeleteCategoryRequest) (*categoryV1.DeleteCategoryResponse, error) {
	if err := requireCategoryManage(ctx, req.Id); err != nil {
		return nil, err
	}
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Category ID required")
	}

	if err := dao.Category.Delete(ctx, &model.Category{Model: stdao.Model{ID: req.Id}}).Error; err != nil {
		s.Log.Errorw("failed to delete category", "err", err)
		return nil, status.Error(codes.Internal, "Failed to delete category")
	}

	return &categoryV1.DeleteCategoryResponse{}, nil
}

func (s *S) CreateMetaTemplate(ctx context.Context, req *categoryV1.CreateMetaTemplateRequest) (*categoryV1.CreateMetaTemplateResponse, error) {
	categoryID := ""
	if req != nil && req.Template != nil {
		categoryID = req.Template.CategoryId
	}
	if err := requireCategoryManage(ctx, categoryID); err != nil {
		return nil, err
	}
	if req.Template == nil || req.Template.CategoryId == "" || req.Template.Key == "" || req.Template.Label == "" {
		return nil, status.Error(codes.InvalidArgument, "CategoryId, Key, and Label are required")
	}

	tpl := &model.CategoryMetaTemplate{
		Model:        stdao.Model{ID: ulid.Make().String()},
		CategoryID:   req.Template.CategoryId,
		Key:          req.Template.Key,
		Label:        req.Template.Label,
		Type:         req.Template.Type,
		Required:     req.Template.Required,
		Options:      req.Template.Options,
		SortOrder:    int(req.Template.SortOrder),
		DefaultValue: req.Template.DefaultValue,
	}

	if err := dao.CategoryMetaTemplate.Create(ctx, tpl); err != nil {
		s.Log.Errorw("failed to create meta template", "err", err)
		return nil, status.Error(codes.Internal, "Failed to create meta template")
	}
	return &categoryV1.CreateMetaTemplateResponse{Id: tpl.ID}, nil
}

func (s *S) UpdateMetaTemplate(ctx context.Context, req *categoryV1.UpdateMetaTemplateRequest) (*categoryV1.UpdateMetaTemplateResponse, error) {
	categoryID := ""
	if req != nil && req.Template != nil {
		categoryID = req.Template.CategoryId
	}
	if err := requireCategoryManage(ctx, categoryID); err != nil {
		return nil, err
	}
	if req.Template == nil || req.Template.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Template ID required")
	}

	tpl := &model.CategoryMetaTemplate{
		Model:        stdao.Model{ID: req.Template.Id},
		CategoryID:   req.Template.CategoryId,
		Key:          req.Template.Key,
		Label:        req.Template.Label,
		Type:         req.Template.Type,
		Required:     req.Template.Required,
		Options:      req.Template.Options,
		SortOrder:    int(req.Template.SortOrder),
		DefaultValue: req.Template.DefaultValue,
	}

	if err := dao.CategoryMetaTemplate.Update(ctx, tpl); err != nil {
		s.Log.Errorw("failed to update meta template", "err", err)
		return nil, status.Error(codes.Internal, "Failed to update meta template")
	}
	return &categoryV1.UpdateMetaTemplateResponse{}, nil
}

func (s *S) DeleteMetaTemplate(ctx context.Context, req *categoryV1.DeleteMetaTemplateRequest) (*categoryV1.DeleteMetaTemplateResponse, error) {
	if err := requireCategoryManage(ctx, ""); err != nil {
		return nil, err
	}
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Template ID required")
	}

	if err := dao.CategoryMetaTemplate.Delete(ctx, &model.CategoryMetaTemplate{Model: stdao.Model{ID: req.Id}}).Error; err != nil {
		s.Log.Errorw("failed to delete meta template", "err", err)
		return nil, status.Error(codes.Internal, "Failed to delete meta template")
	}
	return &categoryV1.DeleteMetaTemplateResponse{}, nil
}
