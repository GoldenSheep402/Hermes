package service

import (
	"context"
	"testing"

	inviteV1 "github.com/GoldenSheep402/Hermes/pkg/proto/invite/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInviteService_Actions_NoContext(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	r1, err1 := s.CreateInviteCode(context.Background(), &inviteV1.CreateInviteCodeRequest{})
	assert.Error(t, err1)
	assert.Nil(t, r1)
	assert.Contains(t, err1.Error(), "unauthenticated")

	r2, err2 := s.ListInviteCodes(context.Background(), &inviteV1.ListInviteCodesRequest{})
	assert.Error(t, err2)
	assert.Nil(t, r2)
	assert.Contains(t, err2.Error(), "unauthenticated")

	r3, err3 := s.ConsumeInviteCode(context.Background(), &inviteV1.ConsumeInviteCodeRequest{})
	assert.Error(t, err3)
	assert.Nil(t, r3)
	assert.Contains(t, err3.Error(), "unauthenticated")
}
