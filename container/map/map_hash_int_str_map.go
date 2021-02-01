package vmap

import (
	"utils/util/json"

	"utils/util/empty"

	"utils/conv"
	"utils/util/rwmutex"
)

type IntStrMap struct {
	mu   rwmutex.RWMutex
	data map[int]string
}

func NewIntStrMap(safe ...bool) *IntStrMap {
	return &IntStrMap{
		mu:   rwmutex.Create(safe...),
		data: make(map[int]string),
	}
}

func NewIntStrMapFrom(data map[int]string, safe ...bool) *IntStrMap {
	return &IntStrMap{
		mu:   rwmutex.Create(safe...),
		data: data,
	}
}

func (m *IntStrMap) Iterator(f func(k int, v string) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

func (m *IntStrMap) Clone() *IntStrMap {
	return NewIntStrMapFrom(m.MapCopy(), !m.mu.IsSafe())
}

func (m *IntStrMap) Map() map[int]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.mu.IsSafe() {
		return m.data
	}
	data := make(map[int]string, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

func (m *IntStrMap) MapStrAny() map[string]interface{} {
	m.mu.RLock()
	data := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		data[conv.String(k)] = v
	}
	m.mu.RUnlock()
	return data
}

func (m *IntStrMap) MapCopy() map[int]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(map[int]string, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

func (m *IntStrMap) FilterEmpty() {
	m.mu.Lock()
	for k, v := range m.data {
		if empty.IsEmpty(v) {
			delete(m.data, k)
		}
	}
	m.mu.Unlock()
}

func (m *IntStrMap) Set(key int, val string) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[int]string)
	}
	m.data[key] = val
	m.mu.Unlock()
}

func (m *IntStrMap) Sets(data map[int]string) {
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

func (m *IntStrMap) Search(key int) (value string, found bool) {
	m.mu.RLock()
	if m.data != nil {
		value, found = m.data[key]
	}
	m.mu.RUnlock()
	return
}

func (m *IntStrMap) Get(key int) (value string) {
	m.mu.RLock()
	if m.data != nil {
		value, _ = m.data[key]
	}
	m.mu.RUnlock()
	return
}

func (m *IntStrMap) Pop() (key int, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for key, value = range m.data {
		delete(m.data, key)
		return
	}
	return
}

func (m *IntStrMap) Pops(size int) map[int]string {
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
		newMap = make(map[int]string, size)
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

func (m *IntStrMap) doSetWithLockCheck(key int, value string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[int]string)
	}
	if v, ok := m.data[key]; ok {
		return v
	}
	m.data[key] = value
	return value
}

func (m *IntStrMap) GetOrSet(key int, value string) string {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

func (m *IntStrMap) GetOrSetFunc(key int, f func() string) string {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

func (m *IntStrMap) GetOrSetFuncLock(key int, f func() string) string {
	if v, ok := m.Search(key); !ok {
		m.mu.Lock()
		defer m.mu.Unlock()
		if m.data == nil {
			m.data = make(map[int]string)
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

func (m *IntStrMap) SetIfNotExist(key int, value string) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

func (m *IntStrMap) SetIfNotExistFunc(key int, f func() string) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

func (m *IntStrMap) SetIfNotExistFuncLock(key int, f func() string) bool {
	if !m.Contains(key) {
		m.mu.Lock()
		defer m.mu.Unlock()
		if m.data == nil {
			m.data = make(map[int]string)
		}
		if _, ok := m.data[key]; !ok {
			m.data[key] = f()
		}
		return true
	}
	return false
}

func (m *IntStrMap) Removes(keys []int) {
	m.mu.Lock()
	if m.data != nil {
		for _, key := range keys {
			delete(m.data, key)
		}
	}
	m.mu.Unlock()
}

func (m *IntStrMap) Remove(key int) (value string) {
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

func (m *IntStrMap) Keys() []int {
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

func (m *IntStrMap) Values() []string {
	m.mu.RLock()
	var (
		values = make([]string, len(m.data))
		index  = 0
	)
	for _, value := range m.data {
		values[index] = value
		index++
	}
	m.mu.RUnlock()
	return values
}

func (m *IntStrMap) Contains(key int) bool {
	var ok bool
	m.mu.RLock()
	if m.data != nil {
		_, ok = m.data[key]
	}
	m.mu.RUnlock()
	return ok
}

func (m *IntStrMap) Size() int {
	m.mu.RLock()
	length := len(m.data)
	m.mu.RUnlock()
	return length
}

func (m *IntStrMap) IsEmpty() bool {
	return m.Size() == 0
}

func (m *IntStrMap) Clear() {
	m.mu.Lock()
	m.data = make(map[int]string)
	m.mu.Unlock()
}

func (m *IntStrMap) Replace(data map[int]string) {
	m.mu.Lock()
	m.data = data
	m.mu.Unlock()
}

func (m *IntStrMap) LockFunc(f func(m map[int]string)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f(m.data)
}

func (m *IntStrMap) RLockFunc(f func(m map[int]string)) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f(m.data)
}

func (m *IntStrMap) Flip() {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := make(map[int]string, len(m.data))
	for k, v := range m.data {
		n[conv.Int(v)] = conv.String(k)
	}
	m.data = n
}

func (m *IntStrMap) Merge(other *IntStrMap) {
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

func (m *IntStrMap) MarshalJSON() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.data)
}

func (m *IntStrMap) UnmarshalJSON(b []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[int]string)
	}
	if err := json.Unmarshal(b, &m.data); err != nil {
		return err
	}
	return nil
}

func (m *IntStrMap) UnmarshalValue(value interface{}) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[int]string)
	}
	switch value.(type) {
	case string, []byte:
		return json.Unmarshal(conv.Bytes(value), &m.data)
	default:
		for k, v := range conv.Map(value) {
			m.data[conv.Int(k)] = conv.String(v)
		}
	}
	return
}
