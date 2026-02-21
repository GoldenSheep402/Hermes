package service

import (
	resourceV1 "github.com/GoldenSheep402/Hermes/pkg/proto/resource/v1"

	"context"

	"go.uber.org/zap"
)

var _ resourceV1.ResourceServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	resourceV1.UnimplementedResourceServiceServer
}

func (s S) CreateResource(ctx context.Context, request *resourceV1.CreateResourceRequest) (*resourceV1.CreateResourceResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) GetResource(ctx context.Context, request *resourceV1.GetResourceRequest) (*resourceV1.GetResourceResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) ListResources(ctx context.Context, request *resourceV1.ListResourcesRequest) (*resourceV1.ListResourcesResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) UpdateResource(ctx context.Context, request *resourceV1.UpdateResourceRequest) (*resourceV1.UpdateResourceResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) DeleteResource(ctx context.Context, request *resourceV1.DeleteResourceRequest) (*resourceV1.DeleteResourceResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) ThankResource(ctx context.Context, request *resourceV1.ThankResourceRequest) (*resourceV1.ThankResourceResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) SetPromotion(ctx context.Context, request *resourceV1.SetPromotionRequest) (*resourceV1.SetPromotionResponse, error) {
	// TODO implement me
	panic("implement me")
}
