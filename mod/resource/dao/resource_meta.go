package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/resource/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type resourceMeta struct {
	stdao.Std[*model.ResourceMeta]
}
