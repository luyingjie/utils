package cache

import vtime "github.com/luyingjie/utils/os/time"

func (item *memCacheItem) IsExpired() bool {
	if item.e >= vtime.TimestampMilli() {
		return false
	}
	return true
}
