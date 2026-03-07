package service

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/torrent/dao"
	"github.com/GoldenSheep402/Hermes/mod/torrent/model"
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

	// 1. Parse torrent
	parsed, err := torrent.Parse(req.TorrentData)
	if err != nil {
		s.Log.Errorw("failed to parse torrent", "error", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid torrent file")
	}

	// 2. Validate hash doesn't already exist
	existing, err := dao.Torrent.GetByHash(ctx, parsed.InfoHash)
	if err == nil && existing != nil {
		return nil, status.Error(codes.AlreadyExists, "Torrent already exists")
	} else if err != nil && status.Code(err) != codes.NotFound {
		s.Log.Errorw("failed to check existing torrent", "error", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	// uploaderID is already authenticated via requireAuth

	// 3. Map to Model
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

	// 4. Map files
	var files []model.TorrentFile
	for _, f := range parsed.Files {
		files = append(files, model.TorrentFile{
			Model:     stdao.Model{ID: ulid.Make().String()},
			TorrentID: torrentModel.ID,
			Path:      f.Path,
			Size:      f.Size,
		})
	}

	// 5. Save to DB
	id, err := dao.Torrent.Create(ctx, torrentModel, files)
	if err != nil {
		s.Log.Errorw("failed to create torrent", "error", err)
		return nil, err // DAO returns proper grpc status
	}

	// 6. Return response
	return &torrentV1.UploadTorrentResponse{
		TorrentId: id,
		InfoHash:  parsed.InfoHash,
	}, nil
}

func (s *S) DownloadTorrent(ctx context.Context, req *torrentV1.DownloadTorrentRequest) (*torrentV1.DownloadTorrentResponse, error) {
	if _, err := requireAuth(ctx); err != nil {
		return nil, err
	}

	if req.TorrentId == "" {
		return nil, status.Error(codes.InvalidArgument, "Torrent ID cannot be empty")
	}

	// For downloading, we need the original .torrent byte stream.
	// Currently, the Hermes Torrent module architecture parses and saves metadata to DB.
	// To generate a .torrent on the fly from the DB requires pulling Tracker URL + Metadata.
	// For now, this requires extending the DB or using a generic Torrent re-generator.
	// We'll return an unimplemented error until the file storage architecture is fully clarified (e.g. MinIO vs DB blob).
	return nil, status.Error(codes.Unimplemented, "DownloadTorrent requires blob storage integration")
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
