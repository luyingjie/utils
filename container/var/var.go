package vvar

import (
	"time"
	vtype "utils/container/type"
	"utils/util/empty"
	"utils/util/json"

	"utils/conv"
	vtime "utils/os/time"
)

type Var struct {
	value interface{}
	safe  bool
}

func New(value interface{}, safe ...bool) *Var {
	v := Var{}
	if len(safe) > 0 && !safe[0] {
		v.safe = true
		v.value = vtype.NewInterface(value)
	} else {
		v.value = value
	}
	return &v
}

func Create(value interface{}, safe ...bool) Var {
	v := Var{}
	if len(safe) > 0 && !safe[0] {
		v.safe = true
		v.value = vtype.NewInterface(value)
	} else {
		v.value = value
	}
	return v
}

func (v *Var) Clone() *Var {
	return New(v.Val(), v.safe)
}

func (v *Var) Set(value interface{}) (old interface{}) {
	if v.safe {
		if t, ok := v.value.(*vtype.Interface); ok {
			old = t.Set(value)
			return
		}
	}
	old = v.value
	v.value = value
	return
}

func (v *Var) Val() interface{} {
	if v == nil {
		return nil
	}
	if v.safe {
		if t, ok := v.value.(*vtype.Interface); ok {
			return t.Val()
		}
	}
	return v.value
}

func (v *Var) Interface() interface{} {
	return v.Val()
}

func (v *Var) IsNil() bool {
	return v.Val() == nil
}

func (v *Var) IsEmpty() bool {
	return empty.IsEmpty(v.Val())
}

func (v *Var) Bytes() []byte {
	return conv.Bytes(v.Val())
}

func (v *Var) String() string {
	return conv.String(v.Val())
}

func (v *Var) Bool() bool {
	return conv.Bool(v.Val())
}

func (v *Var) Int() int {
	return conv.Int(v.Val())
}

func (v *Var) Ints() []int {
	return conv.Ints(v.Val())
}

func (v *Var) Int8() int8 {
	return conv.Int8(v.Val())
}

func (v *Var) Int16() int16 {
	return conv.Int16(v.Val())
}

func (v *Var) Int32() int32 {
	return conv.Int32(v.Val())
}

func (v *Var) Int64() int64 {
	return conv.Int64(v.Val())
}

func (v *Var) Uint() uint {
	return conv.Uint(v.Val())
}

func (v *Var) Uints() []uint {
	return conv.Uints(v.Val())
}

func (v *Var) Uint8() uint8 {
	return conv.Uint8(v.Val())
}

func (v *Var) Uint16() uint16 {
	return conv.Uint16(v.Val())
}

func (v *Var) Uint32() uint32 {
	return conv.Uint32(v.Val())
}

func (v *Var) Uint64() uint64 {
	return conv.Uint64(v.Val())
}

func (v *Var) Float32() float32 {
	return conv.Float32(v.Val())
}

func (v *Var) Float64() float64 {
	return conv.Float64(v.Val())
}

func (v *Var) Floats() []float64 {
	return conv.Floats(v.Val())
}

func (v *Var) Strings() []string {
	return conv.Strings(v.Val())
}

func (v *Var) Interfaces() []interface{} {
	return conv.Interfaces(v.Val())
}

func (v *Var) Slice() []interface{} {
	return v.Interfaces()
}

func (v *Var) Array() []interface{} {
	return v.Interfaces()
}

func (v *Var) Vars() []*Var {
	array := conv.Interfaces(v.Val())
	if len(array) == 0 {
		return nil
	}
	vars := make([]*Var, len(array))
	for k, v := range array {
		vars[k] = New(v)
	}
	return vars
}

func (v *Var) Time(format ...string) time.Time {
	return conv.Time(v.Val(), format...)
}

func (v *Var) Duration() time.Duration {
	return conv.Duration(v.Val())
}

func (v *Var) VTime(format ...string) *vtime.Time {
	return conv.VTime(v.Val(), format...)
}

func (v *Var) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Val())
}

func (v *Var) UnmarshalJSON(b []byte) error {
	var i interface{}
	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	v.Set(i)
	return nil
}

func (v *Var) UnmarshalValue(value interface{}) error {
	v.Set(value)
	return nil
}
