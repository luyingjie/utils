package vmap

import (
	"utils/utils/json"

	"utils/utils/empty"

	"utils/convert/conv"

	"utils/container/list"
	"utils/container/var"
	"utils/utils/rwmutex"
)

type ListMap struct {
	mu   rwmutex.RWMutex
	data map[interface{}]*list.Element
	list *list.List
}

type ListMapNode struct {
	key   interface{}
	value interface{}
}

func NewListMap(safe ...bool) *ListMap {
	return &ListMap{
		mu:   rwmutex.Create(safe...),
		data: make(map[interface{}]*list.Element),
		list: list.New(),
	}
}

func NewListMapFrom(data map[interface{}]interface{}, safe ...bool) *ListMap {
	m := NewListMap(safe...)
	m.Sets(data)
	return m
}

func (m *ListMap) Iterator(f func(key, value interface{}) bool) {
	m.IteratorAsc(f)
}

func (m *ListMap) IteratorAsc(f func(key interface{}, value interface{}) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.list != nil {
		node := (*ListMapNode)(nil)
		m.list.IteratorAsc(func(e *list.Element) bool {
			node = e.Value.(*ListMapNode)
			return f(node.key, node.value)
		})
	}
}

func (m *ListMap) IteratorDesc(f func(key interface{}, value interface{}) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.list != nil {
		node := (*ListMapNode)(nil)
		m.list.IteratorDesc(func(e *list.Element) bool {
			node = e.Value.(*ListMapNode)
			return f(node.key, node.value)
		})
	}
}

func (m *ListMap) Clone(safe ...bool) *ListMap {
	return NewListMapFrom(m.Map(), safe...)
}

func (m *ListMap) Clear() {
	m.mu.Lock()
	m.data = make(map[interface{}]*list.Element)
	m.list = list.New()
	m.mu.Unlock()
}

func (m *ListMap) Replace(data map[interface{}]interface{}) {
	m.mu.Lock()
	m.data = make(map[interface{}]*list.Element)
	m.list = list.New()
	for key, value := range data {
		if e, ok := m.data[key]; !ok {
			m.data[key] = m.list.PushBack(&ListMapNode{key, value})
		} else {
			e.Value = &ListMapNode{key, value}
		}
	}
	m.mu.Unlock()
}

func (m *ListMap) Map() map[interface{}]interface{} {
	m.mu.RLock()
	var node *ListMapNode
	var data map[interface{}]interface{}
	if m.list != nil {
		data = make(map[interface{}]interface{}, len(m.data))
		m.list.IteratorAsc(func(e *list.Element) bool {
			node = e.Value.(*ListMapNode)
			data[node.key] = node.value
			return true
		})
	}
	m.mu.RUnlock()
	return data
}

func (m *ListMap) MapStrAny() map[string]interface{} {
	m.mu.RLock()
	var node *ListMapNode
	var data map[string]interface{}
	if m.list != nil {
		data = make(map[string]interface{}, len(m.data))
		m.list.IteratorAsc(func(e *list.Element) bool {
			node = e.Value.(*ListMapNode)
			data[conv.String(node.key)] = node.value
			return true
		})
	}
	m.mu.RUnlock()
	return data
}

func (m *ListMap) FilterEmpty() {
	m.mu.Lock()
	if m.list != nil {
		keys := make([]interface{}, 0)
		node := (*ListMapNode)(nil)
		m.list.IteratorAsc(func(e *list.Element) bool {
			node = e.Value.(*ListMapNode)
			if empty.IsEmpty(node.value) {
				keys = append(keys, node.key)
			}
			return true
		})
		if len(keys) > 0 {
			for _, key := range keys {
				if e, ok := m.data[key]; ok {
					delete(m.data, key)
					m.list.Remove(e)
				}
			}
		}
	}
	m.mu.Unlock()
}

func (m *ListMap) Set(key interface{}, value interface{}) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[interface{}]*list.Element)
		m.list = list.New()
	}
	if e, ok := m.data[key]; !ok {
		m.data[key] = m.list.PushBack(&ListMapNode{key, value})
	} else {
		e.Value = &ListMapNode{key, value}
	}
	m.mu.Unlock()
}

func (m *ListMap) Sets(data map[interface{}]interface{}) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[interface{}]*list.Element)
		m.list = list.New()
	}
	for key, value := range data {
		if e, ok := m.data[key]; !ok {
			m.data[key] = m.list.PushBack(&ListMapNode{key, value})
		} else {
			e.Value = &ListMapNode{key, value}
		}
	}
	m.mu.Unlock()
}

func (m *ListMap) Search(key interface{}) (value interface{}, found bool) {
	m.mu.RLock()
	if m.data != nil {
		if e, ok := m.data[key]; ok {
			value = e.Value.(*ListMapNode).value
			found = ok
		}
	}
	m.mu.RUnlock()
	return
}

func (m *ListMap) Get(key interface{}) (value interface{}) {
	m.mu.RLock()
	if m.data != nil {
		if e, ok := m.data[key]; ok {
			value = e.Value.(*ListMapNode).value
		}
	}
	m.mu.RUnlock()
	return
}

func (m *ListMap) Pop() (key, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, e := range m.data {
		value = e.Value.(*ListMapNode).value
		delete(m.data, k)
		m.list.Remove(e)
		return k, value
	}
	return
}

