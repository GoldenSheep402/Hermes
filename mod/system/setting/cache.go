package setting

import "github.com/patrickmn/go-cache"

var settingCache = cache.New(cache.NoExpiration, cache.NoExpiration)

func ClearCache() {
	settingCache.Flush()
}
