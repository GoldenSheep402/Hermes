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
	require.Empty(t, rewritten.AnnounceList)
	require.Equal(t, poweredByHermesCommentTag, rewritten.Comment)
}

func TestRewriteDownloadTorrentValidation(t *testing.T) {
	t.Parallel()

	_, err := RewriteDownloadTorrent(nil, "https://tracker.hermes.local/announce/pk")
	require.Error(t, err)

	_, err = RewriteDownloadTorrent([]byte("invalid"), "")
	require.Error(t, err)
}

func TestRewriteUploadTorrent_EnforcePrivateAndRemoveDiscoveryFields(t *testing.T) {
	t.Parallel()

	dummyTorrent := map[string]interface{}{
		"announce":      "http://old-tracker.local/announce",
		"announce-list": [][]string{{"http://old-tracker.local/announce"}, {"http://backup.local/announce"}},
		"nodes":         []interface{}{[]interface{}{"router.bittorrent.com", int64(6881)}},
		"nodes6":        []interface{}{[]interface{}{"router.utorrent.com", int64(6881)}},
		"url-list":      "https://seed.example.com/test",
		"httpseeds":     []string{"https://seed2.example.com/test"},
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
	require.False(t, original.IsPrivate)

	out, err := RewriteUploadTorrent(buf.Bytes(), "https://tracker-a.hermes.local/announce/pk123")
	require.NoError(t, err)

	rewritten, err := Parse(out)
	require.NoError(t, err)
	require.True(t, rewritten.IsPrivate)
	require.NotEqual(t, original.InfoHash, rewritten.InfoHash)
	require.Equal(t, "https://tracker-a.hermes.local/announce/pk123", rewritten.Announce)
	require.Empty(t, rewritten.AnnounceList)

	var root map[string]interface{}
	err = bencode.NewDecoder(bytes.NewReader(out)).Decode(&root)
	require.NoError(t, err)

	_, hasAnnounceList := root["announce-list"]
	_, hasNodes := root["nodes"]
	_, hasNodes6 := root["nodes6"]
	_, hasURLList := root["url-list"]
	_, hasHTTPSeeds := root["httpseeds"]
	require.False(t, hasAnnounceList)
	require.False(t, hasNodes)
	require.False(t, hasNodes6)
	require.False(t, hasURLList)
	require.False(t, hasHTTPSeeds)
}

func TestRewriteUploadTorrentValidation(t *testing.T) {
	t.Parallel()

	_, err := RewriteUploadTorrent(nil, "https://tracker.hermes.local/announce/pk")
	require.Error(t, err)

	_, err = RewriteUploadTorrent([]byte("invalid"), "https://tracker.hermes.local/announce/pk")
	require.Error(t, err)

	_, err = RewriteUploadTorrent([]byte("invalid"), "")
	require.Error(t, err)
}
