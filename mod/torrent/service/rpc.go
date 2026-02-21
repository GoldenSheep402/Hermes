package service

import (
	torrentV1 "github.com/GoldenSheep402/Hermes/pkg/proto/torrent/v1"

	"context"

	"go.uber.org/zap"
)

// TODO: Reimplement torrent service RPC methods with new model fields.
// Old code was deeply coupled to old Torrent model (CategoryID, CreatorID,
// CreatedBy, Length, Pieces, Private, Source, Md5sum, NameUTF8, etc.)
// and old system/user DAO methods.

var _ torrentV1.TorrentServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	torrentV1.UnimplementedTorrentServiceServer
}

func (s S) UploadTorrent(ctx context.Context, request *torrentV1.UploadTorrentRequest) (*torrentV1.UploadTorrentResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) DownloadTorrent(ctx context.Context, request *torrentV1.DownloadTorrentRequest) (*torrentV1.DownloadTorrentResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) GetTorrent(ctx context.Context, request *torrentV1.GetTorrentRequest) (*torrentV1.GetTorrentResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) ListTorrentFiles(ctx context.Context, request *torrentV1.ListTorrentFilesRequest) (*torrentV1.ListTorrentFilesResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s S) DeleteTorrent(ctx context.Context, request *torrentV1.DeleteTorrentRequest) (*torrentV1.DeleteTorrentResponse, error) {
	// TODO implement me
	panic("implement me")
}
