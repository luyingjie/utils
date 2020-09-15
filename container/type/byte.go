package mytype

import (
	"strconv"
	"sync/atomic"
	"utils/convert/conv"
)

type Byte struct {
	value int32
}

func NewByte(value ...byte) *Byte {
	if len(value) > 0 {
		return &Byte{
			value: int32(value[0]),
		}
	}
	return &Byte{}
}

func (v *Byte) Clone() *Byte {
	return NewByte(v.Val())
}

func (v *Byte) Set(value byte) (old byte) {
	return byte(atomic.SwapInt32(&v.value, int32(value)))
}

func (v *Byte) Val() byte {
	return byte(atomic.LoadInt32(&v.value))
}

func (v *Byte) Add(delta byte) (new byte) {
	return byte(atomic.AddInt32(&v.value, int32(delta)))
}

func (v *Byte) Cas(old, new byte) (swapped bool) {
	return atomic.CompareAndSwapInt32(&v.value, int32(old), int32(new))
}

func (v *Byte) String() string {
	return strconv.FormatUint(uint64(v.Val()), 10)
}

func (v *Byte) MarshalJSON() ([]byte, error) {
	return conv.UnsafeStrToBytes(strconv.FormatUint(uint64(v.Val()), 10)), nil
}

func (v *Byte) UnmarshalJSON(b []byte) error {
	v.Set(conv.Uint8(conv.UnsafeBytesToStr(b)))
	return nil
}

func (v *Byte) UnmarshalValue(value interface{}) error {
	v.Set(conv.Byte(value))
	return nil
}
