package mytype

import (
	"sync/atomic"
	"utils/utils/json"

	"utils/convert/conv"
)

type Interface struct {
	value atomic.Value
}

func NewInterface(value ...interface{}) *Interface {
	t := &Interface{}
	if len(value) > 0 && value[0] != nil {
		t.value.Store(value[0])
	}
	return t
}

func (v *Interface) Clone() *Interface {
	return NewInterface(v.Val())
}

func (v *Interface) Set(value interface{}) (old interface{}) {
	old = v.Val()
	v.value.Store(value)
	return
}

func (v *Interface) Val() interface{} {
	return v.value.Load()
}

func (v *Interface) String() string {
	return conv.String(v.Val())
}

func (v *Interface) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Val())
}

func (v *Interface) UnmarshalJSON(b []byte) error {
	var i interface{}
	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	v.Set(i)
	return nil
}

func (v *Interface) UnmarshalValue(value interface{}) error {
	v.Set(value)
	return nil
}
