package torrent

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeebo/bencode"
)

func TestRewriteDownloadTorrent(t *testing.T) {
	t.Parallel()

	dummyTorrent := map[string]interface{}{
		"announce":      "http://old-tracker.local/announce",
		"announce-list": [][]string{{"http://old-tracker.local/announce"}},
		"info": map[string]interface{}{
			"name":         "test.txt",
			"piece length": int64(262144),
			"pieces":       "12345678901234567890",
			"length":       int64(1024),
		},
	}

	var buf bytes.Buffer
	err := bencode.NewEncoder(&buf).Encode(dummyTorrent)
	require.NoError(t, err)

	original, err := Parse(buf.Bytes())
	require.NoError(t, err)
	require.NotEmpty(t, original.InfoHash)

	out, err := RewriteDownloadTorrent(buf.Bytes(), "https://tracker.hermes.local/announce/pk123")
	require.NoError(t, err)
	require.NotEmpty(t, out)

	rewritten, err := Parse(out)
	require.NoError(t, err)
	require.Equal(t, original.InfoHash, rewritten.InfoHash)
	require.Equal(t, "https://tracker.hermes.local/announce/pk123", rewritten.Announce)
	require.Equal(t, [][]string{{"https://tracker.hermes.local/announce/pk123"}}, rewritten.AnnounceList)
}

func TestRewriteDownloadTorrentValidation(t *testing.T) {
	t.Parallel()

	_, err := RewriteDownloadTorrent(nil, "https://tracker.hermes.local/announce/pk")
	require.Error(t, err)

	_, err = RewriteDownloadTorrent([]byte("invalid"), "")
	require.Error(t, err)
}

func TestRewriteDownloadTorrentWithTrackers(t *testing.T) {
	t.Parallel()

	dummyTorrent := map[string]interface{}{
		"announce":      "http://old-tracker.local/announce",
		"announce-list": [][]string{{"http://old-tracker.local/announce"}},
		"info": map[string]interface{}{
			"name":         "test.txt",
			"piece length": int64(262144),
			"pieces":       "12345678901234567890",
			"length":       int64(1024),
		},
	}

	var buf bytes.Buffer
	err := bencode.NewEncoder(&buf).Encode(dummyTorrent)
	require.NoError(t, err)

	out, err := RewriteDownloadTorrentWithTrackers(buf.Bytes(), []string{
		"https://tracker-a.hermes.local/announce/pk123",
		"https://tracker-b.hermes.local/announce/pk123",
	})
	require.NoError(t, err)

	rewritten, err := Parse(out)
	require.NoError(t, err)
	require.Equal(t, "https://tracker-a.hermes.local/announce/pk123", rewritten.Announce)
	require.Equal(t, [][]string{
		{"https://tracker-a.hermes.local/announce/pk123"},
		{"https://tracker-b.hermes.local/announce/pk123"},
	}, rewritten.AnnounceList)
}
