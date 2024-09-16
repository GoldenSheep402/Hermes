package example

import "embed"

// FS contains all files in this module,
// and since embed does not support import relative dir,
// we move embed from "modGenerator" to "example"
//
//go:embed *
var FS embed.FS
