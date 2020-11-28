package cache

import vtime "utils/os/time"

func (item *memCacheItem) IsExpired() bool {
	if item.e >= vtime.TimestampMilli() {
		return false
	}
	return true
}
