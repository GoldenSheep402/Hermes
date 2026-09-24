package authz

import (
	"context"

	"github.com/GoldenSheep402/Hermes/mod/casbinX/rbac"
	torrentDao "github.com/GoldenSheep402/Hermes/mod/torrent/dao"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RequireTorrentSensitiveAccess allows global admins and torrent uploaders to read
// peer lists, snatches, and per-torrent tracker statistics.
func RequireTorrentSensitiveAccess(ctx context.Context, torrentID string) error {
	if torrentID == "" {
		return status.Error(codes.InvalidArgument, "Torrent ID required")
	}
	userID, ok := ctx.Value(ctxKey.UID).(string)
	if !ok || userID == "" {
		return status.Error(codes.Unauthenticated, "unauthenticated")
	}
	isAdmin, err := rbac.CasbinManager.CheckUserIsGlobalAdmin(userID)
	if err != nil {
		return status.Error(codes.Internal, "authorization check failed")
	}
	if isAdmin {
		return nil
	}
	torrent, err := torrentDao.Torrent.GetBase(ctx, torrentID)
	if err != nil {
		return err
	}
	if torrent.UploaderID == userID {
		return nil
	}
	return status.Error(codes.PermissionDenied, "not authorized to access this torrent")
}
