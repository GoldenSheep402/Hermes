package service

import (
	"context"
	categoryDao "github.com/GoldenSheep402/Hermes/mod/category/dao"
	categoryModel "github.com/GoldenSheep402/Hermes/mod/category/model"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	categoryV1 "github.com/GoldenSheep402/Hermes/pkg/proto/category/v1"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ categoryV1.CategoryServiceServer

type S struct {
	Log *zap.SugaredLogger
	categoryV1.UnimplementedCategoryServiceServer
}

// CreateCategory creates a new category given a request context and a CreateCategoryRequest.
// It verifies the user's authentication, checks if the user is an admin, and then proceeds to create the category with its associated metadata.
// Returns a CreateCategoryResponse on success or an error otherwise.
func (s *S) CreateCategory(ctx context.Context, req *categoryV1.CreateCategoryRequest) (*categoryV1.CreateCategoryResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	isAdmin, err := userDao.User.IsAdmin(ctx, UID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	if !isAdmin {
		// TODO: rbac
		return nil, status.Error(codes.PermissionDenied, "You are not admin")
	}

	categoryBasic := categoryModel.Category{
		Name:        req.Category.Name,
		Description: req.Category.Description,
	}

	categoryMetas := make([]categoryModel.Metadata, len(req.Category.MetaData))
	orderMap := make(map[int]bool, len(req.Category.MetaData))

	for v, meta := range req.Category.MetaData {
		order := int(meta.Order)
		if orderMap[order] {
			return nil, status.Error(codes.InvalidArgument, "Order value must be unique")
		}

		orderMap[order] = true

		categoryMetas[v] = categoryModel.Metadata{
			Order:        order,
			Key:          meta.Key,
			Type:         meta.Type,
			Value:        meta.Value,
			DefaultValue: meta.DefaultValue,
		}
	}

	// TODO: rbac
	err = categoryDao.Category.Create(ctx, &categoryBasic, categoryMetas)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &categoryV1.CreateCategoryResponse{}, nil
}

// GetCategory retrieves a single category by its ID, including its metadata, after ensuring the user is authenticated.
// It returns a GetCategoryResponse containing the category details or an error if the operation fails.
func (s *S) GetCategory(ctx context.Context, req *categoryV1.GetCategoryRequest) (*categoryV1.GetCategoryResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	// TODO: rbac
	// like some hidden category that only allows certain user to see
	category, categoryMeta, err := categoryDao.Category.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	metaData := make([]*categoryV1.CategoryMetaData, len(categoryMeta))
	for i, meta := range categoryMeta {
		metaData[i] = &categoryV1.CategoryMetaData{
			Order:        int32(meta.Order),
			Key:          meta.Key,
			Type:         meta.Type,
			Value:        meta.Value,
			DefaultValue: meta.DefaultValue,
		}
	}

	return &categoryV1.GetCategoryResponse{
		Category: &categoryV1.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			MetaData:    metaData,
		},
	}, nil
}

func (s *S) UpdateCategory(ctx context.Context, req *categoryV1.UpdateCategoryRequest) (*categoryV1.UpdateCategoryResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	isAdmin, err := userDao.User.IsAdmin(ctx, UID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	if !isAdmin {
		// TODO: rbac
		return nil, status.Error(codes.PermissionDenied, "You are not admin")
	}

	categoryBasic := categoryModel.Category{
		Model: stdao.Model{
			ID: req.Category.Id,
		},
		Name:        req.Category.Name,
		Description: req.Category.Description,
	}

	categoryMetas := make([]categoryModel.Metadata, len(req.Category.MetaData))
	for v, meta := range req.Category.MetaData {
		categoryMetas[v] = categoryModel.Metadata{
			Model: stdao.Model{
				ID: meta.Id,
			},
			CategoryID:   req.Category.Id,
			Order:        int(meta.Order),
			Key:          meta.Key,
			Type:         meta.Type,
			Value:        meta.Value,
			DefaultValue: meta.DefaultValue,
		}
	}

	err = categoryDao.Category.Update(ctx, &categoryBasic, categoryMetas)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &categoryV1.UpdateCategoryResponse{}, nil
}

func (s *S) DeleteCategory(ctx context.Context, req *categoryV1.DeleteCategoryRequest) (*categoryV1.DeleteCategoryResponse, error) {
	UID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || UID == "" {
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	isAdmin, err := userDao.User.IsAdmin(ctx, UID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	if !isAdmin {
		// TODO: rbac
		return nil, status.Error(codes.PermissionDenied, "You are not admin")
	}

	err = categoryDao.Category.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &categoryV1.DeleteCategoryResponse{}, nil
}
