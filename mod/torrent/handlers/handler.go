package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/GoldenSheep402/Hermes/mod/torrent/common"
	"github.com/GoldenSheep402/Hermes/mod/torrent/dao"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/pkg/auth"
	"github.com/GoldenSheep402/Hermes/pkg/ctxKey"
	"github.com/GoldenSheep402/Hermes/pkg/torrent"
	"github.com/juanjiTech/jin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Registry(jinE *jin.Engine) {
	torrentGroup := jinE.Group("/api/torrent")

	torrentGroup.Group("/download").
		GET("/:id", DownloadTorrent)
}

func DownloadTorrent(c *jin.Context) {
	torrentID := strings.TrimSpace(c.Params.ByName("id"))
	if torrentID == "" {
		writeJSONError(c, http.StatusBadRequest, int32(codes.InvalidArgument), "Torrent ID cannot be empty")
		return
	}

	ctx, authErr := buildAuthContext(c.Request.Context(), c.Request.Header.Get("Authorization"))
	if authErr != nil {
		writeJSONError(c, http.StatusUnauthorized, int32(codes.Unauthenticated), "unauthenticated")
		return
	}

	torrentBase, err := dao.Torrent.GetBase(ctx, torrentID)
	if err != nil {
		writeJSONError(c, http.StatusNotFound, int32(codes.NotFound), "Torrent not found")
		return
	}

	rawData, err := dao.TorrentBlob.GetRawByTorrentID(ctx, torrentID)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			writeJSONError(c, http.StatusPreconditionFailed, int32(codes.FailedPrecondition), "Torrent raw data not available")
			return
		}
		writeJSONError(c, http.StatusInternalServerError, int32(codes.Internal), "failed to load torrent raw data")
		return
	}

	uid, _ := ctx.Value(ctxKey.UID).(string)
	if strings.TrimSpace(uid) == "" {
		writeJSONError(c, http.StatusUnauthorized, int32(codes.Unauthenticated), "unauthenticated")
		return
	}

	user, err := userDao.User.GetByID(ctx, uid)
	if err != nil || user == nil || strings.TrimSpace(user.Passkey) == "" {
		writeJSONError(c, http.StatusUnauthorized, int32(codes.Unauthenticated), "invalid user passkey")
		return
	}

	announceURLs, err := common.BuildAnnounceURLsForPasskey(ctx, user.Passkey)
	if err != nil {
		writeJSONError(c, http.StatusPreconditionFailed, int32(codes.FailedPrecondition), "tracker endpoint not available")
		return
	}

	downloadData, err := torrent.RewriteDownloadTorrentWithTrackers(rawData, announceURLs)
	if err != nil {
		writeJSONError(c, http.StatusInternalServerError, int32(codes.Internal), "failed to rewrite torrent for download")
		return
	}

	fileName := buildDownloadFileName(torrentBase.Name, torrentID)
	c.Writer.Header().Set("Content-Type", "application/x-bittorrent")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	c.Writer.Header().Set("Content-Length", strconv.Itoa(len(downloadData)))
	c.Writer.WriteHeader(http.StatusOK)
	_, _ = c.Writer.Write(downloadData)
}

func buildAuthContext(ctx context.Context, authHeader string) (context.Context, error) {
	token := parseBearerToken(authHeader)
	if token == "" {
		return ctx, status.Error(codes.Unauthenticated, "missing bearer token")
	}

	claims, err := auth.ParseToken(token)
	if err != nil || claims.Info.UID == "" {
		return ctx, status.Error(codes.Unauthenticated, "invalid bearer token")
	}

	newCtx := context.WithValue(ctx, ctxKey.UID, claims.Info.UID)
	newCtx = context.WithValue(newCtx, ctxKey.OrgID, claims.Info.OrgID)
	return newCtx, nil
}

func parseBearerToken(authHeader string) string {
	header := strings.TrimSpace(authHeader)
	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}

	if strings.ToLower(strings.TrimSpace(parts[0])) != "bearer" {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

func buildDownloadFileName(name string, torrentID string) string {
	baseName := strings.TrimSpace(name)
	if baseName == "" {
		baseName = torrentID
	}

	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	baseName = strings.TrimSpace(replacer.Replace(baseName))
	if baseName == "" {
		baseName = torrentID
	}
	if !strings.HasSuffix(strings.ToLower(baseName), ".torrent") {
		baseName += ".torrent"
	}
	return baseName
}

func writeJSONError(c *jin.Context, httpStatus int, grpcCode int32, message string) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(httpStatus)
	payload := map[string]interface{}{
		"code":    grpcCode,
		"message": message,
		"details": []interface{}{},
	}
	data, _ := json.Marshal(payload)
	_, _ = c.Writer.Write(data)
}
