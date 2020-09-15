package cache

import (
	"math"
	"sync"
	"time"

	myvar "utils/container/var"

	mytime "utils/os/time"

	"utils/os/timer"

	"utils/container/list"

	"utils/container/set"

	mytype "utils/container/type"

	"utils/convert/conv"
)

type memCache struct {
	dataMu sync.RWMutex

	expireTimeMu sync.RWMutex

	expireSetMu sync.RWMutex

	cap int

	data map[interface{}]memCacheItem

	expireTimes map[interface{}]int64

	expireSets map[int64]*set.Set

	lru *memCacheLru

	lruGetList *list.List

	eventList *list.List

	closed *mytype.Bool
}

type memCacheItem struct {
	v interface{}
	e int64
}

type memCacheEvent struct {
	k interface{}
	e int64
}

const (
	gDEFAULT_MAX_EXPIRE = 9223372036854
)

func newMemCache(lruCap ...int) *memCache {
	c := &memCache{
		lruGetList:  list.New(true),
		data:        make(map[interface{}]memCacheItem),
		expireTimes: make(map[interface{}]int64),
		expireSets:  make(map[int64]*set.Set),
		eventList:   list.New(true),
		closed:      mytype.NewBool(),
	}
	if len(lruCap) > 0 {
		c.cap = lruCap[0]
		c.lru = newMemCacheLru(c)
	}
	return c
}

func (c *memCache) Set(key interface{}, value interface{}, duration time.Duration) {
	expireTime := c.getInternalExpire(duration)
	c.dataMu.Lock()
	c.data[key] = memCacheItem{
		v: value,
		e: expireTime,
	}
	c.dataMu.Unlock()
	c.eventList.PushBack(&memCacheEvent{
		k: key,
		e: expireTime,
	})
}

func (c *memCache) Update(key interface{}, value interface{}) (oldValue interface{}, exist bool) {
	c.dataMu.Lock()
	defer c.dataMu.Unlock()
	if item, ok := c.data[key]; ok {
		c.data[key] = memCacheItem{
			v: value,
			e: item.e,
		}
		return item.v, true
	}
	return nil, false
}

func (c *memCache) UpdateExpire(key interface{}, duration time.Duration) (oldDuration time.Duration) {
	newExpireTime := c.getInternalExpire(duration)
	c.dataMu.Lock()
	defer c.dataMu.Unlock()
	if item, ok := c.data[key]; ok {
		c.data[key] = memCacheItem{
			v: item.v,
			e: newExpireTime,
		}
		c.eventList.PushBack(&memCacheEvent{
			k: key,
			e: newExpireTime,
		})
		return time.Duration(item.e-mytime.TimestampMilli()) * time.Millisecond
	}
	return -1
}

func (c *memCache) GetExpire(key interface{}) time.Duration {
	c.dataMu.RLock()
	defer c.dataMu.RUnlock()
	if item, ok := c.data[key]; ok {
		return time.Duration(item.e-mytime.TimestampMilli()) * time.Millisecond
	}
	return -1
}

func (c *memCache) doSetWithLockCheck(key interface{}, value interface{}, duration time.Duration) interface{} {
	expireTimestamp := c.getInternalExpire(duration)
	c.dataMu.Lock()
	defer c.dataMu.Unlock()
	if v, ok := c.data[key]; ok && !v.IsExpired() {
		return v.v
	}
	if f, ok := value.(func() interface{}); ok {
		value = f()
		if value == nil {
			return nil
		}
	}
	c.data[key] = memCacheItem{v: value, e: expireTimestamp}
	c.eventList.PushBack(&memCacheEvent{k: key, e: expireTimestamp})
	return value
}

func (c *memCache) getInternalExpire(duration time.Duration) int64 {
	if duration == 0 {
		return gDEFAULT_MAX_EXPIRE
	} else {
		return mytime.TimestampMilli() + duration.Nanoseconds()/1000000
	}
}

func (c *memCache) makeExpireKey(expire int64) int64 {
	return int64(math.Ceil(float64(expire/1000)+1) * 1000)
}

func (c *memCache) getExpireSet(expire int64) (expireSet *set.Set) {
	c.expireSetMu.RLock()
	expireSet, _ = c.expireSets[expire]
	c.expireSetMu.RUnlock()
	return
}

func (c *memCache) getOrNewExpireSet(expire int64) (expireSet *set.Set) {
	if expireSet = c.getExpireSet(expire); expireSet == nil {
		expireSet = set.New(true)
		c.expireSetMu.Lock()
		if es, ok := c.expireSets[expire]; ok {
			expireSet = es
		} else {
			c.expireSets[expire] = expireSet
		}
		c.expireSetMu.Unlock()
	}
	return
}

func (c *memCache) SetIfNotExist(key interface{}, value interface{}, duration time.Duration) bool {
	if !c.Contains(key) {
		c.doSetWithLockCheck(key, value, duration)
		return true
	}
	return false
}

func (c *memCache) Sets(data map[interface{}]interface{}, duration time.Duration) {
	expireTime := c.getInternalExpire(duration)
	for k, v := range data {
		c.dataMu.Lock()
		c.data[k] = memCacheItem{
			v: v,
			e: expireTime,
		}
		c.dataMu.Unlock()
		c.eventList.PushBack(&memCacheEvent{
			k: k,
			e: expireTime,
		})
	}
}

