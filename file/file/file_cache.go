package file

import (
	"time"

	"github.com/luyingjie/utils/util/cmdenv"

	"github.com/luyingjie/utils/storage/cache"

	"github.com/luyingjie/utils/file/fsnotify"
)

const (
	gDEFAULT_CACHE_EXPIRE = time.Minute
)

var (
	cacheExpire = cmdenv.Get("file.cache", gDEFAULT_CACHE_EXPIRE).Duration()
)

func GetContentsWithCache(path string, duration ...time.Duration) string {
	return string(GetBytesWithCache(path, duration...))
}

func GetBytesWithCache(path string, duration ...time.Duration) []byte {
	key := cacheKey(path)
	expire := cacheExpire
	if len(duration) > 0 {
		expire = duration[0]
	}
	r := cache.GetOrSetFuncLock(key, func() interface{} {
		b := GetBytes(path)
		if b != nil {
			_, _ = fsnotify.Add(path, func(event *fsnotify.Event) {
				cache.Remove(key)
				fsnotify.Exit()
			})
		}
		return b
	}, expire)
	if r != nil {
		return r.([]byte)
	}
	return nil
}

func cacheKey(path string) string {
	return "file.cache:" + path
}
