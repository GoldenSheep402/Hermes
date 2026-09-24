package service

import (
	"context"
	"testing"

	bonusV1 "github.com/GoldenSheep402/Hermes/pkg/proto/bonus/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestBonusService_GetBonusBalance(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	resp, err := s.GetBonusBalance(context.Background(), &bonusV1.GetBonusBalanceRequest{})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthenticated")
}

func TestBonusService_ListBonusLogs(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	resp, err := s.ListBonusLogs(context.Background(), &bonusV1.ListBonusLogsRequest{})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthenticated")
}

func TestBonusService_ExchangeBonus(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	resp, err := s.ExchangeBonus(context.Background(), &bonusV1.ExchangeBonusRequest{Item: "upload", Quantity: 1})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthenticated")
}

func TestBonusService_ListShopItems(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	resp, err := s.ListShopItems(context.Background(), &bonusV1.ListShopItemsRequest{})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthenticated")
}
