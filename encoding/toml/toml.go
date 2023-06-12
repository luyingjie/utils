package toml

import (
	"bytes"
	// "utils/util/json"
	json2 "encoding/json"

	"github.com/BurntSushi/toml"
)

func Encode(v interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if err := toml.NewEncoder(buffer).Encode(v); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func Decode(v []byte) (interface{}, error) {
	var result interface{}
	if err := toml.Unmarshal(v, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func DecodeTo(v []byte, result interface{}) error {
	return toml.Unmarshal(v, result)
}

func ToJson(v []byte) ([]byte, error) {
	if r, err := Decode(v); err != nil {
		return nil, err
	} else {
		// 这里工具包github.com/json-iterator/go格式化会出错，以前不会，可能是版本升级后出现的问题，这个问题是升级到gon1.20的时候发现的，所以之前以为是其他问题。
		// 先使用官方序列化包。
		// return json.Marshal(r)
		return json2.Marshal(r)
	}
}
