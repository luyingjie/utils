package mymap

import (
	"utils/convert/conv"
	"utils/utils/json"

	"utils/utils/empty"

	"utils/utils/rwmutex"
)

type IntIntMap struct {
	mu   rwmutex.RWMutex
	data map[int]int
}

func NewIntIntMap(safe ...bool) *IntIntMap {
	return &IntIntMap{
		mu:   rwmutex.Create(safe...),
		data: make(map[int]int),
	}
}

func NewIntIntMapFrom(data map[int]int, safe ...bool) *IntIntMap {
	return &IntIntMap{
		mu:   rwmutex.Create(safe...),
		data: data,
	}
}

func (m *IntIntMap) Iterator(f func(k int, v int) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

func (m *IntIntMap) Clone() *IntIntMap {
	return NewIntIntMapFrom(m.MapCopy(), !m.mu.IsSafe())
}

func (m *IntIntMap) Map() map[int]int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.mu.IsSafe() {
		return m.data
	}
	data := make(map[int]int, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

func (m *IntIntMap) MapStrAny() map[string]interface{} {
	m.mu.RLock()
	data := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		data[conv.String(k)] = v
	}
	m.mu.RUnlock()
	return data
}

func (m *IntIntMap) MapCopy() map[int]int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(map[int]int, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

func (m *IntIntMap) FilterEmpty() {
	m.mu.Lock()
	for k, v := range m.data {
		if empty.IsEmpty(v) {
			delete(m.data, k)
		}
	}
	m.mu.Unlock()
}

func (m *IntIntMap) Set(key int, val int) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[int]int)
	}
	m.data[key] = val
	m.mu.Unlock()
}

func (m *IntIntMap) Sets(data map[int]int) {
	m.mu.Lock()
	if m.data == nil {
		m.data = data
	} else {
		for k, v := range data {
			m.data[k] = v
		}
	}
	m.mu.Unlock()
}

func (m *IntIntMap) Search(key int) (value int, found bool) {
	m.mu.RLock()
	if m.data != nil {
		value, found = m.data[key]
	}
	m.mu.RUnlock()
	return
}

func (m *IntIntMap) Get(key int) (value int) {
	m.mu.RLock()
	if m.data != nil {
		value, _ = m.data[key]
	}
	m.mu.RUnlock()
	return
}

func (m *IntIntMap) Pop() (key, value int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for key, value = range m.data {
		delete(m.data, key)
		return
	}
	return
}

func (m *IntIntMap) Pops(size int) map[int]int {
	m.mu.Lock()
	defer m.mu.Unlock()
	if size > len(m.data) || size == -1 {
		size = len(m.data)
	}
	if size == 0 {
		return nil
	}
	var (
		index  = 0
		newMap = make(map[int]int, size)
	)
	for k, v := range m.data {
		delete(m.data, k)
		newMap[k] = v
		index++
		if index == size {
			break
		}
	}
	return newMap
}

func (m *IntIntMap) doSetWithLockCheck(key int, value int) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[int]int)
	}
	if v, ok := m.data[key]; ok {
		return v
	}
	m.data[key] = value
	return value
}

func (m *IntIntMap) GetOrSet(key int, value int) int {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

func (m *IntIntMap) GetOrSetFunc(key int, f func() int) int {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

func (m *IntIntMap) GetOrSetFuncLock(key int, f func() int) int {
	if v, ok := m.Search(key); !ok {
		m.mu.Lock()
		defer m.mu.Unlock()
		if m.data == nil {
			m.data = make(map[int]int)
		}
		if v, ok = m.data[key]; ok {
			return v
		}
		v = f()
		m.data[key] = v
		return v
	} else {
		return v
	}
}

func (m *IntIntMap) SetIfNotExist(key int, value int) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

func (m *IntIntMap) SetIfNotExistFunc(key int, f func() int) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

func (m *IntIntMap) SetIfNotExistFuncLock(key int, f func() int) bool {
	if !m.Contains(key) {
		m.mu.Lock()
		defer m.mu.Unlock()
		if m.data == nil {
			m.data = make(map[int]int)
		}
		if _, ok := m.data[key]; !ok {
			m.data[key] = f()
		}
		return true
	}
	return false
}

func (m *IntIntMap) Removes(keys []int) {
	m.mu.Lock()
	if m.data != nil {
		for _, key := range keys {
			delete(m.data, key)
		}
	}
	m.mu.Unlock()
}

func (m *IntIntMap) Remove(key int) (value int) {
	m.mu.Lock()
	if m.data != nil {
		var ok bool
		if value, ok = m.data[key]; ok {
			delete(m.data, key)
		}
	}
	m.mu.Unlock()
	return
}

func (m *IntIntMap) Keys() []int {
	m.mu.RLock()
	var (
		keys  = make([]int, len(m.data))
		index = 0
	)
	for key := range m.data {
		keys[index] = key
		index++
	}
	m.mu.RUnlock()
	return keys
}

func (m *IntIntMap) Values() []int {
	m.mu.RLock()
	var (
		values = make([]int, len(m.data))
		index  = 0
	)
	for _, value := range m.data {
		values[index] = value
		index++
	}
	m.mu.RUnlock()
	return values
}

func (m *IntIntMap) Contains(key int) bool {
	var ok bool
	m.mu.RLock()
	if m.data != nil {
		_, ok = m.data[key]
	}
	m.mu.RUnlock()
	return ok
}

func (m *IntIntMap) Size() int {
	m.mu.RLock()
	length := len(m.data)
	m.mu.RUnlock()
	return length
}

func (m *IntIntMap) IsEmpty() bool {
	return m.Size() == 0
}

func (m *IntIntMap) Clear() {
	m.mu.Lock()
	m.data = make(map[int]int)
	m.mu.Unlock()
}

func (m *IntIntMap) Replace(data map[int]int) {
	m.mu.Lock()
	m.data = data
	m.mu.Unlock()
}

func (m *IntIntMap) LockFunc(f func(m map[int]int)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f(m.data)
}

func (m *IntIntMap) RLockFunc(f func(m map[int]int)) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f(m.data)
}

func (m *IntIntMap) Flip() {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := make(map[int]int, len(m.data))
	for k, v := range m.data {
		n[v] = k
	}
	m.data = n
}

func (m *IntIntMap) Merge(other *IntIntMap) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = other.MapCopy()
		return
	}
	if other != m {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	for k, v := range other.data {
		m.data[k] = v
	}
}

func (m *IntIntMap) MarshalJSON() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.data)
}

func (m *IntIntMap) UnmarshalJSON(b []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[int]int)
	}
	if err := json.Unmarshal(b, &m.data); err != nil {
		return err
	}
	return nil
}

func (m *IntIntMap) UnmarshalValue(value interface{}) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[int]int)
	}
	switch value.(type) {
	case string, []byte:
		return json.Unmarshal(conv.Bytes(value), &m.data)
	default:
		for k, v := range conv.Map(value) {
			m.data[conv.Int(k)] = conv.Int(v)
		}
	}
	return
}
