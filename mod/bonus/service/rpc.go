package service

import (
	"context"
	"errors"

	"github.com/GoldenSheep402/Hermes/mod/bonus/dao"
	"github.com/GoldenSheep402/Hermes/mod/bonus/model"
	systemSetting "github.com/GoldenSheep402/Hermes/mod/system/setting"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
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

func requireAuth(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return "", status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return userID, nil
}

func (s S) GetBonusBalance(ctx context.Context, request *bonusV1.GetBonusBalanceRequest) (*bonusV1.GetBonusBalanceResponse, error) {
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}
	user, err := userDao.User.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &bonusV1.GetBonusBalanceResponse{
		Balance:       user.BonusPoints,
		BalancePoints: user.BonusPoints / model.MilliPointsPerPoint,
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

func (s S) ListShopItems(ctx context.Context, request *bonusV1.ListShopItemsRequest) (*bonusV1.ListShopItemsResponse, error) {
	if _, err := requireAuth(ctx); err != nil {
		return nil, err
	}
	return &bonusV1.ListShopItemsResponse{
		Items: []*bonusV1.ShopItem{
			{
				Id:          "upload",
				Name:        "上传量",
				PricePoints: int64(systemSetting.BonusUploadPointsPerGiBValue(ctx)),
				Unit:        "GiB",
			},
			{
				Id:          "invite",
				Name:        "邀请名额",
				PricePoints: int64(systemSetting.BonusInvitePointsValue(ctx)),
				Unit:        "invite",
			},
		},
	}, nil
}

func (s S) ExchangeBonus(ctx context.Context, request *bonusV1.ExchangeBonusRequest) (*bonusV1.ExchangeBonusResponse, error) {
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}

	if request.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be greater than zero")
	}

	var remaining int64
	switch request.Item {
	case "upload":
		remaining, err = dao.BonusLog.ExchangeUpload(ctx, userID, request.Quantity, int64(systemSetting.BonusUploadPointsPerGiBValue(ctx)))
	case "invite":
		remaining, err = dao.BonusLog.ExchangeInvite(ctx, userID, request.Quantity, int64(systemSetting.BonusInvitePointsValue(ctx)))
	default:
		return nil, status.Error(codes.InvalidArgument, "unknown shop item")
	}
	if err != nil {
		if errors.Is(err, dao.ErrInsufficientBonus) {
			return nil, status.Error(codes.FailedPrecondition, "insufficient bonus points")
		}
		s.Log.Errorw("failed to exchange bonus", "item", request.Item, "err", err)
		return nil, status.Error(codes.Internal, "failed to exchange bonus")
	}

	return &bonusV1.ExchangeBonusResponse{RemainingBalance: remaining}, nil
}
