package vtype

import (
	"strconv"
	"sync/atomic"

	"github.com/luyingjie/utils/conv"
)

type Uint32 struct {
	value uint32
}

func NewUint32(value ...uint32) *Uint32 {
	if len(value) > 0 {
		return &Uint32{
			value: value[0],
		}
	}
	return &Uint32{}
}

func (v *Uint32) Clone() *Uint32 {
	return NewUint32(v.Val())
}

func (v *Uint32) Set(value uint32) (old uint32) {
	return atomic.SwapUint32(&v.value, value)
}

func (v *Uint32) Val() uint32 {
	return atomic.LoadUint32(&v.value)
}

func (v *Uint32) Add(delta uint32) (new uint32) {
	return atomic.AddUint32(&v.value, delta)
}

func (v *Uint32) Cas(old, new uint32) (swapped bool) {
	return atomic.CompareAndSwapUint32(&v.value, old, new)
}

func (v *Uint32) String() string {
	return strconv.FormatUint(uint64(v.Val()), 10)
}

func (v *Uint32) MarshalJSON() ([]byte, error) {
	return conv.UnsafeStrToBytes(strconv.FormatUint(uint64(v.Val()), 10)), nil
}

func (v *Uint32) UnmarshalJSON(b []byte) error {
	v.Set(conv.Uint32(conv.UnsafeBytesToStr(b)))
	return nil
}

func (v *Uint32) UnmarshalValue(value interface{}) error {
	v.Set(conv.Uint32(value))
	return nil
}
