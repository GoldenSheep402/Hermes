package dao

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/traffic/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
)

type torrentStats struct {
	stdao.Std[*model.TorrentStats]
}

func (d *torrentStats) GetByTorrentID(ctx context.Context, torrentID string) (*model.TorrentStats, error) {
	var ts model.TorrentStats
	err := d.GetTxFromCtx(ctx).WithContext(ctx).Where("torrent_id = ?", torrentID).First(&ts).Error
	if err != nil {
		return nil, err
	}
	return &ts, nil
}
