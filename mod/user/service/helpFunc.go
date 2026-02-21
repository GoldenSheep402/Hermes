package service

import (
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	userV1 "github.com/GoldenSheep402/Hermes/pkg/proto/user/v1"
)

// helper to convert model.User to userV1.User
func convertUserModelToProto(u *model.User) *userV1.User {
	if u == nil {
		return nil
	}
	protoUser := &userV1.User{
		Id:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		Avatar:      u.Avatar,
		IsAdmin:     u.IsAdmin,
		IsEnabled:   u.IsEnabled,
		GroupId:     u.GroupID,
		BonusPoints: u.BonusPoints,
		Uploaded:    u.Uploaded,
		Downloaded:  u.Downloaded,
		SeedTime:    u.SeedTime,
		InviteCount: int32(u.InviteCount),
		Passkey:     u.Passkey,
		CreatedAt:   u.CreatedAt.String(),
	}
	if u.LastLogin != nil {
		protoUser.LastLogin = u.LastLogin.String()
	}
	return protoUser
}
