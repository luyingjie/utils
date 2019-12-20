package convert

import (
	"encoding/json"
)

func StructToMap(obj interface{}) map[string]interface{} {
	m := make(map[string]interface{})
    j, _ := json.Marshal(obj)
	json.Unmarshal(j, &m)
	return m
}

func StructToMapString(obj interface{}) map[string]string {
	rs := StructToMap(obj)

	var data = make(map[string]string)
	for key, val := range rs {
		data[key] = val.(string)
	}
	return data
}