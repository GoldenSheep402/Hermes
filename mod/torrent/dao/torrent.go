package dao

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/torrent/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type torrent struct {
	stdao.Std[*model.Torrent]
}

var (
	ErrTorrentHashAlreadyExists = status.Error(codes.AlreadyExists, "TorrentHash already exists")
)

func (t *torrent) Init(db *gorm.DB) error {
	return t.Std.Init(db)
}

func (t *torrent) Create(ctx context.Context, torrentBase *model.Torrent, files []model.TorrentFile) (string, error) {
	_ctx := t.SetTxToCtx(ctx, t.DB())
	tx := t.GetTxFromCtx(_ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&model.Torrent{}).Where("info_hash = ?", torrentBase.InfoHash).First(&model.Torrent{}).Error; err == nil {
		tx.Rollback()
		return "", ErrTorrentHashAlreadyExists
	}

	if err := tx.Model(&model.Torrent{}).Create(torrentBase).Error; err != nil {
		tx.Rollback()
		return "", status.Error(codes.Internal, "Internal error")
	}

	if files != nil {
		for i := range files {
			files[i].TorrentID = torrentBase.ID
		}
		if err := tx.Model(&model.TorrentFile{}).Create(files).Error; err != nil {
			tx.Rollback()
			return "", status.Error(codes.Internal, "Internal error")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return "", status.Error(codes.Internal, "Internal error")
	}

	return torrentBase.ID, nil
}

func (t *torrent) Get(ctx context.Context, torrentID string) (*model.Torrent, []model.TorrentFile, error) {
	db := t.DB().WithContext(ctx)
	var torrent model.Torrent
	if err := db.Model(&model.Torrent{}).Where("id = ?", torrentID).First(&torrent).Error; err != nil {
		return nil, nil, status.Error(codes.NotFound, "Torrent not found")
	}

	var files []model.TorrentFile
	if err := db.Model(&model.TorrentFile{}).Where("torrent_id = ?", torrentID).Order("id ASC").Find(&files).Error; err != nil {
		return nil, nil, status.Error(codes.Internal, "Internal error")
	}

	return &torrent, files, nil
}

func (t *torrent) GetByHash(ctx context.Context, hash string) (*model.Torrent, error) {
	db := t.GetTxFromCtx(ctx).WithContext(ctx)
	var torrent model.Torrent
	if err := db.Model(&model.Torrent{}).Where("info_hash = ?", hash).First(&torrent).Error; err != nil {
		return nil, status.Error(codes.NotFound, "Torrent not found")
	}
	return &torrent, nil
}
