package dao

import (
	"errors"
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type groupMembership struct {
	stdao.Std[*model.GroupMembership]
}

func (g *groupMembership) Init(db *gorm.DB) error {
	return g.Std.Init(db)
}

func (g *groupMembership) Check(gid, uid string) (bool, error) {
	err := g.DB().Where("gid = ? AND uid = ?", gid, uid).First(&model.GroupMembership{}).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return false, err
		default:
			return false, err
		}
	}

	return true, nil
}
