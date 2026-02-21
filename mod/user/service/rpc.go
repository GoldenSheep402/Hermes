package service

import (
	userV1 "github.com/GoldenSheep402/Hermes/pkg/proto/user/v1"
	"go.uber.org/zap"
)

// TODO: Reimplement user service RPC methods with new proto message types.
// Old code referenced GetUserInfoRequest, CreateGroupRequest, GetGroupRequest, etc.
// which have been renamed in the new proto definitions.

var _ userV1.UserServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	userV1.UnimplementedUserServiceServer
}
