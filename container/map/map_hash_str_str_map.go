package mymap

import (
	"utils/convert/conv"
	"utils/utils/json"

	"utils/utils/empty"

	"utils/utils/rwmutex"
)

type StrStrMap struct {
	mu   rwmutex.RWMutex
	data map[string]string
}

func NewStrStrMap(safe ...bool) *StrStrMap {
	return &StrStrMap{
		data: make(map[string]string),
		mu:   rwmutex.Create(safe...),
	}
}

func NewStrStrMapFrom(data map[string]string, safe ...bool) *StrStrMap {
	return &StrStrMap{
		mu:   rwmutex.Create(safe...),
		data: data,
	}
}

func (m *StrStrMap) Iterator(f func(k string, v string) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

func (m *StrStrMap) Clone() *StrStrMap {
	return NewStrStrMapFrom(m.MapCopy(), !m.mu.IsSafe())
}

func (m *StrStrMap) Map() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.mu.IsSafe() {
		return m.data
	}
	data := make(map[string]string, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

func (m *StrStrMap) MapStrAny() map[string]interface{} {
	m.mu.RLock()
	data := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	m.mu.RUnlock()
	return data
}

func (m *StrStrMap) MapCopy() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(map[string]string, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

func (m *StrStrMap) FilterEmpty() {
	m.mu.Lock()
	for k, v := range m.data {
		if empty.IsEmpty(v) {
			delete(m.data, k)
		}
	}
	m.mu.Unlock()
}

func (m *StrStrMap) Set(key string, val string) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[string]string)
	}
	m.data[key] = val
	m.mu.Unlock()
}

func (m *StrStrMap) Sets(data map[string]string) {
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

func (m *StrStrMap) Search(key string) (value string, found bool) {
	m.mu.RLock()
	if m.data != nil {
		value, found = m.data[key]
	}
	m.mu.RUnlock()
	return
}

func (m *StrStrMap) Get(key string) (value string) {
	m.mu.RLock()
	if m.data != nil {
		value, _ = m.data[key]
	}
	m.mu.RUnlock()
	return
}

func (m *StrStrMap) Pop() (key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for key, value = range m.data {
		delete(m.data, key)
		return
	}
	return
}

func (m *StrStrMap) Pops(size int) map[string]string {
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
		newMap = make(map[string]string, size)
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

func (m *StrStrMap) doSetWithLockCheck(key string, value string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[string]string)
	}
	if v, ok := m.data[key]; ok {
		return v
	}
	m.data[key] = value
	return value
}

func (m *StrStrMap) GetOrSet(key string, value string) string {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

func (m *StrStrMap) GetOrSetFunc(key string, f func() string) string {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

func (m *StrStrMap) GetOrSetFuncLock(key string, f func() string) string {
	if v, ok := m.Search(key); !ok {
		m.mu.Lock()
		defer m.mu.Unlock()
		if m.data == nil {
			m.data = make(map[string]string)
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

func (m *StrStrMap) SetIfNotExist(key string, value string) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

func (m *StrStrMap) SetIfNotExistFunc(key string, f func() string) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

func (m *StrStrMap) SetIfNotExistFuncLock(key string, f func() string) bool {
	if !m.Contains(key) {
		m.mu.Lock()
		defer m.mu.Unlock()
		if m.data == nil {
			m.data = make(map[string]string)
		}
		if _, ok := m.data[key]; !ok {
			m.data[key] = f()
		}
		return true
	}
	return false
}

func (m *StrStrMap) Removes(keys []string) {
	m.mu.Lock()
	if m.data != nil {
		for _, key := range keys {
			delete(m.data, key)
		}
	}
	m.mu.Unlock()
}

func (m *StrStrMap) Remove(key string) (value string) {
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

func (m *StrStrMap) Keys() []string {
	m.mu.RLock()
	var (
		keys  = make([]string, len(m.data))
		index = 0
	)
	for key := range m.data {
		keys[index] = key
		index++
	}
	m.mu.RUnlock()
	return keys
}

func (m *StrStrMap) Values() []string {
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

func (m *StrStrMap) Contains(key string) bool {
	var ok bool
	m.mu.RLock()
	if m.data != nil {
		_, ok = m.data[key]
	}
	m.mu.RUnlock()
	return ok
}

func (m *StrStrMap) Size() int {
	m.mu.RLock()
	length := len(m.data)
	m.mu.RUnlock()
	return length
}

func (m *StrStrMap) IsEmpty() bool {
	return m.Size() == 0
}

func (m *StrStrMap) Clear() {
	m.mu.Lock()
	m.data = make(map[string]string)
	m.mu.Unlock()
}

func (m *StrStrMap) Replace(data map[string]string) {
	m.mu.Lock()
	m.data = data
	m.mu.Unlock()
}

func (m *StrStrMap) LockFunc(f func(m map[string]string)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f(m.data)
}

func (m *StrStrMap) RLockFunc(f func(m map[string]string)) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f(m.data)
}

func (m *StrStrMap) Flip() {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := make(map[string]string, len(m.data))
	for k, v := range m.data {
		n[v] = k
	}
	m.data = n
}

func (m *StrStrMap) Merge(other *StrStrMap) {
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

func (m *StrStrMap) MarshalJSON() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.data)
}

func (m *StrStrMap) UnmarshalJSON(b []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[string]string)
	}
	if err := json.Unmarshal(b, &m.data); err != nil {
		return err
	}
	return nil
}

func (m *StrStrMap) UnmarshalValue(value interface{}) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = conv.MapStrStr(value)
	return
}
