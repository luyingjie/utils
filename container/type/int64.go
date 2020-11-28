package vtype

import (
	"strconv"
	"sync/atomic"
	"utils/convert/conv"
)

type Int64 struct {
	value int64
}

func NewInt64(value ...int64) *Int64 {
	if len(value) > 0 {
		return &Int64{
			value: value[0],
		}
	}
	return &Int64{}
}

func (v *Int64) Clone() *Int64 {
	return NewInt64(v.Val())
}

func (v *Int64) Set(value int64) (old int64) {
	return atomic.SwapInt64(&v.value, value)
}

func (v *Int64) Val() int64 {
	return atomic.LoadInt64(&v.value)
}

func (v *Int64) Add(delta int64) (new int64) {
	return atomic.AddInt64(&v.value, delta)
}

func (v *Int64) Cas(old, new int64) (swapped bool) {
	return atomic.CompareAndSwapInt64(&v.value, old, new)
}

func (v *Int64) String() string {
	return strconv.FormatInt(v.Val(), 10)
}

func (v *Int64) MarshalJSON() ([]byte, error) {
	return conv.UnsafeStrToBytes(strconv.FormatInt(v.Val(), 10)), nil
}

func (v *Int64) UnmarshalJSON(b []byte) error {
	v.Set(conv.Int64(conv.UnsafeBytesToStr(b)))
	return nil
}

func (v *Int64) UnmarshalValue(value interface{}) error {
	v.Set(conv.Int64(value))
	return nil
}
