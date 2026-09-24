package torrent

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/zeebo/bencode"
)

const poweredByHermesCommentTag = "Powered by Hermes"

// appendPoweredByHermesComment appends a site attribution to the torrent root comment on download.
func appendPoweredByHermesComment(existing string) string {
	existing = strings.TrimSpace(existing)
	if strings.Contains(strings.ToLower(existing), "powered by hermes") {
		return existing
	}
	if existing == "" {
		return poweredByHermesCommentTag
	}
	return existing + " | " + poweredByHermesCommentTag
}

// RewriteDownloadTorrent parses raw torrent bytes, replaces the primary announce URL,
// removes announce-list (clients prefer it over announce when present), and re-encodes.
func RewriteDownloadTorrent(rawData []byte, announceURL string) ([]byte, error) {
	if len(rawData) == 0 {
		return nil, fmt.Errorf("empty torrent data")
	}

	announceURL = strings.TrimSpace(announceURL)
	if announceURL == "" {
		return nil, fmt.Errorf("empty announce url")
	}

	parsed, err := Parse(rawData)
	if err != nil {
		return nil, err
	}

	parsed.Announce = announceURL
	parsed.AnnounceList = nil
	parsed.Comment = appendPoweredByHermesComment(parsed.Comment)

	return parsed.Marshal()
}

// RewriteUploadTorrent rewrites tracker fields for upload persistence and enforces
// private tracker mode:
// - sets info.private = 1
// - removes announce-list and DHT / external peer-discovery fields
func RewriteUploadTorrent(rawData []byte, announceURL string) ([]byte, error) {
	if len(rawData) == 0 {
		return nil, fmt.Errorf("empty torrent data")
	}

	announceURL = strings.TrimSpace(announceURL)
	if announceURL == "" {
		return nil, fmt.Errorf("empty announce url")
	}

	root, err := decodeTorrentRoot(rawData)
	if err != nil {
		return nil, err
	}

	root["announce"] = announceURL
	delete(root, "announce-list")

	info, ok := root["info"].(map[string]interface{})
	if !ok || info == nil {
		return nil, fmt.Errorf("invalid torrent info dict")
	}
	info["private"] = int64(1)
	root["info"] = info

	// Remove non-private peer discovery / external source keys at root level.
	delete(root, "nodes")
	delete(root, "nodes6")
	delete(root, "dht")
	delete(root, "url-list")
	delete(root, "httpseeds")

	var buf bytes.Buffer
	if err := bencode.NewEncoder(&buf).Encode(root); err != nil {
		return nil, fmt.Errorf("failed to encode torrent data: %w", err)
	}
	return buf.Bytes(), nil
}

// RewriteDownloadTorrentWithTrackers is kept for callers that still pass a list;
// only the first non-empty URL is used.
func RewriteDownloadTorrentWithTrackers(rawData []byte, announceURLs []string) ([]byte, error) {
	url, err := firstAnnounceURL(announceURLs)
	if err != nil {
		return nil, err
	}
	return RewriteDownloadTorrent(rawData, url)
}

// RewriteUploadTorrentWithTrackers is kept for callers that still pass a list;
// only the first non-empty URL is used.
func RewriteUploadTorrentWithTrackers(rawData []byte, announceURLs []string) ([]byte, error) {
	url, err := firstAnnounceURL(announceURLs)
	if err != nil {
		return nil, err
	}
	return RewriteUploadTorrent(rawData, url)
}

func firstAnnounceURL(announceURLs []string) (string, error) {
	for _, announceURL := range announceURLs {
		announceURL = strings.TrimSpace(announceURL)
		if announceURL != "" {
			return announceURL, nil
		}
	}
	return "", fmt.Errorf("empty announce url")
}

func decodeTorrentRoot(rawData []byte) (map[string]interface{}, error) {
	var root map[string]interface{}
	if err := bencode.NewDecoder(bytes.NewReader(rawData)).Decode(&root); err != nil {
		return nil, fmt.Errorf("failed to decode torrent data: %w", err)
	}
	if len(root) == 0 {
		return nil, fmt.Errorf("invalid torrent root dict")
	}
	return root, nil
}
