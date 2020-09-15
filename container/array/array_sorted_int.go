package array

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"utils/utils/json"

	"utils/utils/rwmutex"

	"utils/convert/conv"
	"utils/utils/rand"
)

type SortedIntArray struct {
	mu         rwmutex.RWMutex
	array      []int
	unique     bool
	comparator func(a, b int) int
}

func NewSortedIntArray(safe ...bool) *SortedIntArray {
	return NewSortedIntArraySize(0, safe...)
}

func NewSortedIntArrayComparator(comparator func(a, b int) int, safe ...bool) *SortedIntArray {
	array := NewSortedIntArray(safe...)
	array.comparator = comparator
	return array
}

func NewSortedIntArraySize(cap int, safe ...bool) *SortedIntArray {
	return &SortedIntArray{
		mu:         rwmutex.Create(safe...),
		array:      make([]int, 0, cap),
		comparator: defaultComparatorInt,
	}
}

func NewSortedIntArrayRange(start, end, step int, safe ...bool) *SortedIntArray {
	if step == 0 {
		panic(fmt.Sprintf(`invalid step value: %d`, step))
	}
	slice := make([]int, (end-start+1)/step)
	index := 0
	for i := start; i <= end; i += step {
		slice[index] = i
		index++
	}
	return NewSortedIntArrayFrom(slice, safe...)
}

func NewSortedIntArrayFrom(array []int, safe ...bool) *SortedIntArray {
	a := NewSortedIntArraySize(0, safe...)
	a.array = array
	sort.Ints(a.array)
	return a
}

func NewSortedIntArrayFromCopy(array []int, safe ...bool) *SortedIntArray {
	newArray := make([]int, len(array))
	copy(newArray, array)
	return NewSortedIntArrayFrom(newArray, safe...)
}

func (a *SortedIntArray) SetArray(array []int) *SortedIntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.array = array
	quickSortInt(a.array, a.getComparator())
	return a
}

func (a *SortedIntArray) Sort() *SortedIntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	quickSortInt(a.array, a.getComparator())
	return a
}

func (a *SortedIntArray) Add(values ...int) *SortedIntArray {
	return a.Append(values...)
}

func (a *SortedIntArray) Append(values ...int) *SortedIntArray {
	if len(values) == 0 {
		return a
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, value := range values {
		index, cmp := a.binSearch(value, false)
		if a.unique && cmp == 0 {
			continue
		}
		if index < 0 {
			a.array = append(a.array, value)
			continue
		}
		if cmp > 0 {
			index++
		}
		rear := append([]int{}, a.array[index:]...)
		a.array = append(a.array[0:index], value)
		a.array = append(a.array, rear...)
	}
	return a
}

func (a *SortedIntArray) Get(index int) (value int, found bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if index < 0 || index >= len(a.array) {
		return 0, false
	}
	return a.array[index], true
}

func (a *SortedIntArray) Remove(index int) (value int, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.doRemoveWithoutLock(index)
}

func (a *SortedIntArray) doRemoveWithoutLock(index int) (value int, found bool) {
	if index < 0 || index >= len(a.array) {
		return 0, false
	}

	if index == 0 {
		value := a.array[0]
		a.array = a.array[1:]
		return value, true
	} else if index == len(a.array)-1 {
		value := a.array[index]
		a.array = a.array[:index]
		return value, true
	}

	value = a.array[index]
	a.array = append(a.array[:index], a.array[index+1:]...)
	return value, true
}

func (a *SortedIntArray) RemoveValue(value int) bool {
	if i := a.Search(value); i != -1 {
		_, found := a.Remove(i)
		return found
	}
	return false
}

func (a *SortedIntArray) PopLeft() (value int, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.array) == 0 {
		return 0, false
	}
	value = a.array[0]
	a.array = a.array[1:]
	return value, true
}

func (a *SortedIntArray) PopRight() (value int, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	index := len(a.array) - 1
	if index < 0 {
		return 0, false
	}
	value = a.array[index]
	a.array = a.array[:index]
	return value, true
}

func (a *SortedIntArray) PopRand() (value int, found bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.doRemoveWithoutLock(rand.Intn(len(a.array)))
}

func (a *SortedIntArray) PopRands(size int) []int {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size <= 0 || len(a.array) == 0 {
		return nil
	}
	if size >= len(a.array) {
		size = len(a.array)
	}
	array := make([]int, size)
	for i := 0; i < size; i++ {
		array[i], _ = a.doRemoveWithoutLock(rand.Intn(len(a.array)))
	}
	return array
}

