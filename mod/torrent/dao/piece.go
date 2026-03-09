package dao

import (
	"context"
	"strings"

	"github.com/GoldenSheep402/Hermes/mod/torrent/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type torrentPiece struct {
	stdao.Std[*model.TorrentPiece]
}

func (d *torrentPiece) Init(db *gorm.DB) error {
	return d.Std.Init(db)
}

func (d *torrentPiece) ReplaceByTorrentID(ctx context.Context, torrentID string, pieceHashes []string) error {
	if torrentID == "" {
		return status.Error(codes.InvalidArgument, "torrent id is empty")
	}

	db := d.GetTxFromCtx(ctx).WithContext(ctx)
	if err := db.Where("torrent_id = ?", torrentID).Delete(&model.TorrentPiece{}).Error; err != nil {
		return status.Error(codes.Internal, "failed to clear torrent pieces")
	}

	if len(pieceHashes) == 0 {
		return nil
	}

	items := make([]model.TorrentPiece, 0, len(pieceHashes))
	for index, hash := range pieceHashes {
		text := strings.TrimSpace(strings.ToLower(hash))
		if text == "" {
			continue
		}
		items = append(items, model.TorrentPiece{
			TorrentID:  torrentID,
			PieceIndex: index,
			PieceSHA1:  text,
		})
	}

	if len(items) == 0 {
		return nil
	}

	if err := db.Create(&items).Error; err != nil {
		return status.Error(codes.Internal, "failed to save torrent pieces")
	}

	return nil
}
