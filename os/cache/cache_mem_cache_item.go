package cache

import mytime "utils/os/time"

func (item *memCacheItem) IsExpired() bool {
	if item.e >= mytime.TimestampMilli() {
		return false
	}
	return true
}
