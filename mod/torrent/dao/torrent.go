package dao

import (
	"context"
	"errors"
	"strings"

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
	if err := t.Std.Init(db); err != nil {
		return err
	}

	// Ensure info_hash is enforced as unique index.
	// Drop potential old index variants and re-create using current model tags.
	migrator := db.Migrator()
	for _, indexName := range []string{"idx_torrents_info_hash", "info_hash", "InfoHash", "torrents_info_hash_key"} {
		if migrator.HasIndex(&model.Torrent{}, indexName) {
			_ = migrator.DropIndex(&model.Torrent{}, indexName)
		}
	}
	if !migrator.HasIndex(&model.Torrent{}, "InfoHash") {
		if err := migrator.CreateIndex(&model.Torrent{}, "InfoHash"); err != nil {
			return err
		}
	}

	return nil
}

func (t *torrent) Create(
	ctx context.Context,
	torrentBase *model.Torrent,
	files []model.TorrentFile,
	rawData []byte,
	pieceHashes []string,
) (string, error) {
	if len(rawData) == 0 {
		return "", status.Error(codes.InvalidArgument, "torrent raw data is empty")
	}

	_ctx := t.SetTxToCtx(ctx, t.DB())
	tx := t.GetTxFromCtx(_ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&model.Torrent{}).Create(torrentBase).Error; err != nil {
		tx.Rollback()
		if isDuplicateInfoHashError(err) {
			return "", status.Error(codes.AlreadyExists, "Torrent already exists")
		}
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

	blob := &model.TorrentBlob{
		TorrentID: torrentBase.ID,
		RawData:   append([]byte(nil), rawData...),
	}
	if err := tx.Model(&model.TorrentBlob{}).Create(blob).Error; err != nil {
		tx.Rollback()
		return "", status.Error(codes.Internal, "Internal error")
	}

	if len(pieceHashes) > 0 {
		items := make([]model.TorrentPiece, 0, len(pieceHashes))
		for index, hash := range pieceHashes {
			text := strings.TrimSpace(strings.ToLower(hash))
			if text == "" {
				continue
			}
			items = append(items, model.TorrentPiece{
				TorrentID:  torrentBase.ID,
				PieceIndex: index,
				PieceSHA1:  text,
			})
		}

		if len(items) > 0 {
			if err := tx.Model(&model.TorrentPiece{}).Create(&items).Error; err != nil {
				tx.Rollback()
				return "", status.Error(codes.Internal, "Internal error")
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return "", status.Error(codes.Internal, "Internal error")
	}

	return torrentBase.ID, nil
}

func (t *torrent) Get(ctx context.Context, torrentID string) (*model.Torrent, []model.TorrentFile, error) {
	db := t.GetTxFromCtx(ctx).WithContext(ctx)
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

func (t *torrent) GetBase(ctx context.Context, torrentID string) (*model.Torrent, error) {
	db := t.GetTxFromCtx(ctx).WithContext(ctx)
	var entity model.Torrent
	if err := db.Model(&model.Torrent{}).Where("id = ?", torrentID).First(&entity).Error; err != nil {
		return nil, status.Error(codes.NotFound, "Torrent not found")
	}
	return &entity, nil
}

func (t *torrent) GetByHash(ctx context.Context, hash string) (*model.Torrent, error) {
	db := t.GetTxFromCtx(ctx).WithContext(ctx)
	var torrent model.Torrent
	if err := db.Model(&model.Torrent{}).Where("info_hash = ?", hash).Order("created_at ASC").First(&torrent).Error; err != nil {
		return nil, status.Error(codes.NotFound, "Torrent not found")
	}
	return &torrent, nil
}

func (t *torrent) DeleteByID(ctx context.Context, torrentID string) error {
	if strings.TrimSpace(torrentID) == "" {
		return status.Error(codes.InvalidArgument, "Torrent ID cannot be empty")
	}

	db := t.GetTxFromCtx(ctx).WithContext(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		var current model.Torrent
		if err := tx.Model(&model.Torrent{}).Where("id = ?", torrentID).First(&current).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return status.Error(codes.NotFound, "Torrent not found")
			}
			return status.Error(codes.Internal, "Internal error")
		}

		if err := tx.Model(&model.Torrent{}).Where("id = ?", torrentID).Update("is_active", false).Error; err != nil {
			return status.Error(codes.Internal, "Internal error")
		}

		if err := tx.Where("id = ?", torrentID).Delete(&model.Torrent{}).Error; err != nil {
			return status.Error(codes.Internal, "Internal error")
		}

		if err := tx.Where("torrent_id = ?", torrentID).Delete(&model.TorrentFile{}).Error; err != nil {
			return status.Error(codes.Internal, "Internal error")
		}
		if err := tx.Where("torrent_id = ?", torrentID).Delete(&model.TorrentBlob{}).Error; err != nil {
			return status.Error(codes.Internal, "Internal error")
		}
		if err := tx.Where("torrent_id = ?", torrentID).Delete(&model.TorrentPiece{}).Error; err != nil {
			return status.Error(codes.Internal, "Internal error")
		}

		return nil
	})
}

func isDuplicateInfoHashError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	text := strings.ToLower(err.Error())
	return strings.Contains(text, "duplicate key") ||
		strings.Contains(text, "unique constraint") ||
		strings.Contains(text, "unique failed") ||
		strings.Contains(text, "duplicate entry")
}
