package mytype

import (
	"strconv"
	"sync/atomic"
	"utils/convert/conv"
)

type Int32 struct {
	value int32
}

func NewInt32(value ...int32) *Int32 {
	if len(value) > 0 {
		return &Int32{
			value: value[0],
		}
	}
	return &Int32{}
}

func (v *Int32) Clone() *Int32 {
	return NewInt32(v.Val())
}

func (v *Int32) Set(value int32) (old int32) {
	return atomic.SwapInt32(&v.value, value)
}

func (v *Int32) Val() int32 {
	return atomic.LoadInt32(&v.value)
}

func (v *Int32) Add(delta int32) (new int32) {
	return atomic.AddInt32(&v.value, delta)
}

func (v *Int32) Cas(old, new int32) (swapped bool) {
	return atomic.CompareAndSwapInt32(&v.value, old, new)
}

func (v *Int32) String() string {
	return strconv.Itoa(int(v.Val()))
}

func (v *Int32) MarshalJSON() ([]byte, error) {
	return conv.UnsafeStrToBytes(strconv.Itoa(int(v.Val()))), nil
}

func (v *Int32) UnmarshalJSON(b []byte) error {
	v.Set(conv.Int32(conv.UnsafeBytesToStr(b)))
	return nil
}

func (v *Int32) UnmarshalValue(value interface{}) error {
	v.Set(conv.Int32(value))
	return nil
}
