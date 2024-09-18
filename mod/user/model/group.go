package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type Group struct {
	stdao.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}
