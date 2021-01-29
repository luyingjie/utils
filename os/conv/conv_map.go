package conv

import (
	"reflect"
	"strings"
	myerror "utils/os/error"
	"utils/util"
	"utils/util/empty"
	"utils/util/json"
)

func Map(value interface{}, tags ...string) map[string]interface{} {
	return doMapConvert(value, false, tags...)
}

func MapDeep(value interface{}, tags ...string) map[string]interface{} {
	return doMapConvert(value, true, tags...)
}

func doMapConvert(value interface{}, recursive bool, tags ...string) map[string]interface{} {
	if value == nil {
		return nil
	}

	dataMap := make(map[string]interface{})
	switch r := value.(type) {
	case string:
		if len(r) > 0 && r[0] == '{' && r[len(r)-1] == '}' {
			if err := json.Unmarshal([]byte(r), &dataMap); err != nil {
				return nil
			}
		} else {
			return nil
		}
	case []byte:
		if len(r) > 0 && r[0] == '{' && r[len(r)-1] == '}' {
			if err := json.Unmarshal(r, &dataMap); err != nil {
				return nil
			}
		} else {
			return nil
		}
	case map[interface{}]interface{}:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[interface{}]string:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[interface{}]int:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[interface{}]uint:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[interface{}]float32:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[interface{}]float64:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[string]bool:
		for k, v := range r {
			dataMap[k] = v
		}
	case map[string]int:
		for k, v := range r {
			dataMap[k] = v
		}
	case map[string]uint:
		for k, v := range r {
			dataMap[k] = v
		}
	case map[string]float32:
		for k, v := range r {
			dataMap[k] = v
		}
	case map[string]float64:
		for k, v := range r {
			dataMap[k] = v
		}
	case map[string]interface{}:
		return r
	case map[int]interface{}:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[int]string:
		for k, v := range r {
			dataMap[String(k)] = v
		}
	case map[uint]string:
		for k, v := range r {
			dataMap[String(k)] = v
		}

	default:
		// Not a common type, it then uses reflection for conversion.
		var rv reflect.Value
		if v, ok := value.(reflect.Value); ok {
			rv = v
		} else {
			rv = reflect.ValueOf(value)
		}
		kind := rv.Kind()
		// If it is a pointer, we should find its real data type.
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		switch kind {
		// If <value> is type of array, it converts the value of even number index as its key and
		// the value of odd number index as its corresponding value, for example:
		// []string{"k1","v1","k2","v2"} => map[string]interface{}{"k1":"v1", "k2":"v2"}
		// []string{"k1","v1","k2"}      => map[string]interface{}{"k1":"v1", "k2":nil}
		case reflect.Slice, reflect.Array:
			length := rv.Len()
			for i := 0; i < length; i += 2 {
				if i+1 < length {
					dataMap[String(rv.Index(i).Interface())] = rv.Index(i + 1).Interface()
				} else {
					dataMap[String(rv.Index(i).Interface())] = nil
				}
			}
		case reflect.Map:
			ks := rv.MapKeys()
			for _, k := range ks {
				dataMap[String(k.Interface())] = rv.MapIndex(k).Interface()
			}
		case reflect.Struct:
			// Map converting interface check.
			if v, ok := value.(apiMapStrAny); ok {
				return v.MapStrAny()
			}
			// Using reflect for converting.
			var (
				rtField  reflect.StructField
				rvField  reflect.Value
				rt       = rv.Type()
				name     = ""
				tagArray = StructTagPriority
			)
			switch len(tags) {
			case 0:
				// No need handle.
			case 1:
				tagArray = append(strings.Split(tags[0], ","), StructTagPriority...)
			default:
				tagArray = append(tags, StructTagPriority...)
			}
			for i := 0; i < rv.NumField(); i++ {
				rtField = rt.Field(i)
				rvField = rv.Field(i)
				// Only convert the public attributes.
				fieldName := rtField.Name
				if !util.IsLetterUpper(fieldName[0]) {
					continue
				}
				name = ""
				fieldTag := rtField.Tag
				for _, tag := range tagArray {
					if name = fieldTag.Get(tag); name != "" {
						break
					}
				}
				if name == "" {
					name = fieldName
				} else {
					// Support json tag feature: -, omitempty
					name = strings.TrimSpace(name)
					if name == "-" {
						continue
					}
					array := strings.Split(name, ",")
					if len(array) > 1 {
						switch strings.TrimSpace(array[1]) {
						case "omitempty":
							if empty.IsEmpty(rvField.Interface()) {
								continue
							} else {
								name = strings.TrimSpace(array[0])
							}
						default:
							name = strings.TrimSpace(array[0])
						}
					}
				}
				if recursive {
					var (
						rvAttrField = rvField
						rvAttrKind  = rvField.Kind()
					)
					if rvAttrKind == reflect.Ptr {
						rvAttrField = rvField.Elem()
						rvAttrKind = rvAttrField.Kind()
					}
					if rvAttrKind == reflect.Struct {
						var (
							hasNoTag        = name == fieldName
							rvAttrInterface = rvAttrField.Interface()
						)
						if hasNoTag && rtField.Anonymous {
							// It means this attribute field has no tag.
							// Overwrite the attribute with sub-struct attribute fields.
							for k, v := range doMapConvert(rvAttrInterface, recursive, tags...) {
								dataMap[k] = v
							}
						} else {
							// It means this attribute field has desired tag.
							if m := doMapConvert(rvAttrInterface, recursive, tags...); len(m) > 0 {
								dataMap[name] = m
							} else {
								dataMap[name] = rv.Field(i).Interface()
							}
						}
					} else {
						if rvField.IsValid() {
							dataMap[name] = rv.Field(i).Interface()
						} else {
							dataMap[name] = nil
						}
					}
				} else {
					if rvField.IsValid() {
						dataMap[name] = rv.Field(i).Interface()
					} else {
						dataMap[name] = nil
					}
				}
			}
		default:
			return nil
		}
	}
	return dataMap
}

