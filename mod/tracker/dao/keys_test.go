package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetRedisKeyPrefix(t *testing.T) {
	defer SetRedisKeyPrefix("")

	SetRedisKeyPrefix("")
	assert.Equal(t, "tracker:peer:v2:t:p", peerKey("t", "p"))

	SetRedisKeyPrefix("siteA")
	assert.Equal(t, "siteA:tracker:peer:v2:t:p", peerKey("t", "p"))

	SetRedisKeyPrefix("siteB:")
	assert.Equal(t, "siteB:tracker:traffic:dirty:users", trafficDirtyUsersKey())
}
