package vtype

import (
	"strconv"
	"sync/atomic"

	"github.com/luyingjie/utils/conv"
)

type Int struct {
	value int64
}

func NewInt(value ...int) *Int {
	if len(value) > 0 {
		return &Int{
			value: int64(value[0]),
		}
	}
	return &Int{}
}

func (v *Int) Clone() *Int {
	return NewInt(v.Val())
}

func (v *Int) Set(value int) (old int) {
	return int(atomic.SwapInt64(&v.value, int64(value)))
}

func (v *Int) Val() int {
	return int(atomic.LoadInt64(&v.value))
}

func (v *Int) Add(delta int) (new int) {
	return int(atomic.AddInt64(&v.value, int64(delta)))
}

func (v *Int) Cas(old, new int) (swapped bool) {
	return atomic.CompareAndSwapInt64(&v.value, int64(old), int64(new))
}

func (v *Int) String() string {
	return strconv.Itoa(v.Val())
}

func (v *Int) MarshalJSON() ([]byte, error) {
	return conv.UnsafeStrToBytes(strconv.Itoa(v.Val())), nil
}

func (v *Int) UnmarshalJSON(b []byte) error {
	v.Set(conv.Int(conv.UnsafeBytesToStr(b)))
	return nil
}

func (v *Int) UnmarshalValue(value interface{}) error {
	v.Set(conv.Int(value))
	return nil
}
