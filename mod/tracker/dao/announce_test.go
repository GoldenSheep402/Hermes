package dao

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	torrentModel "github.com/GoldenSheep402/Hermes/mod/torrent/model"
	trackerModel "github.com/GoldenSheep402/Hermes/mod/tracker/model"
	trafficModel "github.com/GoldenSheep402/Hermes/mod/traffic/model"
	userModel "github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/alicebob/miniredis/v2"
	"github.com/glebarez/sqlite"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestApplyAnnounceFirstCompletedDoesNotTrustFullCounters(t *testing.T) {
	trafficDAO, peerDAO, cleanup := newTrackerTrafficTestStore(t)
	defer cleanup()

	ctx := context.Background()
	peerData := testPeer("peer-complete", 320, 640, 0)

	result, err := trafficDAO.ApplyAnnounce(ctx, peerData, "completed")
	require.NoError(t, err)
	require.Zero(t, result.UploadDelta)
	require.Zero(t, result.DownloadDelta)
	require.False(t, result.StartedAt.IsZero())

	uploaded, downloaded, exists, err := trafficDAO.GetUserTotals(ctx, peerData.UserID)
	require.NoError(t, err)
	require.True(t, exists)
	require.Zero(t, uploaded)
	require.Zero(t, downloaded)

	storedPeer, err := peerDAO.Get(ctx, peerData.TorrentID, peerData.PeerID)
	require.NoError(t, err)
	require.NotNil(t, storedPeer)
	require.Equal(t, int64(320), storedPeer.Uploaded)
	require.Equal(t, int64(640), storedPeer.Downloaded)
	require.Equal(t, int64(0), storedPeer.Left)
}

func TestApplyAnnounceConcurrentUpdatesStayMonotonic(t *testing.T) {
	trafficDAO, peerDAO, cleanup := newTrackerTrafficTestStore(t)
	defer cleanup()

	ctx := context.Background()
	startedPeer := testPeer("peer-race", 0, 0, 1024)
	_, err := trafficDAO.ApplyAnnounce(ctx, startedPeer, "started")
	require.NoError(t, err)

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	for _, uploaded := range []int64{100, 200} {
		wg.Add(1)
		go func(uploaded int64) {
			defer wg.Done()
			announcePeer := testPeer("peer-race", uploaded, 0, 1024)
			announcePeer.LastAction = time.Now()
			_, applyErr := trafficDAO.ApplyAnnounce(ctx, announcePeer, "")
			errCh <- applyErr
		}(uploaded)
	}
	wg.Wait()
	close(errCh)

	for applyErr := range errCh {
		require.NoError(t, applyErr)
	}

	uploaded, downloaded, exists, err := trafficDAO.GetUserTotals(ctx, startedPeer.UserID)
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, int64(200), uploaded)
	require.Zero(t, downloaded)

	storedPeer, err := peerDAO.Get(ctx, startedPeer.TorrentID, startedPeer.PeerID)
	require.NoError(t, err)
	require.NotNil(t, storedPeer)
	require.Equal(t, int64(200), storedPeer.Uploaded)
	require.Zero(t, storedPeer.Downloaded)
}

