package vtype

import (
	"bytes"
	"sync/atomic"
	"utils/convert/conv"
)

type String struct {
	value atomic.Value
}

func NewString(value ...string) *String {
	t := &String{}
	if len(value) > 0 {
		t.value.Store(value[0])
	}
	return t
}

func (v *String) Clone() *String {
	return NewString(v.Val())
}

func (v *String) Set(value string) (old string) {
	old = v.Val()
	v.value.Store(value)
	return
}

func (v *String) Val() string {
	s := v.value.Load()
	if s != nil {
		return s.(string)
	}
	return ""
}

func (v *String) String() string {
	return v.Val()
}

func (v *String) MarshalJSON() ([]byte, error) {
	return conv.UnsafeStrToBytes(`"` + v.Val() + `"`), nil
}

func (v *String) UnmarshalJSON(b []byte) error {
	v.Set(conv.UnsafeBytesToStr(bytes.Trim(b, `"`)))
	return nil
}

func (v *String) UnmarshalValue(value interface{}) error {
	v.Set(conv.String(value))
	return nil
}
