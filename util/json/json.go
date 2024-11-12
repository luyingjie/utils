// Json的基础操作

package json

import (
	json2 "encoding/json"
	"io"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func MarshalToString(v interface{}) (string, error) {
	return json.MarshalToString(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json2.MarshalIndent(v, prefix, indent)
}

func UnmarshalFromString(str string, v interface{}) error {
	return json.UnmarshalFromString(str, v)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func NewEncoder(writer io.Writer) *json2.Encoder {
	return json2.NewEncoder(writer)
}

func NewDecoder(reader io.Reader) *json2.Decoder {
	return json2.NewDecoder(reader)
}

func Valid(data []byte) bool {
	return json2.Valid(data)
}
