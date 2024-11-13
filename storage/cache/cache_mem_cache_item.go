package cache

import "github.com/luyingjie/utils/container/vtime"

func (item *memCacheItem) IsExpired() bool {
	if item.e >= vtime.TimestampMilli() {
		return false
	}
	return true
}
