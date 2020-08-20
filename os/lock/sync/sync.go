package sync

import (
	"sync"
)

var mu sync.RWMutex

// 获取全局变量
func GetGlobal(global *map[string]interface{}, name string) interface{} {
	// m = new(sync.RWMutex)
	mu.RLock()
	defer mu.RUnlock()
	return (*global)[name]
}

// 设置全局变量
func SetGlobal(global *map[string]interface{}, name string, value interface{}) {
	mu.Lock()
	(*global)[name] = value
	mu.Unlock()
}
