package dao

import (
	"github.com/GoldenSheep402/Hermes/mod/resource/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type tag struct {
	stdao.Std[*model.Tag]
}

type resourceTag struct {
	stdao.Std[*model.ResourceTag]
}
