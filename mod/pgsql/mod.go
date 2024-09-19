package pgsql

import (
	"errors"
	"fmt"
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule // 请为所有Module引入UnimplementedModule
	config                     Config
}

type Config struct {
	Addr     string `yaml:"Addr"`
	PORT     string `yaml:"Port"`
	USER     string `yaml:"User"`
	PASSWORD string `yaml:"Password"`
	DATABASE string `yaml:"Database"`
}

func (m *Mod) Config() any {
	return &m.config
}

func (m *Mod) Name() string {
	return "Pgsql"
}

func (m *Mod) PreInit(hub *kernel.Hub) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Shanghai",
		m.config.Addr, m.config.PORT, m.config.USER, m.config.DATABASE, m.config.PASSWORD)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	hub.Log.Info("pgsql init success")
	hub.Map(&db)
	return nil
}

func (m *Mod) Init(hub *kernel.Hub) error {
	var db *gorm.DB
	if hub.Load(&db) != nil {
		return errors.New("can't load gorm from kernel")
	}

	var tables []string
	result := db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tables)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
