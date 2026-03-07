package service

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/bonus/dao"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	bonusV1 "github.com/GoldenSheep402/Hermes/pkg/proto/bonus/v1"
)

var _ bonusV1.BonusServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	bonusV1.UnimplementedBonusServiceServer
}

// requireAuth checks if caller is authenticated.
func requireAuth(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return "", status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return userID, nil
}

func (s S) GetBonusBalance(ctx context.Context, request *bonusV1.GetBonusBalanceRequest) (*bonusV1.GetBonusBalanceResponse, error) {
	_, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: Fetch user's bonus balance from User DAO once interconnected, or aggregate logs here.
	// For now we return 0 as a placeholder to unblock the API.
	return &bonusV1.GetBonusBalanceResponse{
		Balance: 0,
	}, nil
}

func (s S) ListBonusLogs(ctx context.Context, request *bonusV1.ListBonusLogsRequest) (*bonusV1.ListBonusLogsResponse, error) {
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}

	list, count, err := dao.BonusLog.ListByUserID(ctx, userID, request.Page, request.PageSize)
	if err != nil {
		s.Log.Errorw("failed to list bonus logs", "err", err)
		return nil, status.Error(codes.Internal, "Failed to list bonus logs")
	}

	var pList []*bonusV1.BonusLogItem
	for _, item := range list {
		pList = append(pList, &bonusV1.BonusLogItem{
			Id:          item.ID,
			UserId:      item.UserID,
			Amount:      item.Amount,
			Reason:      item.Reason,
			Description: item.Description,
			RelatedId:   item.RelatedID,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &bonusV1.ListBonusLogsResponse{
		Logs:  pList,
		Total: count,
	}, nil
}

func (s S) ExchangeBonus(ctx context.Context, request *bonusV1.ExchangeBonusRequest) (*bonusV1.ExchangeBonusResponse, error) {
	_, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}

	if request.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Amount must be greater than zero")
	}

	// TODO: Deduct bonus from User model and process exchange logic.
	return nil, status.Error(codes.Unimplemented, "Exchange logic needs User model integration")
}
