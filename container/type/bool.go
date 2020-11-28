package vtype

import (
	"bytes"
	"sync/atomic"

	"utils/convert/conv"
)

type Bool struct {
	value int32
}

var (
	bytesTrue  = []byte("true")
	bytesFalse = []byte("false")
)

func NewBool(value ...bool) *Bool {
	t := &Bool{}
	if len(value) > 0 {
		if value[0] {
			t.value = 1
		} else {
			t.value = 0
		}
	}
	return t
}

func (v *Bool) Clone() *Bool {
	return NewBool(v.Val())
}

func (v *Bool) Set(value bool) (old bool) {
	if value {
		old = atomic.SwapInt32(&v.value, 1) == 1
	} else {
		old = atomic.SwapInt32(&v.value, 0) == 1
	}
	return
}

func (v *Bool) Val() bool {
	return atomic.LoadInt32(&v.value) > 0
}

func (v *Bool) Cas(old, new bool) (swapped bool) {
	var oldInt32, newInt32 int32
	if old {
		oldInt32 = 1
	}
	if new {
		newInt32 = 1
	}
	return atomic.CompareAndSwapInt32(&v.value, oldInt32, newInt32)
}

func (v *Bool) String() string {
	if v.Val() {
		return "true"
	}
	return "false"
}

func (v *Bool) MarshalJSON() ([]byte, error) {
	if v.Val() {
		return bytesTrue, nil
	} else {
		return bytesFalse, nil
	}
}

func (v *Bool) UnmarshalJSON(b []byte) error {
	v.Set(conv.Bool(bytes.Trim(b, `"`)))
	return nil
}

func (v *Bool) UnmarshalValue(value interface{}) error {
	v.Set(conv.Bool(value))
	return nil
}
