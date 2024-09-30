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

var (
	ErrTorrentHashAlreadyExists = status.Error(codes.AlreadyExists, "TorrentHash already exists")
)

func (t *torrent) Init(db *gorm.DB) error {
	return t.Std.Init(db)
}

func (t *torrent) Create(ctx context.Context, torrentBase *model.Torrent, files []model.File, metas []model.TorrentMetadata) (string, error) {
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

		if err := tx.Model(&model.File{}).Create(files).Error; err != nil {
			tx.Rollback()
			return "", status.Error(codes.Internal, "Internal error")
		}
	}

	for i := range metas {
		metas[i].TorrentID = torrentBase.ID
	}

	metaIds := make([]string, 0, len(metas))
	for i := range metas {
		metaIds = append(metaIds, metas[i].MetadataID)
	}

	var metasBase []categoryModel.Metadata
	if err := tx.Model(&categoryModel.Metadata{}).Where("id IN ?", metaIds).Find(&metasBase).Error; err != nil {
		tx.Rollback()
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return "", status.Error(codes.NotFound, "Metadata not found")
		default:
			return "", status.Error(codes.Internal, "Internal error")
		}
	} else {
		if err := tx.Model(&model.TorrentMetadata{}).Create(metas).Error; err != nil {
			tx.Rollback()
			return "", status.Error(codes.Internal, "Internal error")
		}
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
	if err := db.Model(&model.File{}).Where("torrent_id = ?", torrentID).Order("id ASC").Find(&files).Order("id").Error; err != nil {
		return nil, nil, status.Error(codes.Internal, "Internal error")
	}

	//for i := range files {
	//	if files[i].PathUTF8 == "" {
	//		files[i].PathUTF8 = files[i].Path
	//	}
	//}

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

func (t *torrent) GetTorrentList(ctx context.Context, categoryID string, torrentId string, limit int) ([]model.Torrent, error) {
	db := t.GetTxFromCtx(ctx).WithContext(ctx)
	var torrents []model.Torrent
	if categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}
	if torrentId != "" {
		db = db.Where("id < ?", torrentId)
	}
	if limit <= 0 || limit > 1000 {
		db = db.Limit(1000)
	} else {
		db = db.Limit(limit)
	}

	db = db.Order("created_at DESC")
	if err := db.Find(&torrents).Error; err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return torrents, nil
}

// GetTorrentByHash get torrent by hash
// TODO: add cache
func (t *torrent) GetTorrentByHash(ctx context.Context, hash string) (*model.Torrent, error) {
	db := t.GetTxFromCtx(ctx).WithContext(ctx)
	var torrent model.Torrent
	if err := db.Model(&model.Torrent{}).Where("info_hash = ?", hash).First(&torrent).Error; err != nil {
		return nil, status.Error(codes.NotFound, "Torrent not found")
	}

	return &torrent, nil
}