func (a *SortedIntArray) PopLefts(size int) []int {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size <= 0 || len(a.array) == 0 {
		return nil
	}
	if size >= len(a.array) {
		array := a.array
		a.array = a.array[:0]
		return array
	}
	value := a.array[0:size]
	a.array = a.array[size:]
	return value
}

func (a *SortedIntArray) PopRights(size int) []int {
	a.mu.Lock()
	defer a.mu.Unlock()
	if size <= 0 || len(a.array) == 0 {
		return nil
	}
	index := len(a.array) - size
	if index <= 0 {
		array := a.array
		a.array = a.array[:0]
		return array
	}
	value := a.array[index:]
	a.array = a.array[:index]
	return value
}

func (a *SortedIntArray) Range(start int, end ...int) []int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	offsetEnd := len(a.array)
	if len(end) > 0 && end[0] < offsetEnd {
		offsetEnd = end[0]
	}
	if start > offsetEnd {
		return nil
	}
	if start < 0 {
		start = 0
	}
	array := ([]int)(nil)
	if a.mu.IsSafe() {
		array = make([]int, offsetEnd-start)
		copy(array, a.array[start:offsetEnd])
	} else {
		array = a.array[start:offsetEnd]
	}
	return array
}

func (a *SortedIntArray) SubSlice(offset int, length ...int) []int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	size := len(a.array)
	if len(length) > 0 {
		size = length[0]
	}
	if offset > len(a.array) {
		return nil
	}
	if offset < 0 {
		offset = len(a.array) + offset
		if offset < 0 {
			return nil
		}
	}
	if size < 0 {
		offset += size
		size = -size
		if offset < 0 {
			return nil
		}
	}
	end := offset + size
	if end > len(a.array) {
		end = len(a.array)
		size = len(a.array) - offset
	}
	if a.mu.IsSafe() {
		s := make([]int, size)
		copy(s, a.array[offset:])
		return s
	} else {
		return a.array[offset:end]
	}
}

func (a *SortedIntArray) Len() int {
	a.mu.RLock()
	length := len(a.array)
	a.mu.RUnlock()
	return length
}

// Sum returns the sum of values in an array.
func (a *SortedIntArray) Sum() (sum int) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, v := range a.array {
		sum += v
	}
	return
}

func (a *SortedIntArray) Slice() []int {
	array := ([]int)(nil)
	if a.mu.IsSafe() {
		a.mu.RLock()
		defer a.mu.RUnlock()
		array = make([]int, len(a.array))
		copy(array, a.array)
	} else {
		array = a.array
	}
	return array
}

func (a *SortedIntArray) Interfaces() []interface{} {
	a.mu.RLock()
	defer a.mu.RUnlock()
	array := make([]interface{}, len(a.array))
	for k, v := range a.array {
		array[k] = v
	}
	return array
}

func (a *SortedIntArray) Contains(value int) bool {
	return a.Search(value) != -1
}

func (a *SortedIntArray) Search(value int) (index int) {
	if i, r := a.binSearch(value, true); r == 0 {
		return i
	}
	return -1
}

func (a *SortedIntArray) binSearch(value int, lock bool) (index int, result int) {
	if lock {
		a.mu.RLock()
		defer a.mu.RUnlock()
	}
	if len(a.array) == 0 {
		return -1, -2
	}
	min := 0
	max := len(a.array) - 1
	mid := 0
	cmp := -2
	for min <= max {
		mid = (min + max) / 2
		cmp = a.getComparator()(value, a.array[mid])
		switch {
		case cmp < 0:
			max = mid - 1
		case cmp > 0:
			min = mid + 1
		default:
			return mid, cmp
		}
	}
	return mid, cmp
}

func (a *SortedIntArray) SetUnique(unique bool) *SortedIntArray {
	oldUnique := a.unique
	a.unique = unique
	if unique && oldUnique != unique {
		a.Unique()
	}
	return a
}

func (a *SortedIntArray) Unique() *SortedIntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.array) == 0 {
		return a
	}
	i := 0
	for {
		if i == len(a.array)-1 {
			break
		}
		if a.getComparator()(a.array[i], a.array[i+1]) == 0 {
			a.array = append(a.array[:i+1], a.array[i+1+1:]...)
		} else {
			i++
		}
	}
	return a
}

