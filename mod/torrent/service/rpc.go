package service

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/torrent/common"
	"github.com/GoldenSheep402/Hermes/mod/torrent/dao"
	"github.com/GoldenSheep402/Hermes/mod/torrent/model"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	torrentV1 "github.com/GoldenSheep402/Hermes/pkg/proto/torrent/v1"
	"github.com/GoldenSheep402/Hermes/pkg/stdao"
	"github.com/GoldenSheep402/Hermes/pkg/torrent"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ torrentV1.TorrentServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	torrentV1.UnimplementedTorrentServiceServer
}

// requireAuth checks if caller is authenticated.
func requireAuth(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return "", status.Error(codes.Unauthenticated, "unauthenticated")
	}
	return userID, nil
}

func (s *S) UploadTorrent(ctx context.Context, req *torrentV1.UploadTorrentRequest) (*torrentV1.UploadTorrentResponse, error) {
	uploaderID, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}

	if len(req.TorrentData) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Empty torrent data")
	}

	// 1. Parse original torrent first.
	// InfoHash is derived from the info dict and should not depend on announce URLs.
	parsed, err := torrent.Parse(req.TorrentData)
	if err != nil {
		s.Log.Errorw("failed to parse torrent", "error", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid torrent file")
	}

	uploader, err := userDao.User.GetByID(ctx, uploaderID)
	if err != nil || uploader == nil || uploader.Passkey == "" {
		return nil, status.Error(codes.Unauthenticated, "invalid user passkey")
	}

	announceURLs, err := common.BuildAnnounceURLsForPasskey(ctx, uploader.Passkey)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "tracker endpoint not available")
	}

	// Rewrite tracker URLs only for persisted raw bytes.
	normalizedData, err := torrent.RewriteDownloadTorrentWithTrackers(req.TorrentData, announceURLs)
	if err != nil {
		s.Log.Errorw("failed to rewrite torrent at upload", "error", err)
		return nil, status.Error(codes.Internal, "Failed to normalize torrent file")
	}

	// uploaderID is already authenticated via requireAuth

	// 2. Map to Model (from original parse result)
	torrentModel := &model.Torrent{
		Model:        stdao.Model{ID: ulid.Make().String()},
		InfoHash:     parsed.InfoHash,
		UploaderID:   uploaderID,
		Name:         parsed.Name,
		Size:         parsed.Size,
		PieceLength:  parsed.PieceLength,
		PieceCount:   parsed.PieceCount,
		IsSingleFile: len(parsed.Files) == 1,
		FileCount:    len(parsed.Files),
		IsActive:     true,
	}

	// 3. Map files
	var files []model.TorrentFile
	for _, f := range parsed.Files {
		files = append(files, model.TorrentFile{
			Model:     stdao.Model{ID: ulid.Make().String()},
			TorrentID: torrentModel.ID,
			Path:      f.Path,
			Size:      f.Size,
		})
	}

	// 4. Save to DB
	id, err := dao.Torrent.Create(ctx, torrentModel, files, normalizedData, parsed.PieceHashes)
	if err != nil {
		s.Log.Errorw("failed to create torrent", "error", err)
		return nil, err // DAO returns proper grpc status
	}

	// 5. Return response
	return &torrentV1.UploadTorrentResponse{
		TorrentId: id,
		InfoHash:  parsed.InfoHash,
	}, nil
}

