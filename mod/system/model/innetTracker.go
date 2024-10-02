package model

import "github.com/GoldenSheep402/Hermes/pkg/stdao"

type InnetTracker struct {
	stdao.Model
	Address string `json:"address"`
	Enable  bool   `json:"enable"`
}
