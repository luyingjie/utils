package vtype

import (
	"math"
	"strconv"
	"sync/atomic"
	"unsafe"

	"github.com/luyingjie/utils/conv"
)

type Float64 struct {
	value uint64
}

func NewFloat64(value ...float64) *Float64 {
	if len(value) > 0 {
		return &Float64{
			value: math.Float64bits(value[0]),
		}
	}
	return &Float64{}
}

func (v *Float64) Clone() *Float64 {
	return NewFloat64(v.Val())
}

func (v *Float64) Set(value float64) (old float64) {
	return math.Float64frombits(atomic.SwapUint64(&v.value, math.Float64bits(value)))
}

func (v *Float64) Val() float64 {
	return math.Float64frombits(atomic.LoadUint64(&v.value))
}

func (v *Float64) Add(delta float64) (new float64) {
	for {
		old := math.Float64frombits(v.value)
		new = old + delta
		if atomic.CompareAndSwapUint64(
			(*uint64)(unsafe.Pointer(&v.value)),
			math.Float64bits(old),
			math.Float64bits(new),
		) {
			break
		}
	}
	return
}

func (v *Float64) Cas(old, new float64) (swapped bool) {
	return atomic.CompareAndSwapUint64(&v.value, math.Float64bits(old), math.Float64bits(new))
}

func (v *Float64) String() string {
	return strconv.FormatFloat(v.Val(), 'g', -1, 64)
}

func (v *Float64) MarshalJSON() ([]byte, error) {
	return conv.UnsafeStrToBytes(strconv.FormatFloat(v.Val(), 'g', -1, 64)), nil
}

func (v *Float64) UnmarshalJSON(b []byte) error {
	v.Set(conv.Float64(conv.UnsafeBytesToStr(b)))
	return nil
}

func (v *Float64) UnmarshalValue(value interface{}) error {
	v.Set(conv.Float64(value))
	return nil
}
