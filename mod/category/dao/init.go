package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/category/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Category             = &category{}
	CategoryMetaTemplate = &categoryMetaTemplate{}
)

type category struct {
	stdao.Std[*model.Category]
}

type categoryMetaTemplate struct {
	stdao.Std[*model.CategoryMetaTemplate]
}

func Init(db *gorm.DB, rdb *redis.Client) error {
	if err := Category.Std.Init(db); err != nil {
		return err
	}
	if err := CategoryMetaTemplate.Std.Init(db); err != nil {
		return err
	}
	return nil
}