func (c *memCache) Get(key interface{}) interface{} {
	c.dataMu.RLock()
	item, ok := c.data[key]
	c.dataMu.RUnlock()
	if ok && !item.IsExpired() {
		if c.cap > 0 {
			c.lruGetList.PushBack(key)
		}
		return item.v
	}
	return nil
}

func (c *memCache) GetVar(key interface{}) *myvar.Var {
	return myvar.New(c.Get(key))
}

func (c *memCache) GetOrSet(key interface{}, value interface{}, duration time.Duration) interface{} {
	if v := c.Get(key); v == nil {
		return c.doSetWithLockCheck(key, value, duration)
	} else {
		return v
	}
}

func (c *memCache) GetOrSetFunc(key interface{}, f func() interface{}, duration time.Duration) interface{} {
	if v := c.Get(key); v == nil {
		value := f()
		if value == nil {
			return nil
		}
		return c.doSetWithLockCheck(key, value, duration)
	} else {
		return v
	}
}

func (c *memCache) GetOrSetFuncLock(key interface{}, f func() interface{}, duration time.Duration) interface{} {
	if v := c.Get(key); v == nil {
		return c.doSetWithLockCheck(key, f, duration)
	} else {
		return v
	}
}

func (c *memCache) Contains(key interface{}) bool {
	return c.Get(key) != nil
}

func (c *memCache) Remove(keys ...interface{}) (value interface{}) {
	c.dataMu.Lock()
	defer c.dataMu.Unlock()
	for _, key := range keys {
		item, ok := c.data[key]
		if ok {
			value = item.v
			delete(c.data, key)
			c.eventList.PushBack(&memCacheEvent{
				k: key,
				e: mytime.TimestampMilli() - 1000,
			})
		}
	}
	return
}

func (c *memCache) Removes(keys []interface{}) {
	c.Remove(keys...)
}

func (c *memCache) Data() map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	c.dataMu.RLock()
	for k, v := range c.data {
		if !v.IsExpired() {
			m[k] = v.v
		}
	}
	c.dataMu.RUnlock()
	return m
}

func (c *memCache) Keys() []interface{} {
	keys := make([]interface{}, 0)
	c.dataMu.RLock()
	for k, v := range c.data {
		if !v.IsExpired() {
			keys = append(keys, k)
		}
	}
	c.dataMu.RUnlock()
	return keys
}

func (c *memCache) KeyStrings() []string {
	return conv.Strings(c.Keys())
}

func (c *memCache) Values() []interface{} {
	values := make([]interface{}, 0)
	c.dataMu.RLock()
	for _, v := range c.data {
		if !v.IsExpired() {
			values = append(values, v.v)
		}
	}
	c.dataMu.RUnlock()
	return values
}

func (c *memCache) Size() (size int) {
	c.dataMu.RLock()
	size = len(c.data)
	c.dataMu.RUnlock()
	return
}

func (c *memCache) Close() {
	if c.cap > 0 {
		c.lru.Close()
	}
	c.closed.Set(true)
}

func (c *memCache) syncEventAndClearExpired() {
	if c.closed.Val() {
		timer.Exit()
		return
	}
	var (
		event         *memCacheEvent
		oldExpireTime int64
		newExpireTime int64
	)

	for {
		v := c.eventList.PopFront()
		if v == nil {
			break
		}
		event = v.(*memCacheEvent)
		c.expireTimeMu.RLock()
		oldExpireTime = c.expireTimes[event.k]
		c.expireTimeMu.RUnlock()
		newExpireTime = c.makeExpireKey(event.e)
		if newExpireTime != oldExpireTime {
			c.getOrNewExpireSet(newExpireTime).Add(event.k)
			if oldExpireTime != 0 {
				c.getOrNewExpireSet(oldExpireTime).Remove(event.k)
			}
			c.expireTimeMu.Lock()
			c.expireTimes[event.k] = newExpireTime
			c.expireTimeMu.Unlock()
		}
		if c.cap > 0 {
			c.lru.Push(event.k)
		}
	}
	if c.cap > 0 && c.lruGetList.Len() > 0 {
		for {
			if v := c.lruGetList.PopFront(); v != nil {
				c.lru.Push(v)
			} else {
				break
			}
		}
	}
	var (
		expireSet *set.Set
		ek        = c.makeExpireKey(mytime.TimestampMilli())
		eks       = []int64{ek - 1000, ek - 2000, ek - 3000, ek - 4000, ek - 5000}
	)
	for _, expireTime := range eks {
		if expireSet = c.getExpireSet(expireTime); expireSet != nil {
			expireSet.Iterator(func(key interface{}) bool {
				c.clearByKey(key)
				return true
			})
			c.expireSetMu.Lock()
			delete(c.expireSets, expireTime)
			c.expireSetMu.Unlock()
		}
	}
}

func (c *memCache) clearByKey(key interface{}, force ...bool) {
	c.dataMu.Lock()
	if item, ok := c.data[key]; (ok && item.IsExpired()) || (len(force) > 0 && force[0]) {
		delete(c.data, key)
	}
	c.dataMu.Unlock()

	c.expireTimeMu.Lock()
	delete(c.expireTimes, key)
	c.expireTimeMu.Unlock()

	if c.cap > 0 {
		c.lru.Remove(key)
	}
}
