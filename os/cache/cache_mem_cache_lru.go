package cache

import (
	"time"

	"github.com/luyingjie/utils/container/list"
	vmap "github.com/luyingjie/utils/container/map"
	vtype "github.com/luyingjie/utils/container/type"

	"github.com/luyingjie/utils/os/timer"
)

type memCacheLru struct {
	cache   *memCache
	data    *vmap.Map
	list    *list.List
	rawList *list.List
	closed  *vtype.Bool
}

func newMemCacheLru(cache *memCache) *memCacheLru {
	lru := &memCacheLru{
		cache:   cache,
		data:    vmap.New(true),
		list:    list.New(true),
		rawList: list.New(true),
		closed:  vtype.NewBool(),
	}
	timer.AddSingleton(time.Second, lru.SyncAndClear)
	return lru
}

func (lru *memCacheLru) Close() {
	lru.closed.Set(true)
}

func (lru *memCacheLru) Remove(key interface{}) {
	if v := lru.data.Get(key); v != nil {
		lru.data.Remove(key)
		lru.list.Remove(v.(*list.Element))
	}
}

func (lru *memCacheLru) Size() int {
	return lru.data.Size()
}

func (lru *memCacheLru) Push(key interface{}) {
	lru.rawList.PushBack(key)
}

func (lru *memCacheLru) Pop() interface{} {
	if v := lru.list.PopBack(); v != nil {
		lru.data.Remove(v)
		return v
	}
	return nil
}

func (lru *memCacheLru) SyncAndClear() {
	if lru.closed.Val() {
		timer.Exit()
		return
	}
	for {
		if v := lru.rawList.PopFront(); v != nil {
			if v := lru.data.Get(v); v != nil {
				lru.list.Remove(v.(*list.Element))
			}
			lru.data.Set(v, lru.list.PushFront(v))
		} else {
			break
		}
	}
	for i := lru.Size() - lru.cache.cap; i > 0; i-- {
		if s := lru.Pop(); s != nil {
			lru.cache.clearByKey(s, true)
		}
	}
}
