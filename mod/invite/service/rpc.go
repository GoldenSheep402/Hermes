package service

import (
	inviteV1 "github.com/GoldenSheep402/Hermes/pkg/proto/invite/v1"

	"context"

	"go.uber.org/zap"
)

var _ inviteV1.InviteServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	inviteV1.UnimplementedInviteServiceServer
}

func (s S) CreateInviteCode(ctx context.Context, request *inviteV1.CreateInviteCodeRequest) (*inviteV1.CreateInviteCodeResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) ListInviteCodes(ctx context.Context, request *inviteV1.ListInviteCodesRequest) (*inviteV1.ListInviteCodesResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) ConsumeInviteCode(ctx context.Context, request *inviteV1.ConsumeInviteCodeRequest) (*inviteV1.ConsumeInviteCodeResponse, error) {
	// TODO implement me
	panic("implement me")
}
