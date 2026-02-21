package torrent

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/anacrolix/torrent/metainfo"
)

// ParsedFile represents a file inside a torrent
type ParsedFile struct {
	Path string
	Size int64
}

// ParsedTorrent represents the extracted metadata from a .torrent file
type ParsedTorrent struct {
	// rawMeta stores the original parsed MetaInfo to ensure InfoHash consistency
	// when re-marshalling, while allowing modification of root-level fields like Announce.
	rawMeta *metainfo.MetaInfo

	// The unique 40-character hex string identifying this torrent
	InfoHash string
	// The suggested name for the saved file or directory
	Name string
	// The UTF-8 encoded version of the name
	NameUtf8 string
	// Total size of all files in bytes
	Size int64
	// The number of bytes in each piece (typically 256KB, 512KB, etc.)
	PieceLength int64
	// Total number of pieces the torrent is split into
	PieceCount int
	// If true, the client must not use DHT or PEX to find peers
	IsPrivate bool
	// Often used to indicate the site or group that released the torrent
	Source string
	// Metadata for each individual file contained in the torrent
	Files []ParsedFile

	// --- Tracker Information ---
	// The primary tracker URL (Announce URL)
	Announce string
	// A list of backup trackers organized by priority tiers (BEP 12)
	AnnounceList [][]string

	// Unix timestamp of when the torrent was created
	CreationDate int64
	// Free-form text comment by the author
	Comment string
	// The name and version of the program used to create the .torrent file
	CreatedBy string
	// Character encoding used for the 'pieces' part (usually UTF-8)
	Encoding string
	// List of web seeds (URLs for direct HTTP/FTP downloads)
	UrlList []string
}

// Parse parses a .torrent file from raw bytes.
func Parse(data []byte) (*ParsedTorrent, error) {
	mi, err := metainfo.Load(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to load metainfo: %w", err)
	}

	info, err := mi.UnmarshalInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal info dict: %w", err)
	}

	// Calculate InfoHash
	if len(mi.InfoBytes) == 0 {
		return nil, fmt.Errorf("missing info bytes")
	}
	hash := sha1.Sum(mi.InfoBytes)
	infoHashHex := hex.EncodeToString(hash[:])

	pt := &ParsedTorrent{
		rawMeta:      mi, // store original to allow perfect re-encoding without modifying info hash
		InfoHash:     infoHashHex,
		Name:         info.Name,
		NameUtf8:     info.NameUtf8,
		Size:         info.TotalLength(),
		PieceLength:  info.PieceLength,
		PieceCount:   info.NumPieces(),
		IsPrivate:    info.Private != nil && *info.Private,
		Source:       info.Source,
		Announce:     mi.Announce,
		AnnounceList: mi.AnnounceList,
		CreationDate: mi.CreationDate,
		Comment:      mi.Comment,
		CreatedBy:    mi.CreatedBy,
		Encoding:     mi.Encoding,
		UrlList:      mi.UrlList,
	}

	if len(info.Files) > 0 {
		// Multi-file torrent
		for _, f := range info.Files {
			// f.Path is []string, join them
			filePath := filepath.Join(f.Path...)
			pt.Files = append(pt.Files, ParsedFile{
				Path: filePath,
				Size: f.Length,
			})
		}
	} else {
		// Single-file torrent
		pt.Files = append(pt.Files, ParsedFile{
			Path: info.Name,
			Size: info.Length,
		})
	}

	return pt, nil
}

// Marshal re-encodes the ParsedTorrent back into a .torrent file byte array.
// It uses the originally parsed raw string to ensure the InfoHash doesn't change,
// but it applies any modifications made to Announce, Comment, etc.
func (pt *ParsedTorrent) Marshal() ([]byte, error) {
	if pt.rawMeta == nil {
		return nil, fmt.Errorf("cannot marshal ParsedTorrent not created via Parse")
	}

	// Update the updatable fields on the raw metaclass
	pt.rawMeta.Announce = pt.Announce
	pt.rawMeta.AnnounceList = pt.AnnounceList
	pt.rawMeta.Comment = pt.Comment
	pt.rawMeta.CreatedBy = pt.CreatedBy
	pt.rawMeta.CreationDate = pt.CreationDate
	pt.rawMeta.Encoding = pt.Encoding
	pt.rawMeta.UrlList = pt.UrlList

	var buf bytes.Buffer
	err := pt.rawMeta.Write(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to write torrent metaclass: %w", err)
	}

	return buf.Bytes(), nil
}
