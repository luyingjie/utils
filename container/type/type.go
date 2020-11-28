package vtype

type Type = Interface

func New(value ...interface{}) *Type {
	return NewInterface(value...)
}
