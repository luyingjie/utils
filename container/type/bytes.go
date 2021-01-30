package vtype

import (
	"bytes"
	"encoding/base64"
	"sync/atomic"
	"utils/os/conv"
)

type Bytes struct {
	value atomic.Value
}

func NewBytes(value ...[]byte) *Bytes {
	t := &Bytes{}
	if len(value) > 0 {
		t.value.Store(value[0])
	}
	return t
}

func (v *Bytes) Clone() *Bytes {
	return NewBytes(v.Val())
}

func (v *Bytes) Set(value []byte) (old []byte) {
	old = v.Val()
	v.value.Store(value)
	return
}

func (v *Bytes) Val() []byte {
	if s := v.value.Load(); s != nil {
		return s.([]byte)
	}
	return nil
}

func (v *Bytes) String() string {
	return string(v.Val())
}

func (v *Bytes) MarshalJSON() ([]byte, error) {
	val := v.Val()
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(val)))
	base64.StdEncoding.Encode(dst, val)
	return conv.UnsafeStrToBytes(`"` + conv.UnsafeBytesToStr(dst) + `"`), nil
}

func (v *Bytes) UnmarshalJSON(b []byte) error {
	src := make([]byte, base64.StdEncoding.DecodedLen(len(b)))
	n, err := base64.StdEncoding.Decode(src, bytes.Trim(b, `"`))
	if err != nil {
		return nil
	}
	v.Set(src[:n])
	return nil
}

func (v *Bytes) UnmarshalValue(value interface{}) error {
	v.Set(conv.Bytes(value))
	return nil
}