func (s *S) DownloadTorrent(ctx context.Context, req *torrentV1.DownloadTorrentRequest) (*torrentV1.DownloadTorrentResponse, error) {
	uid, err := requireAuth(ctx)
	if err != nil {
		return nil, err
	}

	if req.TorrentId == "" {
		return nil, status.Error(codes.InvalidArgument, "Torrent ID cannot be empty")
	}

	rawData, err := dao.TorrentBlob.GetRawByTorrentID(ctx, req.TorrentId)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			if _, baseErr := dao.Torrent.GetBase(ctx, req.TorrentId); baseErr != nil {
				s.Log.Errorw("failed to get torrent for download", "id", req.TorrentId, "err", baseErr)
				return nil, status.Error(codes.NotFound, "Torrent not found")
			}
			return nil, status.Error(codes.FailedPrecondition, "Torrent raw data not available")
		}
		s.Log.Errorw("failed to get torrent blob for download", "id", req.TorrentId, "err", err)
		return nil, status.Error(codes.Internal, "Failed to download torrent")
	}

	user, err := userDao.User.GetByID(ctx, uid)
	if err != nil || user == nil || user.Passkey == "" {
		s.Log.Errorw("failed to get passkey for download", "uid", uid, "err", err)
		return nil, status.Error(codes.Unauthenticated, "invalid user passkey")
	}

	announceURLs, err := common.BuildAnnounceURLsForPasskey(ctx, user.Passkey)
	if err != nil {
		s.Log.Errorw("failed to resolve tracker endpoint for download", "uid", uid, "err", err)
		return nil, status.Error(codes.FailedPrecondition, "tracker endpoint not available")
	}

	downloadData, err := torrent.RewriteDownloadTorrentWithTrackers(rawData, announceURLs)
	if err != nil {
		s.Log.Errorw("failed to rewrite torrent download data", "id", req.TorrentId, "err", err)
		return nil, status.Error(codes.Internal, "Failed to build torrent download")
	}

	return &torrentV1.DownloadTorrentResponse{
		TorrentData: downloadData,
	}, nil
}

func (s *S) GetTorrent(ctx context.Context, req *torrentV1.GetTorrentRequest) (*torrentV1.GetTorrentResponse, error) {
	if _, err := requireAuth(ctx); err != nil {
		return nil, err
	}

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Torrent ID cannot be empty")
	}

	torrent, _, err := dao.Torrent.Get(ctx, req.Id)
	if err != nil {
		s.Log.Errorw("failed to get torrent", "id", req.Id, "err", err)
		return nil, status.Error(codes.NotFound, "Torrent not found")
	}

	return &torrentV1.GetTorrentResponse{
		Torrent: &torrentV1.TorrentInfo{
			Id:           torrent.ID,
			InfoHash:     torrent.InfoHash,
			UploaderId:   torrent.UploaderID,
			Name:         torrent.Name,
			Size:         torrent.Size,
			IsSingleFile: torrent.IsSingleFile,
			FileCount:    int32(torrent.FileCount),
			SeedCount:    0, // Will be hydrated by Tracker
			LeechCount:   0, // Will be hydrated by Tracker
			SnatchCount:  0, // Will be hydrated by Tracker
			IsActive:     torrent.IsActive,
			CreatedAt:    torrent.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

func (s *S) ListTorrentFiles(ctx context.Context, req *torrentV1.ListTorrentFilesRequest) (*torrentV1.ListTorrentFilesResponse, error) {
	if _, err := requireAuth(ctx); err != nil {
		return nil, err
	}

	if req.TorrentId == "" {
		return nil, status.Error(codes.InvalidArgument, "Torrent ID cannot be empty")
	}

	_, files, err := dao.Torrent.Get(ctx, req.TorrentId)
	if err != nil {
		s.Log.Errorw("failed to get torrent files", "id", req.TorrentId, "err", err)
		return nil, status.Error(codes.NotFound, "Torrent not found")
	}

	var protoFiles []*torrentV1.TorrentFile
	for _, f := range files {
		protoFiles = append(protoFiles, &torrentV1.TorrentFile{
			Id:        f.ID,
			TorrentId: f.TorrentID,
			Path:      f.Path,
			Size:      f.Size,
		})
	}

	return &torrentV1.ListTorrentFilesResponse{
		Files: protoFiles,
	}, nil
}

func (s *S) DeleteTorrent(ctx context.Context, req *torrentV1.DeleteTorrentRequest) (*torrentV1.DeleteTorrentResponse, error) {
	if _, err := requireAuth(ctx); err != nil {
		return nil, err
	}

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "Torrent ID cannot be empty")
	}

	// Soft delete requires dao.Torrent implementation. For now we just return Unimplemented.
	return nil, status.Error(codes.Unimplemented, "implement me")
}
