package vtype

import (
	"strconv"
	"sync/atomic"
	"utils/convert/conv"
)

type Uint struct {
	value uint64
}

func NewUint(value ...uint) *Uint {
	if len(value) > 0 {
		return &Uint{
			value: uint64(value[0]),
		}
	}
	return &Uint{}
}

func (v *Uint) Clone() *Uint {
	return NewUint(v.Val())
}

func (v *Uint) Set(value uint) (old uint) {
	return uint(atomic.SwapUint64(&v.value, uint64(value)))
}

func (v *Uint) Val() uint {
	return uint(atomic.LoadUint64(&v.value))
}

func (v *Uint) Add(delta uint) (new uint) {
	return uint(atomic.AddUint64(&v.value, uint64(delta)))
}

func (v *Uint) Cas(old, new uint) (swapped bool) {
	return atomic.CompareAndSwapUint64(&v.value, uint64(old), uint64(new))
}

func (v *Uint) String() string {
	return strconv.FormatUint(uint64(v.Val()), 10)
}

func (v *Uint) MarshalJSON() ([]byte, error) {
	return conv.UnsafeStrToBytes(strconv.FormatUint(uint64(v.Val()), 10)), nil
}

func (v *Uint) UnmarshalJSON(b []byte) error {
	v.Set(conv.Uint(conv.UnsafeBytesToStr(b)))
	return nil
}

func (v *Uint) UnmarshalValue(value interface{}) error {
	v.Set(conv.Uint(value))
	return nil
}
