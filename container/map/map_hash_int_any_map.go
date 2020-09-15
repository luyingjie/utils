package mymap

import (
	myvar "utils/container/var"
	"utils/utils/json"

	"utils/utils/empty"

	"utils/convert/conv"
	"utils/utils/rwmutex"
)

type IntAnyMap struct {
	mu   rwmutex.RWMutex
	data map[int]interface{}
}

func NewIntAnyMap(safe ...bool) *IntAnyMap {
	return &IntAnyMap{
		mu:   rwmutex.Create(safe...),
		data: make(map[int]interface{}),
	}
}

func NewIntAnyMapFrom(data map[int]interface{}, safe ...bool) *IntAnyMap {
	return &IntAnyMap{
		mu:   rwmutex.Create(safe...),
		data: data,
	}
}

func (m *IntAnyMap) Iterator(f func(k int, v interface{}) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

func (m *IntAnyMap) Clone() *IntAnyMap {
	return NewIntAnyMapFrom(m.MapCopy(), !m.mu.IsSafe())
}

func (m *IntAnyMap) Map() map[int]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.mu.IsSafe() {
		return m.data
	}
	data := make(map[int]interface{}, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

func (m *IntAnyMap) MapStrAny() map[string]interface{} {
	m.mu.RLock()
	data := make(map[string]interface{}, len(m.data))
	for k, v := range m.data {
		data[conv.String(k)] = v
	}
	m.mu.RUnlock()
	return data
}

func (m *IntAnyMap) MapCopy() map[int]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(map[int]interface{}, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

func (m *IntAnyMap) FilterEmpty() {
	m.mu.Lock()
	for k, v := range m.data {
		if empty.IsEmpty(v) {
			delete(m.data, k)
		}
	}
	m.mu.Unlock()
}

func (m *IntAnyMap) FilterNil() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		if empty.IsNil(v) {
			delete(m.data, k)
		}
	}
}

func (m *IntAnyMap) Set(key int, val interface{}) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(map[int]interface{})
	}
	m.data[key] = val
	m.mu.Unlock()
}

func (m *IntAnyMap) Sets(data map[int]interface{}) {
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

func (m *IntAnyMap) Search(key int) (value interface{}, found bool) {
	m.mu.RLock()
	if m.data != nil {
		value, found = m.data[key]
	}
	m.mu.RUnlock()
	return
}

func (m *IntAnyMap) Get(key int) (value interface{}) {
	m.mu.RLock()
	if m.data != nil {
		value, _ = m.data[key]
	}
	m.mu.RUnlock()
	return
}

func (m *IntAnyMap) Pop() (key int, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for key, value = range m.data {
		delete(m.data, key)
		return
	}
	return
}

func (m *IntAnyMap) Pops(size int) map[int]interface{} {
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
		newMap = make(map[int]interface{}, size)
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

func (m *IntAnyMap) doSetWithLockCheck(key int, value interface{}) interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[int]interface{})
	}
	if v, ok := m.data[key]; ok {
		return v
	}
	if f, ok := value.(func() interface{}); ok {
		value = f()
	}
	if value != nil {
		m.data[key] = value
	}
	return value
}

func (m *IntAnyMap) GetOrSet(key int, value interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

func (m *IntAnyMap) GetOrSetFunc(key int, f func() interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

func (m *IntAnyMap) GetOrSetFuncLock(key int, f func() interface{}) interface{} {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, f)
	} else {
		return v
	}
}

func (m *IntAnyMap) GetVar(key int) *myvar.Var {
	return myvar.New(m.Get(key))
}

func (m *IntAnyMap) GetVarOrSet(key int, value interface{}) *myvar.Var {
	return myvar.New(m.GetOrSet(key, value))
}

func (m *IntAnyMap) GetVarOrSetFunc(key int, f func() interface{}) *myvar.Var {
	return myvar.New(m.GetOrSetFunc(key, f))
}

func (m *IntAnyMap) GetVarOrSetFuncLock(key int, f func() interface{}) *myvar.Var {
	return myvar.New(m.GetOrSetFuncLock(key, f))
}

func (m *IntAnyMap) SetIfNotExist(key int, value interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

func (m *IntAnyMap) SetIfNotExistFunc(key int, f func() interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

func (m *IntAnyMap) SetIfNotExistFuncLock(key int, f func() interface{}) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

func (m *IntAnyMap) Removes(keys []int) {
	m.mu.Lock()
	if m.data != nil {
		for _, key := range keys {
			delete(m.data, key)
		}
	}
	m.mu.Unlock()
}

func (m *IntAnyMap) Remove(key int) (value interface{}) {
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

func (m *IntAnyMap) Keys() []int {
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

func (m *IntAnyMap) Values() []interface{} {
	m.mu.RLock()
	var (
		values = make([]interface{}, len(m.data))
		index  = 0
	)
	for _, value := range m.data {
		values[index] = value
		index++
	}
	m.mu.RUnlock()
	return values
}

func (m *IntAnyMap) Contains(key int) bool {
	var ok bool
	m.mu.RLock()
	if m.data != nil {
		_, ok = m.data[key]
	}
	m.mu.RUnlock()
	return ok
}

func (m *IntAnyMap) Size() int {
	m.mu.RLock()
	length := len(m.data)
	m.mu.RUnlock()
	return length
}

func (m *IntAnyMap) IsEmpty() bool {
	return m.Size() == 0
}

func (m *IntAnyMap) Clear() {
	m.mu.Lock()
	m.data = make(map[int]interface{})
	m.mu.Unlock()
}

func (m *IntAnyMap) Replace(data map[int]interface{}) {
	m.mu.Lock()
	m.data = data
	m.mu.Unlock()
}

func (m *IntAnyMap) LockFunc(f func(m map[int]interface{})) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f(m.data)
}

func (m *IntAnyMap) RLockFunc(f func(m map[int]interface{})) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f(m.data)
}

func (m *IntAnyMap) Flip() {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := make(map[int]interface{}, len(m.data))
	for k, v := range m.data {
		n[conv.Int(v)] = k
	}
	m.data = n
}

func (m *IntAnyMap) Merge(other *IntAnyMap) {
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

func (m *IntAnyMap) MarshalJSON() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.data)
}

func (m *IntAnyMap) UnmarshalJSON(b []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[int]interface{})
	}
	if err := json.Unmarshal(b, &m.data); err != nil {
		return err
	}
	return nil
}

func (m *IntAnyMap) UnmarshalValue(value interface{}) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(map[int]interface{})
	}
	switch value.(type) {
	case string, []byte:
		return json.Unmarshal(conv.Bytes(value), &m.data)
	default:
		for k, v := range conv.Map(value) {
			m.data[conv.Int(k)] = v
		}
	}
	return
}
