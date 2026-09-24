package dao

import (
	"fmt"
	"testing"
	"time"

	"github.com/GoldenSheep402/Hermes/mod/tracker/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type legacySnatch struct {
	stdao.Model
	TorrentID  string
	UserID     string
	Uploaded   int64
	Downloaded int64
	SeedTime   int64
	IsActive   bool
	FinishedAt *time.Time
	LastAction time.Time
}

func (legacySnatch) TableName() string {
	return "snatches"
}

func TestSnatchInitCompactsDuplicates(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%d?mode=memory&cache=shared", time.Now().UnixNano())), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&legacySnatch{}))

	firstAction := time.Now().Add(-2 * time.Hour).UTC().Truncate(time.Second)
	secondAction := time.Now().UTC().Truncate(time.Second)
	firstFinished := firstAction.Add(5 * time.Minute)
	secondFinished := secondAction.Add(-10 * time.Minute)
	require.NoError(t, db.Create([]legacySnatch{
		{
			Model:      stdao.Model{ID: "sn-1"},
			TorrentID:  "torrent-1",
			UserID:     "user-1",
			Uploaded:   100,
			Downloaded: 200,
			SeedTime:   300,
			IsActive:   true,
			FinishedAt: &firstFinished,
			LastAction: firstAction,
		},
		{
			Model:      stdao.Model{ID: "sn-2"},
			TorrentID:  "torrent-1",
			UserID:     "user-1",
			Uploaded:   250,
			Downloaded: 150,
			SeedTime:   900,
			IsActive:   false,
			FinishedAt: &secondFinished,
			LastAction: secondAction,
		},
	}).Error)

	store := &snatch{}
	require.NoError(t, store.Init(db))

	var rows []model.Snatch
	require.NoError(t, db.Find(&rows).Error)
	require.Len(t, rows, 1)
	require.Equal(t, int64(250), rows[0].Uploaded)
	require.Equal(t, int64(200), rows[0].Downloaded)
	require.Equal(t, int64(900), rows[0].SeedTime)
	require.False(t, rows[0].IsActive)
	require.NotNil(t, rows[0].FinishedAt)
	require.Equal(t, firstFinished.Unix(), rows[0].FinishedAt.UTC().Unix())
	require.Equal(t, secondAction.Unix(), rows[0].LastAction.UTC().Unix())
}
