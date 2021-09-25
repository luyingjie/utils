package vcache

import (
	"sync"
	"time"
)

type CacheItem struct {
	sync.RWMutex

	// 缓存 Key
	key interface{}
	// 缓存 值
	data interface{}
	// 不被访问时的存活时间
	lifeSpan time.Duration

	// 创建时间
	createdOn time.Time
	// 最后一次访问时间
	accessedOn time.Time
	// 访问次数
	accessCount int64

	// 缓存被清理的时候的回调
	aboutToExpire []func(key interface{})
}

func NewCacheItem(key interface{}, lifeSpan time.Duration, data interface{}) *CacheItem {
	t := time.Now()
	return &CacheItem{
		key:           key,
		lifeSpan:      lifeSpan,
		createdOn:     t,
		accessedOn:    t,
		accessCount:   0,
		aboutToExpire: nil,
		data:          data,
	}
}

// KeepAlive 标记一个条目将在另一个过期时间内保留。
func (item *CacheItem) KeepAlive() {
	item.Lock()
	defer item.Unlock()
	item.accessedOn = time.Now()
	item.accessCount++
}

// LifeSpan 返回该项的过期时间。
func (item *CacheItem) LifeSpan() time.Duration {
	// immutable
	return item.lifeSpan
}

// AccessedOn 返回最后一次访问此项的时间。
func (item *CacheItem) AccessedOn() time.Time {
	item.RLock()
	defer item.RUnlock()
	return item.accessedOn
}

// CreatedOn 返回将此项添加到缓存的时间。
func (item *CacheItem) CreatedOn() time.Time {
	// immutable
	return item.createdOn
}

// AccessCount 返回该项被访问的频率。
func (item *CacheItem) AccessCount() int64 {
	item.RLock()
	defer item.RUnlock()
	return item.accessCount
}

// Key 返回此缓存项的键。
func (item *CacheItem) Key() interface{} {
	// immutable
	return item.key
}

// Data 返回此缓存项的值。
func (item *CacheItem) Data() interface{} {
	// immutable
	return item.data
}

// SetAboutToExpireCallback 配置一个回调函数， 在项目即将从缓存中删除之前。
func (item *CacheItem) SetAboutToExpireCallback(f func(interface{})) {
	if len(item.aboutToExpire) > 0 {
		item.RemoveAboutToExpireCallback()
	}
	item.Lock()
	defer item.Unlock()
	item.aboutToExpire = append(item.aboutToExpire, f)
}

// AddAboutToExpireCallback 添加一个新的回调到AboutToExpire队列
func (item *CacheItem) AddAboutToExpireCallback(f func(interface{})) {
	item.Lock()
	defer item.Unlock()
	item.aboutToExpire = append(item.aboutToExpire, f)
}

// RemoveAboutToExpireCallback 清空即将过期的回调队列
func (item *CacheItem) RemoveAboutToExpireCallback() {
	item.Lock()
	defer item.Unlock()
	item.aboutToExpire = nil
}
