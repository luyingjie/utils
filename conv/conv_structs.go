package conv

import (
	"fmt"
	"reflect"
)

func Structs(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	return doStructs(params, pointer, false, mapping...)
}

func StructsDeep(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	return doStructs(params, pointer, true, mapping...)
}

func doStructs(params interface{}, pointer interface{}, deep bool, mapping ...map[string]string) (err error) {
	if params == nil {
		return fmt.Errorf("params cannot be nil")
	}
	if pointer == nil {
		return fmt.Errorf("object pointer cannot be nil")
	}
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	pointerRv, ok := pointer.(reflect.Value)
	if !ok {
		pointerRv = reflect.ValueOf(pointer)
		if kind := pointerRv.Kind(); kind != reflect.Ptr {
			return fmt.Errorf("pointer should be type of pointer, but got: %v", kind)
		}
	}
	params = Maps(params)
	var (
		reflectValue = reflect.ValueOf(params)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		// If <params> is an empty slice, no conversion.
		if reflectValue.Len() == 0 {
			return nil
		}
		var (
			array    = reflect.MakeSlice(pointerRv.Type().Elem(), reflectValue.Len(), reflectValue.Len())
			itemType = array.Index(0).Type()
		)
		for i := 0; i < reflectValue.Len(); i++ {
			if itemType.Kind() == reflect.Ptr {
				// Slice element is type pointer.
				e := reflect.New(itemType.Elem()).Elem()
				if deep {
					if err = StructDeep(reflectValue.Index(i).Interface(), e, mapping...); err != nil {
						return err
					}
				} else {
					if err = Struct(reflectValue.Index(i).Interface(), e, mapping...); err != nil {
						return err
					}
				}
				array.Index(i).Set(e.Addr())
			} else {
				// Slice element is not type of pointer.
				e := reflect.New(itemType).Elem()
				if deep {
					if err = StructDeep(reflectValue.Index(i).Interface(), e, mapping...); err != nil {
						return err
					}
				} else {
					if err = Struct(reflectValue.Index(i).Interface(), e, mapping...); err != nil {
						return err
					}
				}
				array.Index(i).Set(e)
			}
		}
		pointerRv.Elem().Set(array)
		return nil
	default:
		return fmt.Errorf("params should be type of slice, but got: %v", reflectKind)
	}
}
