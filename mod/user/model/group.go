package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Group struct {
	stdao.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (g *Group) BeforeCreate(_ *gorm.DB) (err error) {
	g.ID = ulid.Make().String()
	return
}
