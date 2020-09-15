package json

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"utils/utils/json"

	"utils/utils/rwmutex"

	"utils/text/regex"

	"utils/convert/conv"

	"utils/os/file"

	"utils/encoding/ini"
	"utils/encoding/toml"
	"utils/encoding/xml"
	"utils/encoding/yaml"
)

func New(data interface{}, safe ...bool) *Json {
	return NewWithTag(data, "json", safe...)
}

func NewWithTag(data interface{}, tags string, safe ...bool) *Json {
	j := (*Json)(nil)
	switch data.(type) {
	case string, []byte:
		if r, err := LoadContent(conv.Bytes(data)); err == nil {
			j = r
		} else {
			j = &Json{
				p:  &data,
				c:  byte(gDEFAULT_SPLIT_CHAR),
				vc: false,
			}
		}
	default:
		rv := reflect.ValueOf(data)
		kind := rv.Kind()
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		switch kind {
		case reflect.Slice, reflect.Array:
			i := interface{}(nil)
			i = conv.Interfaces(data)
			j = &Json{
				p:  &i,
				c:  byte(gDEFAULT_SPLIT_CHAR),
				vc: false,
			}
		case reflect.Map, reflect.Struct:
			i := interface{}(nil)
			i = conv.Map(data, tags)
			j = &Json{
				p:  &i,
				c:  byte(gDEFAULT_SPLIT_CHAR),
				vc: false,
			}
		default:
			j = &Json{
				p:  &data,
				c:  byte(gDEFAULT_SPLIT_CHAR),
				vc: false,
			}
		}
	}
	j.mu = rwmutex.New(safe...)
	return j
}

func Load(path string, safe ...bool) (*Json, error) {
	if p, err := file.Search(path); err != nil {
		return nil, err
	} else {
		path = p
	}
	return doLoadContent(file.Ext(path), file.GetBytesWithCache(path), safe...)
}

func LoadJson(data interface{}, safe ...bool) (*Json, error) {
	return doLoadContent("json", conv.Bytes(data), safe...)
}

func LoadXml(data interface{}, safe ...bool) (*Json, error) {
	return doLoadContent("xml", conv.Bytes(data), safe...)
}

func LoadIni(data interface{}, safe ...bool) (*Json, error) {
	return doLoadContent("ini", conv.Bytes(data), safe...)
}

func LoadYaml(data interface{}, safe ...bool) (*Json, error) {
	return doLoadContent("yaml", conv.Bytes(data), safe...)
}

func LoadToml(data interface{}, safe ...bool) (*Json, error) {
	return doLoadContent("toml", conv.Bytes(data), safe...)
}

func doLoadContent(dataType string, data []byte, safe ...bool) (*Json, error) {
	var err error
	var result interface{}
	if len(data) == 0 {
		return New(nil, safe...), nil
	}
	if dataType == "" {
		dataType = checkDataType(data)
	}
	switch dataType {
	case "json", ".json", ".js":

	case "xml", ".xml":
		if data, err = xml.ToJson(data); err != nil {
			return nil, err
		}

	case "yml", "yaml", ".yml", ".yaml":
		if data, err = yaml.ToJson(data); err != nil {
			return nil, err
		}

	case "toml", ".toml":
		if data, err = toml.ToJson(data); err != nil {
			return nil, err
		}
	case "ini", ".ini":
		if data, err = ini.ToJson(data); err != nil {
			return nil, err
		}
	default:
		err = errors.New("unsupported type for loading")
	}
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))

	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	switch result.(type) {
	case string, []byte:
		return nil, fmt.Errorf(`json decoding failed for content: %s`, string(data))
	}
	return New(result, safe...), nil
}

func LoadContent(data interface{}, safe ...bool) (*Json, error) {
	content := conv.Bytes(data)
	if len(content) == 0 {
		return New(nil, safe...), nil
	}
	return LoadContentType(checkDataType(content), content, safe...)
}

func LoadContentType(dataType string, data interface{}, safe ...bool) (*Json, error) {
	content := conv.Bytes(data)
	if len(content) == 0 {
		return New(nil, safe...), nil
	}
	//ignore UTF8-BOM
	if content[0] == 0xEF && content[1] == 0xBB && content[2] == 0xBF {
		content = content[3:]
	}
	return doLoadContent(dataType, content, safe...)
}

func checkDataType(content []byte) string {
	if json.Valid(content) {
		return "json"
	} else if regex.IsMatch(`^<.+>[\S\s]+<.+>$`, content) {
		return "xml"
	} else if (regex.IsMatch(`^[\n\r]*[\w\-\s\t]+\s*:\s*".+"`, content) || regex.IsMatch(`^[\n\r]*[\w\-\s\t]+\s*:\s*\w+`, content)) ||
		(regex.IsMatch(`[\n\r]+[\w\-\s\t]+\s*:\s*".+"`, content) || regex.IsMatch(`[\n\r]+[\w\-\s\t]+\s*:\s*\w+`, content)) {
		return "yml"
	} else if !regex.IsMatch(`^[\s\t\n\r]*;.+`, content) &&
		!regex.IsMatch(`[\s\t\n\r]+;.+`, content) &&
		!regex.IsMatch(`[\n\r]+[\s\t\w\-]+\.[\s\t\w\-]+\s*=\s*.+`, content) &&
		(regex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*".+"`, content) || regex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`, content)) {
		return "toml"
	} else if regex.IsMatch(`\[[\w\.]+\]`, content) &&
		(regex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*".+"`, content) || regex.IsMatch(`[\n\r]*[\s\t\w\-\."]+\s*=\s*\w+`, content)) {
		return "ini"
	} else {
		return ""
	}
}
