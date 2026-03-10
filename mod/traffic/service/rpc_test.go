package service

import (
	"context"
	"testing"

	trafficV1 "github.com/GoldenSheep402/Hermes/pkg/proto/traffic/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTrafficService_ListTransferHistory(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	resp, err := s.ListTransferHistory(context.Background(), &trafficV1.ListTransferHistoryRequest{UserId: ""})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthenticated")
}

func TestTrafficService_GetTorrentStats(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	resp, err := s.GetTorrentStats(context.Background(), &trafficV1.GetTorrentStatsRequest{TorrentId: ""})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthenticated")
}
