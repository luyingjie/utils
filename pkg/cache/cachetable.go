package cache

import (
	"log"
	"sort"
	"sync"
	"time"

	verror "github.com/luyingjie/utils/os/error"
)

type CacheTable struct {
	sync.RWMutex

	// 表名称
	name string
	// 缓存的 Item.
	items map[interface{}]*CacheItem

	// 清理定时器
	cleanupTimer *time.Timer
	// 清理的时长。
	cleanupInterval time.Duration

	// 日志器
	logger *log.Logger

	// 当试图加载不存在的键时触发回调方法。
	loadData func(key interface{}, args ...interface{}) *CacheItem
	// 在向缓存添加新项时触发的回调方法。
	addedItem []func(item *CacheItem)
	// 在从缓存中删除项之前触发的回调方法。
	aboutToDeleteItem []func(item *CacheItem)
}

// Count 返回当前缓存中存储的项数。
func (table *CacheTable) Count() int {
	table.RLock()
	defer table.RUnlock()
	return len(table.items)
}

// Foreach 遍历所有Item
func (table *CacheTable) Foreach(trans func(key interface{}, item *CacheItem)) {
	table.RLock()
	defer table.RUnlock()

	for k, v := range table.items {
		trans(k, v)
	}
}

// SetDataLoader 当尝试访问一个不存在的键的回调函数。
func (table *CacheTable) SetDataLoader(f func(interface{}, ...interface{}) *CacheItem) {
	table.Lock()
	defer table.Unlock()
	table.loadData = f
}

// setadditemcallback 设置新添加Item时候的回调。。
func (table *CacheTable) SetAddedItemCallback(f func(*CacheItem)) {
	if len(table.addedItem) > 0 {
		table.RemoveAddedItemCallbacks()
	}
	table.Lock()
	defer table.Unlock()
	table.addedItem = append(table.addedItem, f)
}

// AddAddedItemCallback 添加一个新的回调到addedItem队列
func (table *CacheTable) AddAddedItemCallback(f func(*CacheItem)) {
	table.Lock()
	defer table.Unlock()
	table.addedItem = append(table.addedItem, f)
}

// RemoveAddedItemCallbacks 清空添加的项目回调队列
func (table *CacheTable) RemoveAddedItemCallbacks() {
	table.Lock()
	defer table.Unlock()
	table.addedItem = nil
}

// SetAboutToDeleteItemCallback 配置删除缓存的回调函数。
func (table *CacheTable) SetAboutToDeleteItemCallback(f func(*CacheItem)) {
	if len(table.aboutToDeleteItem) > 0 {
		table.RemoveAboutToDeleteItemCallback()
	}
	table.Lock()
	defer table.Unlock()
	table.aboutToDeleteItem = append(table.aboutToDeleteItem, f)
}

// AddAboutToDeleteItemCallback 添加一个新的回调到AboutToDeleteItem队列
func (table *CacheTable) AddAboutToDeleteItemCallback(f func(*CacheItem)) {
	table.Lock()
	defer table.Unlock()
	table.aboutToDeleteItem = append(table.aboutToDeleteItem, f)
}

// RemoveAboutToDeleteItemCallback 清空即将删除的项目回调队列
func (table *CacheTable) RemoveAboutToDeleteItemCallback() {
	table.Lock()
	defer table.Unlock()
	table.aboutToDeleteItem = nil
}

// SetLogger 写入日志
func (table *CacheTable) SetLogger(logger *log.Logger) {
	table.Lock()
	defer table.Unlock()
	table.logger = logger
}

// Expiration 检查循环，由自调整定时器触发。
func (table *CacheTable) expirationCheck() {
	table.Lock()
	if table.cleanupTimer != nil {
		table.cleanupTimer.Stop()
	}
	if table.cleanupInterval > 0 {
		table.log("Expiration check triggered after", table.cleanupInterval, "for table", table.name)
	} else {
		table.log("Expiration check installed for table", table.name)
	}

	now := time.Now()
	smallestDuration := 0 * time.Second
	for key, item := range table.items {
		item.RLock()
		lifeSpan := item.lifeSpan
		accessedOn := item.accessedOn
		item.RUnlock()

		if lifeSpan == 0 {
			continue
		}
		if now.Sub(accessedOn) >= lifeSpan {
			table.deleteInternal(key)
		} else {
			if smallestDuration == 0 || lifeSpan-now.Sub(accessedOn) < smallestDuration {
				smallestDuration = lifeSpan - now.Sub(accessedOn)
			}
		}
	}

	table.cleanupInterval = smallestDuration
	if smallestDuration > 0 {
		table.cleanupTimer = time.AfterFunc(smallestDuration, func() {
			go table.expirationCheck()
		})
	}
	table.Unlock()
}

