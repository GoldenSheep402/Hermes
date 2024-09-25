package dao

import (
	"context"
	"errors"
	categoryModel "github.com/GoldenSheep402/Hermes/mod/category/model"
	"github.com/GoldenSheep402/Hermes/mod/torrent/model"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type torrent struct {
	stdao.Std[*model.Torrent]
}

func (t *torrent) Init(db *gorm.DB) error {
	return t.Std.Init(db)
}

func (t *torrent) Create(ctx context.Context, torrentBase *model.Torrent, files []model.File) (string, error) {
	_ctx := t.SetTxToCtx(ctx, t.DB())
	tx := t.GetTxFromCtx(_ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&model.Torrent{}).Create(torrentBase).Error; err != nil {
		tx.Rollback()
		return "", status.Error(codes.Internal, "Internal error")
	}

	for i, _ := range files {
		files[i].TorrentID = torrentBase.ID
	}

	if err := tx.Model(&model.File{}).Create(files).Error; err != nil {
		tx.Rollback()
		return "", status.Error(codes.Internal, "Internal error")
	}

	if err := tx.Commit().Error; err != nil {
		return "", status.Error(codes.Internal, "Internal error")
	}

	return torrentBase.ID, nil
}

func (t *torrent) Get(ctx context.Context, torrentID string) (*model.Torrent, []model.File, error) {
	db := t.DB().WithContext(ctx)
	var torrent model.Torrent
	if err := db.Model(&model.Torrent{}).Where("id = ?", torrentID).First(&torrent).Error; err != nil {
		return nil, nil, status.Error(codes.NotFound, "Torrent not found")
	}

	var files []model.File
	if err := db.Model(&model.File{}).Where("torrent_id = ?", torrentID).Find(&files).Error; err != nil {
		return nil, nil, status.Error(codes.Internal, "Internal error")
	}

	return &torrent, files, nil
}

func (t *torrent) GetTorrentMetadata(ctx context.Context, torrentID string) ([]categoryModel.Metadata, error) {
	db := t.DB().WithContext(ctx)

	var torrent model.Torrent
	if err := db.Model(&model.Torrent{}).Where("id = ?", torrentID).First(&torrent).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, "Torrent not found")
		default:

			return nil, status.Error(codes.Internal, "Internal error")
		}
	}

	var metadataBase []categoryModel.Metadata
	if err := db.Model(&categoryModel.Metadata{}).Where("category_id = ?", torrent.CategoryID).Find(&metadataBase).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, "Torrent not found")
		default:

			return nil, status.Error(codes.Internal, "Internal error")
		}
	}

	// Get metadata values
	var metadataValues []model.TorrentMetadata
	if err := db.Model(&model.TorrentMetadata{}).
		Where("torrent_id = ? AND category_id = ?", torrentID, torrent.CategoryID).
		Find(&metadataValues).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, status.Error(codes.NotFound, "Torrent not found")
		default:

			return nil, status.Error(codes.Internal, "Internal error")
		}
	}

	// Combine metadata base and metadata values
	metadataMap := make(map[string]categoryModel.Metadata)
	for _, metadata := range metadataBase {
		metadataMap[metadata.ID] = metadata
	}

	for _, metadataValue := range metadataValues {
		if metadata, exists := metadataMap[metadataValue.MetadataID]; exists {
			metadata.Value = metadataValue.Value
			metadataMap[metadataValue.MetadataID] = metadata
		}
	}

	// Convert the map back to a slice
	var result []categoryModel.Metadata
	for _, metadata := range metadataMap {
		result = append(result, metadata)
	}

	return result, nil
}
