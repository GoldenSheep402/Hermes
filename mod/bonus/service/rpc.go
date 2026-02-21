package service

import (
	"context"

	"go.uber.org/zap"

	bonusV1 "github.com/GoldenSheep402/Hermes/pkg/proto/bonus/v1"
)

var _ bonusV1.BonusServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	bonusV1.UnimplementedBonusServiceServer
}

func (s S) GetBonusBalance(ctx context.Context, request *bonusV1.GetBonusBalanceRequest) (*bonusV1.GetBonusBalanceResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) ListBonusLogs(ctx context.Context, request *bonusV1.ListBonusLogsRequest) (*bonusV1.ListBonusLogsResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) ExchangeBonus(ctx context.Context, request *bonusV1.ExchangeBonusRequest) (*bonusV1.ExchangeBonusResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) mustEmbedUnimplementedBonusServiceServer() {
	// TODO implement me
	panic("implement me")
}
