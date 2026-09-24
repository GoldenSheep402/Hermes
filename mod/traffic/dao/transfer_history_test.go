package dao

import (
	"fmt"
	"testing"
	"time"

	"github.com/GoldenSheep402/Hermes/mod/traffic/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type legacyTransferHistory struct {
	stdao.Model
	UserID     string
	TorrentID  string
	Uploaded   int64
	Downloaded int64
	SeedTime   int64
	IsFinished bool
	IsActive   bool
	LastAction time.Time
}

func (legacyTransferHistory) TableName() string {
	return "transfer_histories"
}

func TestTransferHistoryInitCompactsDuplicates(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%d?mode=memory&cache=shared", time.Now().UnixNano())), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&legacyTransferHistory{}))

	firstAction := time.Now().Add(-1 * time.Hour).UTC().Truncate(time.Second)
	secondAction := time.Now().UTC().Truncate(time.Second)
	require.NoError(t, db.Create([]legacyTransferHistory{
		{
			Model:      stdao.Model{ID: "th-1"},
			UserID:     "user-1",
			TorrentID:  "torrent-1",
			Uploaded:   100,
			Downloaded: 200,
			SeedTime:   300,
			IsActive:   true,
			LastAction: firstAction,
		},
		{
			Model:      stdao.Model{ID: "th-2"},
			UserID:     "user-1",
			TorrentID:  "torrent-1",
			Uploaded:   250,
			Downloaded: 150,
			SeedTime:   900,
			IsFinished: true,
			LastAction: secondAction,
		},
	}).Error)

	store := &transferHistory{}
	require.NoError(t, store.Init(db))

	var rows []model.TransferHistory
	require.NoError(t, db.Find(&rows).Error)
	require.Len(t, rows, 1)
	require.Equal(t, int64(250), rows[0].Uploaded)
	require.Equal(t, int64(200), rows[0].Downloaded)
	require.Equal(t, int64(900), rows[0].SeedTime)
	require.True(t, rows[0].IsFinished)
	require.False(t, rows[0].IsActive)
	require.Equal(t, secondAction.Unix(), rows[0].LastAction.UTC().Unix())
}
