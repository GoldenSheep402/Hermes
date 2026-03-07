package dao

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/traffic/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type userTraffic struct {
	stdao.Std[*model.UserTraffic]
}

func (d *userTraffic) GetByUserID(ctx context.Context, userID string) (*model.UserTraffic, error) {
	var ut model.UserTraffic
	// Always use GetTxFromCtx for transaction safety per global architecture rules
	err := d.GetTxFromCtx(ctx).WithContext(ctx).Where("user_id = ?", userID).First(&ut).Error
	if err != nil {
		return nil, err
	}
	return &ut, nil
}
