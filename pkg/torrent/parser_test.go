package torrent

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeebo/bencode"
)

func TestParseSingleFile(t *testing.T) {
	// Create a dummy single-file torrent purely via bencode
	dummyTorrent := map[string]interface{}{
		"announce":      "http://tracker.example.com",
		"announce-list": [][]string{{"http://tracker1.example.com", "http://tracker2.example.com"}},
		"creation date": int64(1600000000),
		"comment":       "test comment",
		"created by":    "test client",
		"encoding":      "UTF-8",
		"url-list":      []string{"http://webseed1.example.com"},
		"info": map[string]interface{}{
			"name":         "test.txt",
			"name.utf-8":   "test-utf8.txt",
			"piece length": int64(262144),
			"pieces":       "12345678901234567890", // 20 bytes dummy hash
			"length":       int64(1024),
			"private":      1,
			"source":       "test-source",
		},
	}

	var buf bytes.Buffer
	err := bencode.NewEncoder(&buf).Encode(dummyTorrent)
	assert.NoError(t, err)

	parsed, err := Parse(buf.Bytes())
	assert.NoError(t, err)
	assert.NotNil(t, parsed)

	assert.Equal(t, "test.txt", parsed.Name)
	assert.Equal(t, "test-utf8.txt", parsed.NameUtf8)
	assert.Equal(t, int64(1024), parsed.Size)
	assert.Equal(t, int64(262144), parsed.PieceLength)
	assert.Equal(t, 1, parsed.PieceCount)
	assert.True(t, parsed.IsPrivate)
	assert.Equal(t, "test-source", parsed.Source)
	assert.Equal(t, "http://tracker.example.com", parsed.Announce)
	assert.Len(t, parsed.AnnounceList, 1)
	assert.Len(t, parsed.AnnounceList[0], 2)
	assert.Equal(t, int64(1600000000), parsed.CreationDate)
	assert.Equal(t, "test comment", parsed.Comment)
	assert.Equal(t, "test client", parsed.CreatedBy)
	assert.Equal(t, "UTF-8", parsed.Encoding)
	assert.Len(t, parsed.UrlList, 1)

	assert.Len(t, parsed.Files, 1)
	assert.Equal(t, "test.txt", parsed.Files[0].Path)
	assert.Equal(t, int64(1024), parsed.Files[0].Size)
	assert.NotEmpty(t, parsed.InfoHash)
}

func TestParseMultiFile(t *testing.T) {
	dummyTorrent := map[string]interface{}{
		"announce": "http://tracker.example.com",
		"info": map[string]interface{}{
			"name":         "test_dir",
			"piece length": int64(524288),
			"pieces":       "1234567890123456789012345678901234567890", // 40 bytes (2 pieces)
			"files": []interface{}{
				map[string]interface{}{
					"length": int64(1024),
					"path":   []interface{}{"file1.txt"},
				},
				map[string]interface{}{
					"length": int64(2048),
					"path":   []interface{}{"subdir", "file2.txt"},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := bencode.NewEncoder(&buf).Encode(dummyTorrent)
	assert.NoError(t, err)

	parsed, err := Parse(buf.Bytes())
	assert.NoError(t, err)
	assert.NotNil(t, parsed)

	assert.Equal(t, "test_dir", parsed.Name)
	assert.Equal(t, int64(3072), parsed.Size)
	assert.Equal(t, int64(524288), parsed.PieceLength)
	assert.Equal(t, 2, parsed.PieceCount)
	assert.False(t, parsed.IsPrivate)
	assert.Len(t, parsed.Files, 2)
	assert.Equal(t, "file1.txt", parsed.Files[0].Path)
	assert.Equal(t, int64(1024), parsed.Files[0].Size)
	// path should be joined correctly depending on OS, assuming Unix-like for output comparison
	assert.Contains(t, parsed.Files[1].Path, "file2.txt")
	assert.Equal(t, int64(2048), parsed.Files[1].Size)
}

func TestParseRealFile(t *testing.T) {
	data, err := os.ReadFile("debian-13.3.0-amd64-netinst.iso.torrent")
	if err != nil {
		t.Skip("debian-13.3.0-amd64-netinst.iso.torrent not found, skipping real file test")
	}

	parsed, err := Parse(data)
	assert.NoError(t, err)
	assert.NotNil(t, parsed)

	// Since we don't know the exact hash beforehand, just assert the structure looks sound over an actual file
	assert.NotEmpty(t, parsed.InfoHash)
	assert.Equal(t, int64(40), int64(len(parsed.InfoHash))) // hex encoded sha1 is 40 chars
	assert.True(t, parsed.Size > 0)
	assert.True(t, len(parsed.Files) > 0)
}

func TestMarshal(t *testing.T) {
	data, err := os.ReadFile("debian-13.3.0-amd64-netinst.iso.torrent")
	if err != nil {
		t.Skip("debian-13.3.0-amd64-netinst.iso.torrent not found, skipping marshal test")
	}

	parsed, err := Parse(data)
	assert.NoError(t, err)

	originalHash := parsed.InfoHash
	assert.NotEmpty(t, originalHash)

	// Modify announce struct as the tracker would
	parsed.Announce = "http://my-tracker.hermes.local:8080/announce"
	parsed.AnnounceList = [][]string{
		{"http://my-tracker.hermes.local:8080/announce"},
		{"http://backup-tracker.hermes.local:8080/announce"},
	}

	newTorrentBytes, err := parsed.Marshal()
	assert.NoError(t, err)
	assert.NotEmpty(t, newTorrentBytes)

	// Re-parse to verify
	reParsed, err := Parse(newTorrentBytes)
	assert.NoError(t, err)

	// InfoHash should exactly match original since we didn't touch the Info dictionary
	assert.Equal(t, originalHash, reParsed.InfoHash)

	// Updated fields should be reflected
	assert.Equal(t, "http://my-tracker.hermes.local:8080/announce", reParsed.Announce)
	assert.Equal(t, 2, len(reParsed.AnnounceList))
	assert.Equal(t, "http://backup-tracker.hermes.local:8080/announce", reParsed.AnnounceList[1][0])
}
