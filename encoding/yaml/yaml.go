package yaml

import (
	"utils/util/json"

	"gopkg.in/yaml.v3"

	"utils/os/conv"
)

func Encode(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func Decode(v []byte) (interface{}, error) {
	var result map[string]interface{}
	if err := yaml.Unmarshal(v, &result); err != nil {
		return nil, err
	}
	return conv.MapDeep(result), nil
}

func DecodeTo(v []byte, result interface{}) error {
	return yaml.Unmarshal(v, result)
}

func ToJson(v []byte) ([]byte, error) {
	if r, err := Decode(v); err != nil {
		return nil, err
	} else {
		return json.Marshal(r)
	}
}
