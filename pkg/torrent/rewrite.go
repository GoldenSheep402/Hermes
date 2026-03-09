package torrent

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/zeebo/bencode"
)

// RewriteDownloadTorrent parses raw torrent bytes, replaces tracker announce fields,
// and re-encodes the torrent for client download.
func RewriteDownloadTorrent(rawData []byte, announceURL string) ([]byte, error) {
	return RewriteDownloadTorrentWithTrackers(rawData, []string{announceURL})
}

// RewriteDownloadTorrentWithTrackers parses raw torrent bytes, replaces tracker announce fields
// with a prioritized list, and re-encodes the torrent for client download.
func RewriteDownloadTorrentWithTrackers(rawData []byte, announceURLs []string) ([]byte, error) {
	if len(rawData) == 0 {
		return nil, fmt.Errorf("empty torrent data")
	}

	normalized, err := normalizeAnnounceURLs(announceURLs)
	if err != nil {
		return nil, err
	}

	parsed, err := Parse(rawData)
	if err != nil {
		return nil, err
	}

	parsed.Announce = normalized[0]
	announceList := make([][]string, 0, len(normalized))
	for _, trackerURL := range normalized {
		announceList = append(announceList, []string{trackerURL})
	}
	parsed.AnnounceList = announceList

	return parsed.Marshal()
}

// RewriteUploadTorrentWithTrackers rewrites tracker fields for upload persistence and enforces
// private tracker mode:
// - sets info.private = 1
// - removes DHT / external peer-discovery related fields from root dict
func RewriteUploadTorrentWithTrackers(rawData []byte, announceURLs []string) ([]byte, error) {
	if len(rawData) == 0 {
		return nil, fmt.Errorf("empty torrent data")
	}

	normalized, err := normalizeAnnounceURLs(announceURLs)
	if err != nil {
		return nil, err
	}

	root, err := decodeTorrentRoot(rawData)
	if err != nil {
		return nil, err
	}

	root["announce"] = normalized[0]
	announceList := make([][]string, 0, len(normalized))
	for _, trackerURL := range normalized {
		announceList = append(announceList, []string{trackerURL})
	}
	root["announce-list"] = announceList

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

func normalizeAnnounceURLs(announceURLs []string) ([]string, error) {
	normalized := make([]string, 0, len(announceURLs))
	for _, announceURL := range announceURLs {
		announceURL = strings.TrimSpace(announceURL)
		if announceURL == "" {
			continue
		}
		normalized = append(normalized, announceURL)
	}
	if len(normalized) == 0 {
		return nil, fmt.Errorf("empty announce url")
	}
	return normalized, nil
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
