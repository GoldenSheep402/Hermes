package service

import (
	categoryV1 "github.com/GoldenSheep402/Hermes/pkg/proto/category/v1"
	"go.uber.org/zap"
)

var _ categoryV1.CategoryServiceServer

type S struct {
	Log *zap.SugaredLogger
	categoryV1.UnimplementedCategoryServiceServer
}
