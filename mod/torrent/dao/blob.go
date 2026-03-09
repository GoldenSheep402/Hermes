package dao

import (
	"context"
	"errors"

	"github.com/GoldenSheep402/Hermes/mod/torrent/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type torrentBlob struct {
	stdao.Std[*model.TorrentBlob]
}

func (d *torrentBlob) Init(db *gorm.DB) error {
	return d.Std.Init(db)
}

func (d *torrentBlob) UpsertByTorrentID(ctx context.Context, torrentID string, rawData []byte) error {
	if torrentID == "" {
		return status.Error(codes.InvalidArgument, "torrent id is empty")
	}
	if len(rawData) == 0 {
		return status.Error(codes.InvalidArgument, "torrent raw data is empty")
	}

	db := d.GetTxFromCtx(ctx).WithContext(ctx)
	update := db.Model(&model.TorrentBlob{}).Where("torrent_id = ?", torrentID).Update("raw_data", rawData)
	if update.Error != nil {
		return status.Error(codes.Internal, "failed to upsert torrent blob")
	}
	if update.RowsAffected > 0 {
		return nil
	}

	entity := &model.TorrentBlob{
		TorrentID: torrentID,
		RawData:   append([]byte(nil), rawData...),
	}
	if err := db.Create(entity).Error; err != nil {
		return status.Error(codes.Internal, "failed to create torrent blob")
	}
	return nil
}

func (d *torrentBlob) GetRawByTorrentID(ctx context.Context, torrentID string) ([]byte, error) {
	if torrentID == "" {
		return nil, status.Error(codes.InvalidArgument, "torrent id is empty")
	}

	db := d.GetTxFromCtx(ctx).WithContext(ctx)
	var entity model.TorrentBlob
	if err := db.Where("torrent_id = ?", torrentID).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Legacy fallback: try old raw_data column in torrents table
			var legacy struct {
				RawData []byte `gorm:"column:raw_data"`
			}
			legacyErr := db.Table("torrents").Select("raw_data").Where("id = ?", torrentID).Take(&legacy).Error
			if legacyErr == nil && len(legacy.RawData) > 0 {
				_ = d.UpsertByTorrentID(ctx, torrentID, legacy.RawData)
				return append([]byte(nil), legacy.RawData...), nil
			}
			return nil, status.Error(codes.NotFound, "torrent blob not found")
		}
		return nil, status.Error(codes.Internal, "failed to get torrent blob")
	}

	if len(entity.RawData) == 0 {
		return nil, status.Error(codes.NotFound, "torrent blob not found")
	}

	return append([]byte(nil), entity.RawData...), nil
}
