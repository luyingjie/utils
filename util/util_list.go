package util

import (
	"reflect"
)

func ListItemValues(list interface{}, key interface{}, subKey ...interface{}) (values []interface{}) {
	var (
		reflectValue = reflect.ValueOf(list)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	values = []interface{}{}
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		if reflectValue.Len() == 0 {
			return
		}
		for i := 0; i < reflectValue.Len(); i++ {
			if value, ok := doItemValue(reflectValue.Index(i), key); ok {
				if len(subKey) > 0 && subKey[0] != nil {
					if subValue, ok := doItemValue(value, subKey[0]); ok {
						values = append(values, subValue)
					}
				} else {
					values = append(values, value)
				}
			}
		}
	}
	return
}

func ItemValue(item interface{}, key interface{}) (value interface{}) {
	value, _ = doItemValue(item, key)
	return
}

func doItemValue(item interface{}, key interface{}) (value interface{}, found bool) {
	var reflectValue reflect.Value
	if v, ok := item.(reflect.Value); ok {
		reflectValue = v
	} else {
		reflectValue = reflect.ValueOf(item)
	}
	reflectKind := reflectValue.Kind()
	if reflectKind == reflect.Interface {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	var keyValue reflect.Value
	if v, ok := key.(reflect.Value); ok {
		keyValue = v
	} else {
		keyValue = reflect.ValueOf(key)
	}
	switch reflectKind {
	case reflect.Map:
		v := reflectValue.MapIndex(keyValue)
		if v.IsValid() {
			found = true
			value = v.Interface()
		}

	case reflect.Struct:
		v := reflectValue.FieldByName(keyValue.String())
		if v.IsValid() {
			found = true
			value = v.Interface()
		}
	}
	return
}

func ListItemValuesUnique(list interface{}, key string, subKey ...interface{}) []interface{} {
	values := ListItemValues(list, key, subKey...)
	if len(values) > 0 {
		var (
			ok bool
			m  = make(map[interface{}]struct{}, len(values))
		)
		for i := 0; i < len(values); {
			if _, ok = m[values[i]]; ok {
				values = SliceDelete(values, i)
			} else {
				m[values[i]] = struct{}{}
				i++
			}
		}
	}
	return values
}
