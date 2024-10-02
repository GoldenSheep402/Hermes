package model

import "github.com/GoldenSheep402/Hermes/pkg/stdao"

type Subnet struct {
	stdao.Model
	CIDR    string `gorm:"column:cidr"`
	IsAllow bool   `gorm:"column:is_allow"`
}
