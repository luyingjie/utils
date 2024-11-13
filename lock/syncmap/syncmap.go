package syncmap

import (
	"sync"
)

// 线程安全的map，使用以下方法：
// Global.Store("greece", 97)  写入
// Global.Load("london")       读取
// Global.Delete("london")     删除
// 遍历
// scene.Range(func(k, v interface{}) bool {
// 	return true // true 继续迭代   false  停止
// })
// var Global sync.Map

// type k struct{
// 	sync.RWMutex or mx sync.RWMutex
// }
// k.Lock() or k.mx.Lock()

var mu sync.RWMutex

// 获取全局变量
func Get(m *map[string]interface{}, name string) interface{} {
	// m = new(sync.RWMutex)
	mu.RLock()
	defer mu.RUnlock()
	return (*m)[name]
}

// 设置全局变量
func Set(m *map[string]interface{}, name string, value interface{}) {
	mu.Lock()
	(*m)[name] = value
	mu.Unlock()
}
