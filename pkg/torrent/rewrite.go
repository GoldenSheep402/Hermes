package torrent

import (
	"fmt"
	"strings"
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
