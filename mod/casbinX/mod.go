package casbinX

import (
	"fmt"
	"github.com/GoldenSheep402/Hermes/core/kernel"
	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbac"
	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbacValues"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var _ kernel.Module = (*Mod)(nil)

type Mod struct {
	kernel.UnimplementedModule
	ef *casbin.Enforcer
}

func (m *Mod) Name() string {
	return "casbinX"
}

func (m *Mod) Load(hub *kernel.Hub) error {
	var db *gorm.DB
	if hub.Load(&db) != nil {
		return fmt.Errorf("can't load gorm from kernel")
	}

	mo, _ := model.NewModelFromString(rbacValues.RbacRule)

	ad, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return err
	}

	ef, err := casbin.NewEnforcer(mo, ad)
	if err != nil {
		return err
	}

	err = ef.LoadPolicy()
	if err != nil {
		return err
	}
	ef.EnableAutoSave(true)
	m.ef = ef
	rbac.Init(m.ef)
	hub.Map(&m.ef)
	return nil
}
