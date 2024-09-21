package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"gorm.io/gorm"
)

type group struct {
	stdao.Std[*model.Group]
}

func (g *group) Init(db *gorm.DB) error {
	return g.Std.Init(db)
}

func (g *group) Create(ctx context.Context, group *model.Group, groupMetas []model.GroupMetadata, userID string) error {
	_ctx := g.SetTxToCtx(ctx, g.DB())
	tx := g.GetTxFromCtx(_ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(group).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range groupMetas {
		groupMetas[i].GID = group.ID
	}

	if err := tx.Create(groupMetas).Error; err != nil {
		tx.Rollback()
		return err
	}

	if userID != "" {
		var groupMembership model.GroupMembership
		groupMembership.GID = group.ID
		groupMembership.UID = userID
		if err := tx.Create(&groupMembership).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (g *group) Get(ctx context.Context, id string) (*model.Group, []model.GroupMetadata, error) {
	_ctx := g.SetTxToCtx(ctx, g.DB())
	db := g.GetTxFromCtx(_ctx)

	var group model.Group
	if err := db.Where("id = ?", id).First(&group).Error; err != nil {
		return nil, nil, err
	}

	var groupMetas []model.GroupMetadata
	if err := db.Where("gid = ?", id).Find(&groupMetas).Error; err != nil {
		return nil, nil, err
	}

	return &group, groupMetas, nil
}

func (g *group) Update(ctx context.Context, group *model.Group, groupMetas []model.GroupMetadata) error {
	_ctx := g.SetTxToCtx(ctx, g.DB())
	tx := g.GetTxFromCtx(_ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingGroup model.Group
	if err := tx.Where("id = ?", group.ID).First(&existingGroup).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fmt.Errorf("group with ID %v not found", group.ID)
		}
		return err
	}

	if err := tx.Model(&existingGroup).Updates(group).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, meta := range groupMetas {
		meta.GID = group.ID
		if meta.ID == "" {
			if err := tx.Create(&meta).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			var existingMeta model.GroupMetadata
			if err := tx.Where("id = ?", meta.ID).First(&existingMeta).Error; err != nil {
				if errors.Is(gorm.ErrRecordNotFound, err) {
					tx.Rollback()
					return fmt.Errorf("group metadata with ID %v not found", meta.ID)
				}
				tx.Rollback()
				return err
			}
			if err := tx.Model(&existingMeta).Updates(meta).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (g *group) AddUser(ctx context.Context, groupID, userID string, metas []model.GroupMembershipMetadata) error {
	_ctx := g.SetTxToCtx(ctx, g.DB())
	tx := g.GetTxFromCtx(_ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingGroup model.Group
	if err := tx.Where("id = ?", groupID).First(&existingGroup).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fmt.Errorf("group with ID %v not found", groupID)
		}
		return err
	}

	var membership model.GroupMembership
	membership.GID = groupID
	membership.UID = userID
	if err := tx.Create(&membership).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range metas {
		metas[i].MembershipID = membership.ID

		var _meta model.GroupMembershipMetadata
		if err := tx.Model(&model.GroupMembershipMetadata{}).
			Where("group_metadata_id = ?", metas[i].GroupMetadataID).
			First(&_meta).Error; err != nil {
			tx.Rollback()
			return err
		}

		metas[i].Description = _meta.Description
		metas[i].Key = _meta.Key
		metas[i].Value = _meta.Value
		metas[i].Order = _meta.Order
		metas[i].Type = _meta.Type
	}

	if err := tx.CreateInBatches(&metas, len(metas)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (g *group) RemoveUser(ctx context.Context, groupID, userID string) error {
	_ctx := g.SetTxToCtx(ctx, g.DB())
	tx := g.GetTxFromCtx(_ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingGroup model.Group
	if err := tx.Where("id = ?", groupID).First(&existingGroup).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fmt.Errorf("group with ID %v not found", groupID)
		}
		return err
	}

	var membership model.GroupMembership
	if err := tx.Model(&model.GroupMembership{}).
		Where("uid = ? AND gid = ?", userID, groupID).
		First(&membership).Error; err != nil {
		return err
	}

	// delete membership
	if err := tx.Delete(&membership).Error; err != nil {
		tx.Rollback()
		return err
	}

	// delete membership metadata
	if err := tx.Model(&model.GroupMembershipMetadata{}).
		Where("id = ?", membership.ID).
		Delete(&model.GroupMembershipMetadata{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// UpdateUser updates the metadata for a user's membership in a specific group.
// It takes a context, group ID, user ID, and a slice of GroupMembershipMetadata to update.
// Returns an error if the group or membership is not found, or if there's an issue with the database operation.
func (g *group) UpdateUser(ctx context.Context, groupID, userID string, metas []model.GroupMembershipMetadata) error {
	_ctx := g.SetTxToCtx(ctx, g.DB())
	tx := g.GetTxFromCtx(_ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingGroup model.Group
	if err := tx.Where("id = ?", groupID).First(&existingGroup).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fmt.Errorf("group with ID %v not found", groupID)
		}
		return err
	}

	var membership model.GroupMembership
	if err := tx.Model(&model.GroupMembership{}).
		Where("uid = ? AND gid = ?", userID, groupID).
		First(&membership).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return fmt.Errorf("membership for user ID %v in group ID %v not found", userID, groupID)
		}
		return err
	}

	ids := make([]string, len(metas))
	for i, meta := range metas {
		ids[i] = meta.GroupMetadataID
	}
	var allMetadata []model.GroupMetadata
	if err := tx.Where("group_metadata_id IN (?)", ids).Find(&allMetadata).Error; err != nil {
		return err
	}

	// GroupMembershipMetadata.MembershipID -> GroupMembership.ID
	// GroupMembershipMetadata.GroupMetadataID -> GroupMetadata.ID
	metadataMap := make(map[string]model.GroupMetadata, len(allMetadata))
	for _, m := range allMetadata {
		metadataMap[m.ID] = m
	}

	for i := range metas {
		metas[i].MembershipID = membership.ID
		if meta, exists := metadataMap[metas[i].GroupMetadataID]; exists {
			metas[i].Description = meta.Description
			metas[i].Key = meta.Key
			metas[i].Order = meta.Order
			metas[i].Type = meta.Type
		} else {
			return fmt.Errorf("metadata with ID %v not found", metas[i].GroupMetadataID)
		}
	}

	if err := tx.Model(&model.GroupMembershipMetadata{}).
		Where("membership_id = ?", membership.ID).
		Updates(metas).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
