package dao

import "strings"

// redisKeyPrefix isolates Redis keys when multiple logical sites share one Redis instance.
// Set via SetRedisKeyPrefix (e.g. "prod" -> "prod:tracker:..."). Empty means no prefix.
var redisKeyPrefix string

// SetRedisKeyPrefix configures a namespace prefix for all tracker Redis keys.
// Non-empty values are normalized to end with ":".
func SetRedisKeyPrefix(p string) {
	p = strings.TrimSpace(p)
	if p != "" && !strings.HasSuffix(p, ":") {
		p += ":"
	}
	redisKeyPrefix = p
}

func prefixed(s string) string {
	return redisKeyPrefix + s
}
