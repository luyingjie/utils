package util

import (
	"bytes"
	"encoding/gob"

	"github.com/luyingjie/utils/util/empty"
)

func Throw(exception interface{}) {
	panic(exception)
}

func TryCatch(try func(), catch ...func(exception interface{})) {
	defer func() {
		if e := recover(); e != nil && len(catch) > 0 {
			catch[0](e)
		}
	}()
	try()
}

func IsEmpty(value interface{}) bool {
	return empty.IsEmpty(value)
}

// MapDeepCopy Map深克隆
func MapDeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = MapDeepCopy(v)
		}

		return newMap
	} else if valueMap, ok := value.(map[string]string); ok {
		newMap := make(map[string]string)
		for k, v := range valueMap {
			newMap[k] = v
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = MapDeepCopy(v)
		}

		return newSlice
	}

	return value
}

// DeepCopy 深克隆
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
