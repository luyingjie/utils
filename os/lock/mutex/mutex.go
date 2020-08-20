package mutex

// 线程安全的map，使用以下方法：
// Global.Store("greece", 97)  写入
// Global.Load("london")       读取
// Global.Delete("london")     删除
// 遍历
// scene.Range(func(k, v interface{}) bool {
// 	return true // true 继续迭代   false  停止
// })
// var Global sync.Map
