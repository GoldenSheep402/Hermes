package service

import (
	torrentV1 "github.com/GoldenSheep402/Hermes/pkg/proto/torrent/v1"
	"go.uber.org/zap"
)

var _ torrentV1.TorrentServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	torrentV1.UnimplementedTorrentServiceServer
}
