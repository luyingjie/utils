package vtype

import (
	"strconv"
	"sync/atomic"
	"utils/convert/conv"
)

type Uint64 struct {
	value uint64
}

func NewUint64(value ...uint64) *Uint64 {
	if len(value) > 0 {
		return &Uint64{
			value: value[0],
		}
	}
	return &Uint64{}
}

func (v *Uint64) Clone() *Uint64 {
	return NewUint64(v.Val())
}

func (v *Uint64) Set(value uint64) (old uint64) {
	return atomic.SwapUint64(&v.value, value)
}

func (v *Uint64) Val() uint64 {
	return atomic.LoadUint64(&v.value)
}

func (v *Uint64) Add(delta uint64) (new uint64) {
	return atomic.AddUint64(&v.value, delta)
}

func (v *Uint64) Cas(old, new uint64) (swapped bool) {
	return atomic.CompareAndSwapUint64(&v.value, old, new)
}

func (v *Uint64) String() string {
	return strconv.FormatUint(v.Val(), 10)
}

func (v *Uint64) MarshalJSON() ([]byte, error) {
	return conv.UnsafeStrToBytes(strconv.FormatUint(v.Val(), 10)), nil
}

func (v *Uint64) UnmarshalJSON(b []byte) error {
	v.Set(conv.Uint64(conv.UnsafeBytesToStr(b)))
	return nil
}

func (v *Uint64) UnmarshalValue(value interface{}) error {
	v.Set(conv.Uint64(value))
	return nil
}