func TestApplyAnnounceTracksTorrentAndUserActivity(t *testing.T) {
	trafficDAO, _, cleanup := newTrackerTrafficTestStore(t)
	defer cleanup()

	ctx := context.Background()
	startedAt := time.Now().Add(-2 * time.Minute)

	started := testPeer("peer-activity", 0, 0, 1024)
	started.LastAction = startedAt
	_, err := trafficDAO.ApplyAnnounce(ctx, started, "started")
	require.NoError(t, err)

	stats, err := trafficDAO.GetTorrentSnapshot(ctx, started.TorrentID)
	require.NoError(t, err)
	require.Equal(t, 0, stats.SeedCount)
	require.Equal(t, 1, stats.LeechCount)
	require.Equal(t, 0, stats.SnatchCount)

	seeding, downloading, err := trafficDAO.CountUserActive(ctx, started.UserID)
	require.NoError(t, err)
	require.Zero(t, seeding)
	require.Equal(t, int64(1), downloading)

	completed := testPeer("peer-activity", 50, 1024, 0)
	completed.LastAction = startedAt.Add(30 * time.Second)
	_, err = trafficDAO.ApplyAnnounce(ctx, completed, "completed")
	require.NoError(t, err)

	stats, err = trafficDAO.GetTorrentSnapshot(ctx, started.TorrentID)
	require.NoError(t, err)
	require.Equal(t, 1, stats.SeedCount)
	require.Equal(t, 0, stats.LeechCount)
	require.Equal(t, 1, stats.SnatchCount)

	seeding, downloading, err = trafficDAO.CountUserActive(ctx, started.UserID)
	require.NoError(t, err)
	require.Equal(t, int64(1), seeding)
	require.Zero(t, downloading)

	stopped := testPeer("peer-activity", 50, 1024, 0)
	stopped.LastAction = startedAt.Add(45 * time.Second)
	_, err = trafficDAO.ApplyAnnounce(ctx, stopped, "stopped")
	require.NoError(t, err)

	stats, err = trafficDAO.GetTorrentSnapshot(ctx, started.TorrentID)
	require.NoError(t, err)
	require.Zero(t, stats.SeedCount)
	require.Zero(t, stats.LeechCount)
	require.Equal(t, 1, stats.SnatchCount)

	seeding, downloading, err = trafficDAO.CountUserActive(ctx, started.UserID)
	require.NoError(t, err)
	require.Zero(t, seeding)
	require.Zero(t, downloading)
}

func TestApplyAnnounceAccumulatesSeedTimeAndFlushes(t *testing.T) {
	trafficDAO, _, cleanup := newTrackerTrafficTestStore(t)
	defer cleanup()

	ctx := context.Background()
	startedAt := time.Now().Add(-30 * time.Minute)

	completed := testPeer("peer-seed-time", 0, 2048, 0)
	completed.LastAction = startedAt
	result, err := trafficDAO.ApplyAnnounce(ctx, completed, "completed")
	require.NoError(t, err)
	require.Zero(t, result.SeedTimeDelta)

	seeding := testPeer("peer-seed-time", 256, 2048, 0)
	seeding.LastAction = startedAt.Add(20 * time.Minute)
	result, err = trafficDAO.ApplyAnnounce(ctx, seeding, "")
	require.NoError(t, err)
	require.Equal(t, int64((20*time.Minute)/time.Second), result.SeedTimeDelta)

	upload, download, seedTime, found, err := trafficDAO.GetUserSnapshot(ctx, completed.UserID)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, int64(256), upload)
	require.Equal(t, int64(0), download)
	require.Equal(t, int64((20*time.Minute)/time.Second), seedTime)

	flushed, err := trafficDAO.FlushPending(ctx, 50)
	require.NoError(t, err)
	require.GreaterOrEqual(t, flushed, 2)

	var user userModel.User
	require.NoError(t, trafficDAO.db.WithContext(ctx).Where("id = ?", completed.UserID).First(&user).Error)
	require.Equal(t, int64((20*time.Minute)/time.Second), user.SeedTime)

	var history trafficModel.TransferHistory
	require.NoError(t, trafficDAO.db.WithContext(ctx).Where("user_id = ? AND torrent_id = ?", completed.UserID, completed.TorrentID).First(&history).Error)
	require.Equal(t, int64((20*time.Minute)/time.Second), history.SeedTime)
	require.True(t, history.IsFinished)

	var snatch trackerModel.Snatch
	require.NoError(t, trafficDAO.db.WithContext(ctx).Where("torrent_id = ? AND user_id = ?", completed.TorrentID, completed.UserID).First(&snatch).Error)
	require.Equal(t, int64((20*time.Minute)/time.Second), snatch.SeedTime)
	require.True(t, snatch.IsActive)
}

