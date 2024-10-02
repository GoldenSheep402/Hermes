package model

import "github.com/GoldenSheep402/Hermes/pkg/stdao"

type Setting struct {
	stdao.Model
	PeerExpireTime int    `gorm:"column:peer_expire_time;default:60"`
	SmtpEnable     bool   `gorm:"column:smtp_enable;default:false"`
	SmtpHost       string `gorm:"column:smtp_host"`
	SmtpPort       int    `gorm:"column:smtp_port"`
	SmtpUser       string `gorm:"column:smtp_user"`
	SmtpPass       string `gorm:"column:smtp_pass"`
	RegisterEnable bool   `gorm:"column:register_enable;default:true"`
	LoginEnable    bool   `gorm:"column:login_enable;default:true"`
	PublishEnable  bool   `gorm:"column:publish_enable;default:true"`
}
