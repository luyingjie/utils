package cache

import (
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/luyingjie/utils/os/timer"
)

type Cache struct {
	*memCache
}

func New(lruCap ...int) *Cache {
	c := &Cache{
		memCache: newMemCache(lruCap...),
	}
	timer.AddSingleton(time.Second, c.syncEventAndClearExpired)
	return c
}

func (c *Cache) Clear() {
	old := atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&c.memCache)), unsafe.Pointer(newMemCache()))
	(*memCache)(old).Close()
}
