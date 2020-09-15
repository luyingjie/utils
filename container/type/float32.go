package mytype

import (
	"math"
	"strconv"
	"sync/atomic"
	"unsafe"
	"utils/convert/conv"
)

type Float32 struct {
	value uint32
}

func NewFloat32(value ...float32) *Float32 {
	if len(value) > 0 {
		return &Float32{
			value: math.Float32bits(value[0]),
		}
	}
	return &Float32{}
}

func (v *Float32) Clone() *Float32 {
	return NewFloat32(v.Val())
}

func (v *Float32) Set(value float32) (old float32) {
	return math.Float32frombits(atomic.SwapUint32(&v.value, math.Float32bits(value)))
}

func (v *Float32) Val() float32 {
	return math.Float32frombits(atomic.LoadUint32(&v.value))
}

func (v *Float32) Add(delta float32) (new float32) {
	for {
		old := math.Float32frombits(v.value)
		new = old + delta
		if atomic.CompareAndSwapUint32(
			(*uint32)(unsafe.Pointer(&v.value)),
			math.Float32bits(old),
			math.Float32bits(new),
		) {
			break
		}
	}
	return
}

func (v *Float32) Cas(old, new float32) (swapped bool) {
	return atomic.CompareAndSwapUint32(&v.value, math.Float32bits(old), math.Float32bits(new))
}

func (v *Float32) String() string {
	return strconv.FormatFloat(float64(v.Val()), 'g', -1, 32)
}

func (v *Float32) MarshalJSON() ([]byte, error) {
	return conv.UnsafeStrToBytes(strconv.FormatFloat(float64(v.Val()), 'g', -1, 32)), nil
}

func (v *Float32) UnmarshalJSON(b []byte) error {
	v.Set(conv.Float32(conv.UnsafeBytesToStr(b)))
	return nil
}

func (v *Float32) UnmarshalValue(value interface{}) error {
	v.Set(conv.Float32(value))
	return nil
}
