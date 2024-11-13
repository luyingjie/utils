package vvar

import (
	"github.com/luyingjie/utils/conv"
)

func (v *Var) Struct(pointer interface{}, mapping ...map[string]string) error {
	return conv.Struct(v.Val(), pointer, mapping...)
}

func (v *Var) StructDeep(pointer interface{}, mapping ...map[string]string) error {
	return conv.StructDeep(v.Val(), pointer, mapping...)
}

func (v *Var) Structs(pointer interface{}, mapping ...map[string]string) error {
	return conv.Structs(v.Val(), pointer, mapping...)
}

func (v *Var) StructsDeep(pointer interface{}, mapping ...map[string]string) error {
	return conv.StructsDeep(v.Val(), pointer, mapping...)
}

func (v *Var) Scan(pointer interface{}, mapping ...map[string]string) error {
	return conv.Scan(v.Val(), pointer, mapping...)
}

func (v *Var) ScanDeep(pointer interface{}, mapping ...map[string]string) error {
	return conv.ScanDeep(v.Val(), pointer, mapping...)
}
