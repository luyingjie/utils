package conv

import (
	"reflect"
	myerror "utils/os/error"
)

func Scan(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	t := reflect.TypeOf(pointer)
	k := t.Kind()
	if k != reflect.Ptr {
		return myerror.Newf("params should be type of pointer, but got: %v", k)
	}
	switch t.Elem().Kind() {
	case reflect.Array, reflect.Slice:
		return Structs(params, pointer, mapping...)
	default:
		return Struct(params, pointer, mapping...)
	}
}

func ScanDeep(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	t := reflect.TypeOf(pointer)
	k := t.Kind()
	if k != reflect.Ptr {
		return myerror.Newf("params should be type of pointer, but got: %v", k)
	}
	switch t.Elem().Kind() {
	case reflect.Array, reflect.Slice:
		return StructsDeep(params, pointer, mapping...)
	default:
		return StructDeep(params, pointer, mapping...)
	}
}
