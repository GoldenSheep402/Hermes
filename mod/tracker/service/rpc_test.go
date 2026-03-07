package service

import (
	"context"
	"testing"

	trackerV1 "github.com/GoldenSheep402/Hermes/pkg/proto/tracker/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTrackerService_GetTorrentPeers(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	resp, err := s.GetTorrentPeers(context.Background(), &trackerV1.GetTorrentPeersRequest{TorrentId: ""})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthenticated")
}

func TestTrackerService_ListSnatches(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	resp, err := s.ListSnatches(context.Background(), &trackerV1.ListSnatchesRequest{TorrentId: ""})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthenticated")
}

func TestTrackerService_GetUserSnatches(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	resp, err := s.GetUserSnatches(context.Background(), &trackerV1.GetUserSnatchesRequest{UserId: ""})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthenticated")
}
