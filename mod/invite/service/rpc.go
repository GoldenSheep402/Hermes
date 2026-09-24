package service

import (
	inviteV1 "github.com/GoldenSheep402/Hermes/pkg/proto/invite/v1"

	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/GoldenSheep402/Hermes/mod/invite/dao"
	"github.com/GoldenSheep402/Hermes/mod/invite/model"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	userModel "github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var _ inviteV1.InviteServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	inviteV1.UnimplementedInviteServiceServer
}

func requireAuth(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return "", status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return userID, nil
}

func generateSecureCode() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func (s S) CreateInviteCode(ctx context.Context, req *inviteV1.CreateInviteCodeRequest) (*inviteV1.CreateInviteCodeResponse, error) {
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}

	count := req.Count
	if count <= 0 || count > 10 {
		count = 1 // Default to 1, max 10
	}

	user, err := userDao.User.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if !user.IsAdmin && user.InviteCount < int(count) {
		return nil, status.Error(codes.FailedPrecondition, "insufficient invite quota")
	}

	var codesList []string
	err = dao.InviteCode.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if !user.IsAdmin {
			res := tx.Model(&userModel.User{}).
				Where("id = ? AND invite_count >= ?", userID, count).
				Update("invite_count", gorm.Expr("invite_count - ?", count))
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return status.Error(codes.FailedPrecondition, "insufficient invite quota")
			}
		}

		for i := 0; i < int(count); i++ {
			codeStr := generateSecureCode()
			invite := &model.InviteCode{
				Model:     stdao.Model{ID: ulid.Make().String()},
				Code:      codeStr,
				SenderID:  userID,
				IsUsed:    false,
				ExpiredAt: time.Now().Add(72 * time.Hour),
			}
			if err := tx.Create(invite).Error; err != nil {
				return err
			}
			codesList = append(codesList, codeStr)
		}
		return nil
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return nil, st.Err()
		}
		s.Log.Errorw("failed to create invite code", "err", err)
		return nil, status.Error(codes.Internal, "Failed to allocate invite codes")
	}

	return &inviteV1.CreateInviteCodeResponse{Codes: codesList}, nil
}

func (s S) ListInviteCodes(ctx context.Context, req *inviteV1.ListInviteCodesRequest) (*inviteV1.ListInviteCodesResponse, error) {
	userID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}

	list, total, err := dao.InviteCode.ListBySender(ctx, userID, req.Page, req.PageSize)
	if err != nil {
		s.Log.Errorw("failed to list invite codes", "err", err)
		return nil, status.Error(codes.Internal, "Failed to list invite codes")
	}

	var pbCodes []*inviteV1.InviteCodeInfo
	for _, code := range list {
		var expiredStr, usedStr, recID string
		if !code.ExpiredAt.IsZero() {
			expiredStr = code.ExpiredAt.Format("2006-01-02T15:04:05Z07:00")
		}
		if code.UsedAt != nil && !code.UsedAt.IsZero() {
			usedStr = code.UsedAt.Format("2006-01-02T15:04:05Z07:00")
		}
		if code.ReceiverID != nil {
			recID = *code.ReceiverID
		}

		pbCodes = append(pbCodes, &inviteV1.InviteCodeInfo{
			Id:           code.ID,
			Code:         code.Code,
			SenderId:     code.SenderID,
			SenderName:   "User " + code.SenderID,
			ReceiverId:   recID,
			ReceiverName: "User " + recID,
			IsUsed:       code.IsUsed,
			ExpiredAt:    expiredStr,
			UsedAt:       usedStr,
			CreatedAt:    code.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &inviteV1.ListInviteCodesResponse{Codes: pbCodes, Total: total}, nil
}

func (s S) ConsumeInviteCode(ctx context.Context, req *inviteV1.ConsumeInviteCodeRequest) (*inviteV1.ConsumeInviteCodeResponse, error) {
	userID, err := requireAuth(ctx) // Often during signup, caller might not be fully authed yet, but assuming user gets a provisional token or is an admin consuming
	if err != nil {
		return nil, err
	}

	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "Code is required")
	}

	codeMod, err := dao.InviteCode.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Invite code not found or invalid")
	}

	if codeMod.IsUsed {
		return nil, status.Error(codes.FailedPrecondition, "Invite code has already been used")
	}

	if time.Now().After(codeMod.ExpiredAt) {
		return nil, status.Error(codes.FailedPrecondition, "Invite code has expired")
	}

	if err := dao.InviteCode.Consume(ctx, req.Code, userID); err != nil {
		s.Log.Errorw("failed to consume invite code", "err", err)
		return nil, status.Error(codes.Internal, "Failed to process invite consumption")
	}

	return &inviteV1.ConsumeInviteCodeResponse{Success: true}, nil
}
