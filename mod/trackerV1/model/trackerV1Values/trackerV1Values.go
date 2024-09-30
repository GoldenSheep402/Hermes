package trackerV1Values

import "time"

const (
	OK     = "ok"
	Banned = "Banned"
)

const (
	TTL = 10 * time.Minute
)

const (
	Downloading = iota
	Seeding
	Finished
	Stopped
	ReadySeeding
	Completed
)

const (
	UserDownload_Download = iota
	UserDownload_Upload
	UserDownload_Finished
	UserDownload_Seeding
)
