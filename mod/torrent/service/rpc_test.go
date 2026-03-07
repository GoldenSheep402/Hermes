package service

import (
	"context"
	"testing"

	torrentV1 "github.com/GoldenSheep402/Hermes/pkg/proto/torrent/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// A mocked/in-memory setup for testing can be integrated here.
// For now, we outline the expected behaviors using a placeholder test setup.
// If Hermes uses a global test DB init similar to `ip_test.go`, we can bootstrap that.

func TestTorrentService_GetTorrent(t *testing.T) {
	// Because dao.Torrent requires an active DB connection, we can test standard error flows or utilize an initialized test DB.
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	// Test case: Missing ID
	resp, err := s.GetTorrent(context.Background(), &torrentV1.GetTorrentRequest{Id: ""})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthenticated")

	// The remaining tests would interact directly with `dao.Torrent.Get()`
}

func TestTorrentService_ListTorrentFiles(t *testing.T) {
	s := &S{
		Log: zap.NewNop().Sugar(),
	}

	resp, err := s.ListTorrentFiles(context.Background(), &torrentV1.ListTorrentFilesRequest{TorrentId: ""})
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unauthenticated")
}