func (table *CacheTable) addInternal(item *CacheItem) {
	table.log("Adding item with key", item.key, "and lifespan of", item.lifeSpan, "to table", table.name)
	table.items[item.key] = item

	expDur := table.cleanupInterval
	addedItem := table.addedItem
	table.Unlock()

	if addedItem != nil {
		for _, callback := range addedItem {
			callback(item)
		}
	}

	if item.lifeSpan > 0 && (expDur == 0 || item.lifeSpan < expDur) {
		table.expirationCheck()
	}
}

func (table *CacheTable) Add(key interface{}, lifeSpan time.Duration, data interface{}) *CacheItem {
	item := NewCacheItem(key, lifeSpan, data)

	table.Lock()
	table.addInternal(item)

	return item
}

func (table *CacheTable) deleteInternal(key interface{}) (*CacheItem, error) {
	r, ok := table.items[key]
	if !ok {
		return nil, verror.New("ErrKeyNotFound")
	}

	aboutToDeleteItem := table.aboutToDeleteItem
	table.Unlock()

	if aboutToDeleteItem != nil {
		for _, callback := range aboutToDeleteItem {
			callback(r)
		}
	}

	r.RLock()
	defer r.RUnlock()
	if r.aboutToExpire != nil {
		for _, callback := range r.aboutToExpire {
			callback(key)
		}
	}

	table.Lock()
	table.log("Deleting item with key", key, "created on", r.createdOn, "and hit", r.accessCount, "times from table", table.name)
	delete(table.items, key)

	return r, nil
}

func (table *CacheTable) Delete(key interface{}) (*CacheItem, error) {
	table.Lock()
	defer table.Unlock()

	return table.deleteInternal(key)
}

func (table *CacheTable) Exists(key interface{}) bool {
	table.RLock()
	defer table.RUnlock()
	_, ok := table.items[key]

	return ok
}

func (table *CacheTable) NotFoundAdd(key interface{}, lifeSpan time.Duration, data interface{}) bool {
	table.Lock()

	if _, ok := table.items[key]; ok {
		table.Unlock()
		return false
	}

	item := NewCacheItem(key, lifeSpan, data)
	table.addInternal(item)

	return true
}

func (table *CacheTable) Value(key interface{}, args ...interface{}) (*CacheItem, error) {
	table.RLock()
	r, ok := table.items[key]
	loadData := table.loadData
	table.RUnlock()

	if ok {
		r.KeepAlive()
		return r, nil
	}

	if loadData != nil {
		item := loadData(key, args...)
		if item != nil {
			table.Add(key, item.lifeSpan, item.data)
			return item, nil
		}

		return nil, verror.New("ErrKeyNotFoundOrLoadable")
	}

	return nil, verror.New("ErrKeyNotFound")
}

func (table *CacheTable) Flush() {
	table.Lock()
	defer table.Unlock()

	table.log("Flushing table", table.name)

	table.items = make(map[interface{}]*CacheItem)
	table.cleanupInterval = 0
	if table.cleanupTimer != nil {
		table.cleanupTimer.Stop()
	}
}

type CacheItemPair struct {
	Key         interface{}
	AccessCount int64
}

type CacheItemPairList []CacheItemPair

func (p CacheItemPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p CacheItemPairList) Len() int           { return len(p) }
func (p CacheItemPairList) Less(i, j int) bool { return p[i].AccessCount > p[j].AccessCount }

func (table *CacheTable) MostAccessed(count int64) []*CacheItem {
	table.RLock()
	defer table.RUnlock()

	p := make(CacheItemPairList, len(table.items))
	i := 0
	for k, v := range table.items {
		p[i] = CacheItemPair{k, v.accessCount}
		i++
	}
	sort.Sort(p)

	var r []*CacheItem
	c := int64(0)
	for _, v := range p {
		if c >= count {
			break
		}

		item, ok := table.items[v.Key]
		if ok {
			r = append(r, item)
		}
		c++
	}

	return r
}

func (table *CacheTable) log(v ...interface{}) {
	if table.logger == nil {
		return
	}

	table.logger.Println(v...)
}