func TestApplyAnnounceTracksSiteTotals(t *testing.T) {
	trafficDAO, _, cleanup := newTrackerTrafficTestStore(t)
	defer cleanup()

	ctx := context.Background()

	started := testPeer("peer-site", 128, 256, 1024)
	_, err := trafficDAO.ApplyAnnounce(ctx, started, "started")
	require.NoError(t, err)

	upload, download, found, err := trafficDAO.GetSiteTotals(ctx)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, int64(128), upload)
	require.Equal(t, int64(256), download)
}

func TestEnsureTorrentTotalsInitializesMissingSnatchCountFromTorrent(t *testing.T) {
	trafficDAO, _, cleanup := newTrackerTrafficTestStore(t)
	defer cleanup()

	ctx := context.Background()
	require.NoError(t, trafficDAO.db.WithContext(ctx).
		Model(&torrentModel.Torrent{}).
		Where("id = ?", "torrent-test").
		Update("snatch_count", 9).Error)

	require.NoError(t, trafficDAO.rds.HSet(ctx, torrentTrafficKey("torrent-test"),
		"total_upload", 123,
		"total_download", 456,
		trafficVersionField, 0,
	).Err())
	require.NoError(t, trafficDAO.rds.SAdd(ctx, trafficDirtyTorrentsKey(), "torrent-test").Err())

	stats, err := trafficDAO.GetTorrentSnapshot(ctx, "torrent-test")
	require.NoError(t, err)
	require.Equal(t, 9, stats.SnatchCount)

	flushed, err := trafficDAO.FlushPending(ctx, 10)
	require.NoError(t, err)
	require.GreaterOrEqual(t, flushed, 1)

	var torrent torrentModel.Torrent
	require.NoError(t, trafficDAO.db.WithContext(ctx).Where("id = ?", "torrent-test").First(&torrent).Error)
	require.Equal(t, 9, torrent.SnatchCount)
}

func newTrackerTrafficTestStore(t *testing.T) (*traffic, *peer, func()) {
	t.Helper()

	mr, err := miniredis.Run()
	require.NoError(t, err)

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", ulid.Make().String())), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&userModel.User{},
		&torrentModel.Torrent{},
		&trackerModel.Snatch{},
		&trafficModel.UserTraffic{},
		&trafficModel.TorrentStats{},
		&trafficModel.TransferHistory{},
	))
	require.NoError(t, db.Create(&userModel.User{
		Model:    stdao.Model{ID: "user-test"},
		Username: "user-test",
		Email:    "user-test@example.com",
		Password: "password",
		Salt:     "salt",
		Passkey:  "passkey-user-test",
	}).Error)
	require.NoError(t, db.Create(&torrentModel.Torrent{
		Model:       stdao.Model{ID: "torrent-test"},
		InfoHash:    "hash-test",
		UploaderID:  "user-test",
		Name:        "torrent-test",
		PieceLength: 16384,
		PieceCount:  1,
		FileCount:   1,
		IsActive:    true,
	}).Error)

	trafficDAO := &traffic{}
	require.NoError(t, trafficDAO.Init(db, rdb))

	peerDAO := &peer{}
	require.NoError(t, peerDAO.Init(db, rdb))

	cleanup := func() {
		_ = rdb.Close()
		mr.Close()
	}
	return trafficDAO, peerDAO, cleanup
}

func testPeer(peerID string, uploaded, downloaded, left int64) *trackerModel.Peer {
	now := time.Now()
	return &trackerModel.Peer{
		Model:      stdao.Model{ID: ulid.Make().String()},
		TorrentID:  "torrent-test",
		UserID:     "user-test",
		PeerID:     peerID,
		IP:         "127.0.0.1",
		Port:       6881,
		Uploaded:   uploaded,
		Downloaded: downloaded,
		Left:       left,
		IsSeeder:   left == 0,
		StartedAt:  now,
		LastAction: now,
	}
}
