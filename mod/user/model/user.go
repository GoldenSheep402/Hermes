package model

import (
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type User struct {
	stdao.Model
	Name     string `json:"name"`
	Salt     string `gorm:"not null;column:salt"`     // 加密盐
	Password string `gorm:"not null;column:password"` // 加密密码
	Limit    int    `gorm:"column:limit"`
	Key      string `gorm:"column:key"`
	IsAdmin  bool   `gorm:"not null;column:is_admin" json:"is_admin"`
}