func (m *ListMap) Pops(size int) map[interface{}]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	if size > len(m.data) || size == -1 {
		size = len(m.data)
	}
	if size == 0 {
		return nil
	}
	index := 0
	newMap := make(map[interface{}]interface{}, size)
	for k, e := range m.data {
		value := e.Value.(*ListMapNode).value
		delete(m.data, k)
		m.list.Remove(e)
		newMap[k] = value
		index++
		if index == size {
			break
		}
	}
	return newMap
}

func (m *ListMap) doSetWithLockCheck(key interface{}, value interface{}) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[interface{}]*list.Element)
		m.list = list.New()
	}
	if e, ok := m.data[key]; ok {
		return e.Value.(*ListMapNode).value
	}
	if f, ok := value.(func() interface{}); ok {
		value = f()
	}
	if value != nil {
		m.data[key] = m.list.PushBack(&ListMapNode{key, value})
	}
	return value
}

func (m *ListMap) GetOrSet(key interface{}, value interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

func (m *ListMap) GetOrSetFunc(key interface{}, f func() interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

func (m *ListMap) GetOrSetFuncLock(key interface{}, f func() interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, f)
	} else {
		return v
	}
}

func (m *ListMap) GetVar(key interface{}) *vvar.Var {
	return vvar.New(m.Get(key))
}

func (m *ListMap) GetVarOrSet(key interface{}, value interface{}) *vvar.Var {
	return vvar.New(m.GetOrSet(key, value))
}

func (m *ListMap) GetVarOrSetFunc(key interface{}, f func() interface{}) *vvar.Var {
	return vvar.New(m.GetOrSetFunc(key, f))
}

func (m *ListMap) GetVarOrSetFuncLock(key interface{}, f func() interface{}) *vvar.Var {
	return vvar.New(m.GetOrSetFuncLock(key, f))
}

func (m *ListMap) SetIfNotExist(key interface{}, value interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

func (m *ListMap) SetIfNotExistFunc(key interface{}, f func() interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

func (m *ListMap) SetIfNotExistFuncLock(key interface{}, f func() interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

func (m *ListMap) Remove(key interface{}) (value interface{}) {
	m.mu.Lock()
	if m.data != nil {
		if e, ok := m.data[key]; ok {
			value = e.Value.(*ListMapNode).value
			delete(m.data, key)
			m.list.Remove(e)
		}
	}
	m.mu.Unlock()
	return
}

func (m *ListMap) Removes(keys []interface{}) {
	m.mu.Lock()
	if m.data != nil {
		for _, key := range keys {
			if e, ok := m.data[key]; ok {
				delete(m.data, key)
				m.list.Remove(e)
			}
		}
	}
	m.mu.Unlock()
}

func (m *ListMap) Keys() []interface{} {
	m.mu.RLock()
	var (
		keys  = make([]interface{}, m.list.Len())
		index = 0
	)
	if m.list != nil {
		m.list.IteratorAsc(func(e *list.Element) bool {
			keys[index] = e.Value.(*ListMapNode).key
			index++
			return true
		})
	}
	m.mu.RUnlock()
	return keys
}

func (m *ListMap) Values() []interface{} {
	m.mu.RLock()
	var (
		values = make([]interface{}, m.list.Len())
		index  = 0
	)
	if m.list != nil {
		m.list.IteratorAsc(func(e *list.Element) bool {
			values[index] = e.Value.(*ListMapNode).value
			index++
			return true
		})
	}
	m.mu.RUnlock()
	return values
}

func (m *ListMap) Contains(key interface{}) (ok bool) {
	m.mu.RLock()
	if m.data != nil {
		_, ok = m.data[key]
	}
	m.mu.RUnlock()
	return
}

func (m *ListMap) Size() (size int) {
	m.mu.RLock()
	size = len(m.data)
	m.mu.RUnlock()
	return
}

func (m *ListMap) IsEmpty() bool {
	return m.Size() == 0
}

func (m *ListMap) Flip() {
	data := m.Map()
	m.Clear()
	for key, value := range data {
		m.Set(value, key)
	}
}

func (m *ListMap) Merge(other *ListMap) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[interface{}]*list.Element)
		m.list = list.New()
	}
	if other != m {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	node := (*ListMapNode)(nil)
	other.list.IteratorAsc(func(e *list.Element) bool {
		node = e.Value.(*ListMapNode)
		if e, ok := m.data[node.key]; !ok {
			m.data[node.key] = m.list.PushBack(&ListMapNode{node.key, node.value})
		} else {
			e.Value = &ListMapNode{node.key, node.value}
		}
		return true
	})
}

func (m *ListMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(conv.Map(m.Map()))
}

func (m *ListMap) UnmarshalJSON(b []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[interface{}]*list.Element)
		m.list = list.New()
	}
	var data map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	for key, value := range data {
		if e, ok := m.data[key]; !ok {
			m.data[key] = m.list.PushBack(&ListMapNode{key, value})
		} else {
			e.Value = &ListMapNode{key, value}
		}
	}
	return nil
}

func (m *ListMap) UnmarshalValue(value interface{}) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[interface{}]*list.Element)
		m.list = list.New()
	}
	for k, v := range conv.Map(value) {
		if e, ok := m.data[k]; !ok {
			m.data[k] = m.list.PushBack(&ListMapNode{k, v})
		} else {
			e.Value = &ListMapNode{k, v}
		}
	}
	return
}
