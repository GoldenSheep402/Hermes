package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Category struct {
	stdao.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *Category) BeforeCreate(_ *gorm.DB) (err error) {
	c.ID = ulid.Make().String()

	return
}
