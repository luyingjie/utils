package cache

import (
	"time"

	"github.com/luyingjie/utils/container/vvar"
)

var defaultCache = New()

func Set(key interface{}, value interface{}, duration time.Duration) {
	defaultCache.Set(key, value, duration)
}

func Update(key interface{}, value interface{}) (oldValue interface{}, exist bool) {
	return defaultCache.Update(key, value)
}

func SetIfNotExist(key interface{}, value interface{}, duration time.Duration) bool {
	return defaultCache.SetIfNotExist(key, value, duration)
}

func Sets(data map[interface{}]interface{}, duration time.Duration) {
	defaultCache.Sets(data, duration)
}

func Get(key interface{}) interface{} {
	return defaultCache.Get(key)
}

func GetVar(key interface{}) *vvar.Var {
	return defaultCache.GetVar(key)
}

func GetOrSet(key interface{}, value interface{}, duration time.Duration) interface{} {
	return defaultCache.GetOrSet(key, value, duration)
}

func GetOrSetFunc(key interface{}, f func() interface{}, duration time.Duration) interface{} {
	return defaultCache.GetOrSetFunc(key, f, duration)
}

func GetOrSetFuncLock(key interface{}, f func() interface{}, duration time.Duration) interface{} {
	return defaultCache.GetOrSetFuncLock(key, f, duration)
}

func Contains(key interface{}) bool {
	return defaultCache.Contains(key)
}

func Remove(keys ...interface{}) (value interface{}) {
	return defaultCache.Remove(keys...)
}

func Removes(keys []interface{}) {
	defaultCache.Removes(keys)
}

func Data() map[interface{}]interface{} {
	return defaultCache.Data()
}

func Keys() []interface{} {
	return defaultCache.Keys()
}

func KeyStrings() []string {
	return defaultCache.KeyStrings()
}

func Values() []interface{} {
	return defaultCache.Values()
}

func Size() int {
	return defaultCache.Size()
}

func GetExpire(key interface{}) time.Duration {
	return defaultCache.GetExpire(key)
}

func UpdateExpire(key interface{}, duration time.Duration) (oldDuration time.Duration) {
	return defaultCache.UpdateExpire(key, duration)
}