func MapStrStr(value interface{}, tags ...string) map[string]string {
	if r, ok := value.(map[string]string); ok {
		return r
	}
	m := Map(value, tags...)
	if len(m) > 0 {
		vMap := make(map[string]string, len(m))
		for k, v := range m {
			vMap[k] = String(v)
		}
		return vMap
	}
	return nil
}

func MapStrStrDeep(value interface{}, tags ...string) map[string]string {
	if r, ok := value.(map[string]string); ok {
		return r
	}
	m := MapDeep(value, tags...)
	if len(m) > 0 {
		vMap := make(map[string]string, len(m))
		for k, v := range m {
			vMap[k] = String(v)
		}
		return vMap
	}
	return nil
}

func MapToMap(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doMapToMap(params, pointer, false, mapping...)
}

func MapToMapDeep(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doMapToMap(params, pointer, true, mapping...)
}

func doMapToMap(params interface{}, pointer interface{}, deep bool, mapping ...map[string]string) (err error) {
	var (
		paramsRv   = reflect.ValueOf(params)
		paramsKind = paramsRv.Kind()
	)
	if paramsKind == reflect.Ptr {
		paramsRv = paramsRv.Elem()
		paramsKind = paramsRv.Kind()
	}
	if paramsKind != reflect.Map {
		return myerror.New("params should be type of map")
	}
	// Empty params map, no need continue.
	if paramsRv.Len() == 0 {
		return nil
	}
	var pointerRv reflect.Value
	if v, ok := pointer.(reflect.Value); ok {
		pointerRv = v
	} else {
		pointerRv = reflect.ValueOf(pointer)
	}
	pointerKind := pointerRv.Kind()
	for pointerKind == reflect.Ptr {
		pointerRv = pointerRv.Elem()
		pointerKind = pointerRv.Kind()
	}
	if pointerKind != reflect.Map {
		return myerror.New("pointer should be type of *map")
	}
	defer func() {
		if e := recover(); e != nil {
			err = myerror.Newf("%v", e)
		}
	}()
	var (
		paramsKeys       = paramsRv.MapKeys()
		pointerKeyType   = pointerRv.Type().Key()
		pointerValueType = pointerRv.Type().Elem()
		pointerValueKind = pointerValueType.Kind()
		dataMap          = reflect.MakeMapWithSize(pointerRv.Type(), len(paramsKeys))
	)
	// Retrieve the true element type of target map.
	if pointerValueKind == reflect.Ptr {
		pointerValueKind = pointerValueType.Elem().Kind()
	}
	for _, key := range paramsKeys {
		e := reflect.New(pointerValueType).Elem()
		switch pointerValueKind {
		case reflect.Map, reflect.Struct:
			if deep {
				if err = StructDeep(paramsRv.MapIndex(key).Interface(), e, mapping...); err != nil {
					return err
				}
			} else {
				if err = Struct(paramsRv.MapIndex(key).Interface(), e, mapping...); err != nil {
					return err
				}
			}
		default:
			e.Set(
				reflect.ValueOf(
					Convert(
						paramsRv.MapIndex(key).Interface(),
						pointerValueType.String(),
					),
				),
			)
		}
		dataMap.SetMapIndex(
			reflect.ValueOf(
				Convert(
					key.Interface(),
					pointerKeyType.Name(),
				),
			),
			e,
		)
	}
	pointerRv.Set(dataMap)
	return nil
}

func MapToMaps(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doMapToMaps(params, pointer, false, mapping...)
}

func MapToMapsDeep(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	return doMapToMaps(params, pointer, true, mapping...)
}

func doMapToMaps(params interface{}, pointer interface{}, deep bool, mapping ...map[string]string) (err error) {
	var (
		paramsRv   = reflect.ValueOf(params)
		paramsKind = paramsRv.Kind()
	)
	if paramsKind == reflect.Ptr {
		paramsRv = paramsRv.Elem()
		paramsKind = paramsRv.Kind()
	}
	if paramsKind != reflect.Map {
		return myerror.New("params should be type of map")
	}
	if paramsRv.Len() == 0 {
		return nil
	}
	var (
		pointerRv   = reflect.ValueOf(pointer)
		pointerKind = pointerRv.Kind()
	)
	for pointerKind == reflect.Ptr {
		pointerRv = pointerRv.Elem()
		pointerKind = pointerRv.Kind()
	}
	if pointerKind != reflect.Map {
		return myerror.New("pointer should be type of *map/**map")
	}
	defer func() {
		// Catch the panic, especially the reflect operation panics.
		if e := recover(); e != nil {
			err = myerror.Newf("%v", e)
		}
	}()
	var (
		paramsKeys       = paramsRv.MapKeys()
		pointerKeyType   = pointerRv.Type().Key()
		pointerValueType = pointerRv.Type().Elem()
		dataMap          = reflect.MakeMapWithSize(pointerRv.Type(), len(paramsKeys))
	)
	for _, key := range paramsKeys {
		e := reflect.New(pointerValueType).Elem()
		if deep {
			if err = StructsDeep(paramsRv.MapIndex(key).Interface(), e.Addr(), mapping...); err != nil {
				return err
			}
		} else {
			if err = Structs(paramsRv.MapIndex(key).Interface(), e.Addr(), mapping...); err != nil {
				return err
			}
		}
		dataMap.SetMapIndex(
			reflect.ValueOf(
				Convert(
					key.Interface(),
					pointerKeyType.Name(),
				),
			),
			e,
		)
	}
	pointerRv.Set(dataMap)
	return nil
}
