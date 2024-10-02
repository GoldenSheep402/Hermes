package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Bind                    = &bind{}
	User                    = &user{}
	Group                   = &group{}
	GroupMetadata           = &groupMetadata{}
	GroupMembership         = &groupMembership{}
	GroupMembershipMetadata = &groupMembershipMetadata{}
)

func Init(DB *gorm.DB, rds *redis.Client) error {
	err := Bind.Init(DB)
	if err != nil {
		return err
	}

	err = User.Init(DB, rds)
	if err != nil {
		return err
	}

	err = Group.Init(DB)
	if err != nil {
		return err
	}

	err = GroupMetadata.Init(DB)
	if err != nil {
		return err
	}

	err = GroupMembership.Init(DB)
	if err != nil {
		return err
	}

	err = GroupMembershipMetadata.Init(DB)
	if err != nil {
		return err
	}

	// err = DB.Transaction(func(tx *gorm.DB) error {
	// 	presetGroups := []struct {
	// 		Name        string
	// 		Description string
	// 		Level       int
	// 	}{
	// 		{
	// 			Name:        "normalUser",
	// 			Description: "normal user",
	// 			Level:       rbacValues.Level10,
	// 		},
	// 		{
	// 			Name:        "rbac",
	// 			Description: "rbac user",
	// 			Level:       rbacValues.Level1,
	// 		},
	// 		{
	// 			Name:        "admin",
	// 			Description: "admin",
	// 			Level:       rbacValues.Level0,
	// 		},
	// 	}
	//
	// 	for _, pg := range presetGroups {
	// 		var count int64
	// 		if err := tx.Model(&model.Group{}).Where("name = ?", pg.Name).Count(&count).Error; err != nil {
	// 			return err
	// 		}
	// 		var idMap map[string]string
	// 		if count == 0 {
	// 			group := model.Group{
	// 				Name:        pg.Name,
	// 				Description: pg.Description,
	// 			}
	// 			if err := tx.Create(&group).Error; err != nil {
	// 				return err
	// 			}
	// 			groupMetadata := model.GroupMetadata{
	// 				GID:         group.ID,
	// 				Key:         "Level",
	// 				Value:       strconv.Itoa(pg.Level),
	// 				Type:        "number",
	// 				Description: "It is the level of the group.",
	// 				Order:       0,
	// 			}
	// 			if err := tx.Create(&groupMetadata).Error; err != nil {
	// 				return err
	// 			}
	// 			idMap[pg.Name] = group.ID
	// 		}
	//
	// 		err := rbac.CasbinManager.SetSubgroup(idMap["admin"], idMap["rbac"])
	// 		if err != nil {
	// 			return err
	// 		}
	//
	// 		err = rbac.CasbinManager.SetSubgroup(idMap["rbac"], idMap["normalUser"])
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	//
	// 	return nil
	// })
	//
	// if err != nil {
	// 	return err
	// }

	return nil
}