func (a *SortedIntArray) Clone() (newArray *SortedIntArray) {
	a.mu.RLock()
	array := make([]int, len(a.array))
	copy(array, a.array)
	a.mu.RUnlock()
	return NewSortedIntArrayFrom(array, !a.mu.IsSafe())
}

func (a *SortedIntArray) Clear() *SortedIntArray {
	a.mu.Lock()
	if len(a.array) > 0 {
		a.array = make([]int, 0)
	}
	a.mu.Unlock()
	return a
}

func (a *SortedIntArray) LockFunc(f func(array []int)) *SortedIntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	f(a.array)
	return a
}

func (a *SortedIntArray) RLockFunc(f func(array []int)) *SortedIntArray {
	a.mu.RLock()
	defer a.mu.RUnlock()
	f(a.array)
	return a
}

func (a *SortedIntArray) Merge(array interface{}) *SortedIntArray {
	return a.Add(conv.Ints(array)...)
}

func (a *SortedIntArray) Chunk(size int) [][]int {
	if size < 1 {
		return nil
	}
	a.mu.RLock()
	defer a.mu.RUnlock()
	length := len(a.array)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n [][]int
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, a.array[i*size:end])
		i++
	}
	return n
}

func (a *SortedIntArray) Rand() (value int, found bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.array) == 0 {
		return 0, false
	}
	return a.array[rand.Intn(len(a.array))], true
}

func (a *SortedIntArray) Rands(size int) []int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if size <= 0 || len(a.array) == 0 {
		return nil
	}
	array := make([]int, size)
	for i := 0; i < size; i++ {
		array[i] = a.array[rand.Intn(len(a.array))]
	}
	return array
}

func (a *SortedIntArray) Join(glue string) string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.array) == 0 {
		return ""
	}
	buffer := bytes.NewBuffer(nil)
	for k, v := range a.array {
		buffer.WriteString(conv.String(v))
		if k != len(a.array)-1 {
			buffer.WriteString(glue)
		}
	}
	return buffer.String()
}

func (a *SortedIntArray) CountValues() map[int]int {
	m := make(map[int]int)
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, v := range a.array {
		m[v]++
	}
	return m
}

func (a *SortedIntArray) Iterator(f func(k int, v int) bool) {
	a.IteratorAsc(f)
}

func (a *SortedIntArray) IteratorAsc(f func(k int, v int) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for k, v := range a.array {
		if !f(k, v) {
			break
		}
	}
}

func (a *SortedIntArray) IteratorDesc(f func(k int, v int) bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for i := len(a.array) - 1; i >= 0; i-- {
		if !f(i, a.array[i]) {
			break
		}
	}
}

func (a *SortedIntArray) String() string {
	return "[" + a.Join(",") + "]"
}

func (a SortedIntArray) MarshalJSON() ([]byte, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return json.Marshal(a.array)
}

func (a *SortedIntArray) UnmarshalJSON(b []byte) error {
	if a.comparator == nil {
		a.array = make([]int, 0)
		a.comparator = defaultComparatorInt
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := json.Unmarshal(b, &a.array); err != nil {
		return err
	}
	if a.array != nil {
		sort.Ints(a.array)
	}
	return nil
}

func (a *SortedIntArray) UnmarshalValue(value interface{}) (err error) {
	if a.comparator == nil {
		a.comparator = defaultComparatorInt
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	switch value.(type) {
	case string, []byte:
		err = json.Unmarshal(conv.Bytes(value), &a.array)
	default:
		a.array = conv.SliceInt(value)
	}
	if a.array != nil {
		sort.Ints(a.array)
	}
	return err
}

func (a *SortedIntArray) FilterEmpty() *SortedIntArray {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i := 0; i < len(a.array); {
		if a.array[i] == 0 {
			a.array = append(a.array[:i], a.array[i+1:]...)
		} else {
			break
		}
	}
	for i := len(a.array) - 1; i >= 0; {
		if a.array[i] == 0 {
			a.array = append(a.array[:i], a.array[i+1:]...)
		} else {
			break
		}
	}
	return a
}

func (a *SortedIntArray) Walk(f func(value int) int) *SortedIntArray {
	a.mu.Lock()
	defer a.mu.Unlock()

	defer quickSortInt(a.array, a.getComparator())

	for i, v := range a.array {
		a.array[i] = f(v)
	}
	return a
}

func (a *SortedIntArray) IsEmpty() bool {
	return a.Len() == 0
}

func (a *SortedIntArray) getComparator() func(a, b int) int {
	if a.comparator == nil {
		return defaultComparatorInt
	}
	return a.comparator
}
