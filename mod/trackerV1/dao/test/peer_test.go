package test

import (
	"context"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/dao"
	"github.com/GoldenSheep402/Hermes/mod/trackerV1/model"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"strconv"
	"testing"
	"time"
)

func TestPeer_AddPeer(t *testing.T) {
	rds, mock := redismock.NewClientMock()

	dao.Init(&gorm.DB{}, rds)

	p := dao.Peer
	err := p.Init(&gorm.DB{}, rds)
	assert.NoError(t, err)

	ctx := context.Background()
	infoHash := "test_info_hash"
	peerData := &model.Peer{
		PeerID: "test_peer_id",
		IP:     "127.0.0.1",
		Port:   6881,
	}
	key := "torrent:" + infoHash
	field := peerData.PeerID
	value := peerData.IP + ":" + strconv.Itoa(peerData.Port) + ":" + strconv.FormatInt(time.Now().Unix(), 10)

	mock.ExpectHSet(key, field, value).SetVal(1)
	mock.ExpectExpire(key, time.Hour).SetVal(true)

	err = p.AddPeer(ctx, infoHash, peerData)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
