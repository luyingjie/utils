package json

import (
	"bytes"
	"utils/util/json"

	"utils/os/conv"
)

func Valid(data interface{}) bool {
	return json.Valid(conv.Bytes(data))
}

func Encode(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func Decode(data interface{}) (interface{}, error) {
	var value interface{}
	if err := DecodeTo(conv.Bytes(data), &value); err != nil {
		return nil, err
	} else {
		return value, nil
	}
}

func DecodeTo(data interface{}, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(conv.Bytes(data)))
	return decoder.Decode(v)
}

func DecodeToJson(data interface{}, safe ...bool) (*Json, error) {
	if v, err := Decode(conv.Bytes(data)); err != nil {
		return nil, err
	} else {
		return New(v, safe...), nil
	}
}
