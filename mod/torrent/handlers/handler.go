package handlers

import (
	"github.com/juanjiTech/jin"
)

// TODO: Reimplement torrent download/upload HTTP handlers with new Torrent model.
// Old code referenced: CreatedBy, Comment, FileInfo, Length, Md5sum, Pieces, Private, Source, NameUTF8
// and old user DAO methods (CheckKey).

func Registry(jinE *jin.Engine) {
	torrentGroup := jinE.Group("/api/torrent")

	torrentGroup.Group("/download").
		GET("/:id", DownloadTorrent)
}

func DownloadTorrent(c *jin.Context) {
	// TODO: implement with new model
	c.Writer.WriteHeader(501)
	c.Writer.Write([]byte("Not implemented"))
}
