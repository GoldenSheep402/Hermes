package service

import (
	"context"
	"testing"

	systemV1 "github.com/GoldenSheep402/Hermes/pkg/proto/system/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSystemService_Actions_NoContext(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	r1, err1 := s.GetSettings(context.Background(), &systemV1.GetSettingsRequest{})
	assert.Error(t, err1)
	assert.Nil(t, r1)
	assert.Contains(t, err1.Error(), "unauthenticated")

	r2, err2 := s.GetSetting(context.Background(), &systemV1.GetSettingRequest{})
	assert.Error(t, err2)
	assert.Nil(t, r2)
	assert.Contains(t, err2.Error(), "unauthenticated")

	r3, err3 := s.SetSettings(context.Background(), &systemV1.SetSettingsRequest{})
	assert.Error(t, err3)
	assert.Nil(t, r3)
	assert.Contains(t, err3.Error(), "unauthenticated")

	r4, err4 := s.DeleteSetting(context.Background(), &systemV1.DeleteSettingRequest{})
	assert.Error(t, err4)
	assert.Nil(t, r4)
	assert.Contains(t, err4.Error(), "unauthenticated")

	r5, err5 := s.GetSiteStats(context.Background(), &systemV1.GetSiteStatsRequest{})
	assert.Error(t, err5)
	assert.Nil(t, r5)
	assert.Contains(t, err5.Error(), "unauthenticated")
}
