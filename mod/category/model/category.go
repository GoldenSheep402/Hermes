package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type Category struct {
	stdao.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}
