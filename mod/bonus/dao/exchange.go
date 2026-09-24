package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/GoldenSheep402/Hermes/mod/bonus/model"
	trafficModel "github.com/GoldenSheep402/Hermes/mod/traffic/model"
	userModel "github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	bytesPerGiB = int64(1024 * 1024 * 1024)
)

// CreditSeedingBonus adds milli-points for seeding and writes a ledger row.
func (d *bonusLog) CreditSeedingBonus(ctx context.Context, userID, torrentID string, milliPoints int64, description string) error {
	if milliPoints <= 0 || userID == "" {
		return nil
	}
	db := d.GetTxFromCtx(ctx).WithContext(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&userModel.User{}).
			Where("id = ?", userID).
			Update("bonus_points", gorm.Expr("bonus_points + ?", milliPoints))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		log := &model.BonusLog{
			Model:       stdao.Model{ID: ulid.Make().String()},
			UserID:      userID,
			Amount:      milliPoints,
			Reason:      model.BonusReasonSeeding,
			Description: description,
			RelatedID:   torrentID,
		}
		return tx.Create(log).Error
	})
}

// ExchangeUpload spends milli-points for upload credit (bytes).
func (d *bonusLog) ExchangeUpload(ctx context.Context, userID string, quantityGiB, pricePointsPerGiB int64) (remaining int64, err error) {
	if quantityGiB <= 0 || pricePointsPerGiB <= 0 {
		return 0, fmt.Errorf("invalid quantity or price")
	}
	costMilli := quantityGiB * pricePointsPerGiB * model.MilliPointsPerPoint
	uploadBytes := quantityGiB * bytesPerGiB

	db := d.GetTxFromCtx(ctx).WithContext(ctx)
	err = db.Transaction(func(tx *gorm.DB) error {
		var user userModel.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", userID).First(&user).Error; err != nil {
			return err
		}
		if user.BonusPoints < costMilli {
			return ErrInsufficientBonus
		}
		user.BonusPoints -= costMilli
		user.Uploaded += uploadBytes
		if err := tx.Model(&userModel.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
			"bonus_points": user.BonusPoints,
			"uploaded":     user.Uploaded,
		}).Error; err != nil {
			return err
		}

		ut := &trafficModel.UserTraffic{
			Model:       stdao.Model{ID: ulid.Make().String()},
			UserID:      userID,
			BonusUpload: uploadBytes,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"bonus_upload": gorm.Expr("bonus_upload + ?", uploadBytes),
			}),
		}).Create(ut).Error; err != nil {
			return err
		}

		log := &model.BonusLog{
			Model:       stdao.Model{ID: ulid.Make().String()},
			UserID:      userID,
			Amount:      -costMilli,
			Reason:      model.BonusReasonExchangeUpload,
			Description: fmt.Sprintf("exchange %d GiB upload", quantityGiB),
		}
		if err := tx.Create(log).Error; err != nil {
			return err
		}
		remaining = user.BonusPoints
		return nil
	})
	return remaining, err
}

// ExchangeInvite spends milli-points for invite inventory.
func (d *bonusLog) ExchangeInvite(ctx context.Context, userID string, quantity, pricePoints int64) (remaining int64, err error) {
	if quantity <= 0 || pricePoints <= 0 {
		return 0, fmt.Errorf("invalid quantity or price")
	}
	costMilli := quantity * pricePoints * model.MilliPointsPerPoint

	db := d.GetTxFromCtx(ctx).WithContext(ctx)
	err = db.Transaction(func(tx *gorm.DB) error {
		var user userModel.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", userID).First(&user).Error; err != nil {
			return err
		}
		if user.BonusPoints < costMilli {
			return ErrInsufficientBonus
		}
		user.BonusPoints -= costMilli
		user.InviteCount += int(quantity)
		if err := tx.Model(&userModel.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
			"bonus_points": user.BonusPoints,
			"invite_count": user.InviteCount,
		}).Error; err != nil {
			return err
		}

		log := &model.BonusLog{
			Model:       stdao.Model{ID: ulid.Make().String()},
			UserID:      userID,
			Amount:      -costMilli,
			Reason:      model.BonusReasonExchangeInvite,
			Description: fmt.Sprintf("exchange %d invite(s)", quantity),
		}
		if err := tx.Create(log).Error; err != nil {
			return err
		}
		remaining = user.BonusPoints
		return nil
	})
	return remaining, err
}

var ErrInsufficientBonus = errors.New("insufficient bonus points")
