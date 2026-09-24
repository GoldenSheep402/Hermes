package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedactPasskeyInPath(t *testing.T) {
	assert.Equal(t, "/announce/[passkey]", RedactPasskeyInPath("/announce/secretkey123"))
	assert.Equal(t, "/api/announce/[passkey]", RedactPasskeyInPath("/api/announce/secretkey123"))
	assert.Equal(t, "/scrape/[passkey]", RedactPasskeyInPath("/scrape/abc"))
	assert.Equal(t, "/foo/bar", RedactPasskeyInPath("/foo/bar"))
}
