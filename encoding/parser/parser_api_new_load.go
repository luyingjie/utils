package parser

import (
	vjson "github.com/luyingjie/utils/encoding/json"
)

// New creates a Parser object with any variable type of <data>, but <data> should be a map, struct or
// slice for data access reason, or it will make no sense.
//
// The parameter <safe> specifies whether using this Json object in concurrent-safe context, which
// is false in default.
func New(data interface{}, safe ...bool) *Parser {
	return vjson.New(data, safe...)
}

// NewWithTag creates a Parser object with any variable type of <data>, but <data> should be a map
// or slice for data access reason, or it will make no sense.
//
// The parameter <tags> specifies priority tags for struct conversion to map, multiple tags joined
// with char ','.
//
// The parameter <safe> specifies whether using this Json object in concurrent-safe context, which
// is false in default.
func NewWithTag(data interface{}, tags string, safe ...bool) *Parser {
	return vjson.NewWithTag(data, tags, safe...)
}

// Load loads content from specified file <path>,
// and creates a Parser object from its content.
func Load(path string, safe ...bool) (*Parser, error) {
	return vjson.Load(path, safe...)
}

// LoadContent creates a Parser object from given content,
// it checks the data type of <content> automatically,
// supporting JSON, XML, INI, YAML and TOML types of data.
func LoadContent(data interface{}, safe ...bool) (*Parser, error) {
	return vjson.LoadContent(data, safe...)
}

func LoadJson(data interface{}, safe ...bool) (*Parser, error) {
	return vjson.LoadJson(data, safe...)
}

func LoadXml(data interface{}, safe ...bool) (*Parser, error) {
	return vjson.LoadXml(data, safe...)
}

func LoadYaml(data interface{}, safe ...bool) (*Parser, error) {
	return vjson.LoadYaml(data, safe...)
}

func LoadToml(data interface{}, safe ...bool) (*Parser, error) {
	return vjson.LoadToml(data, safe...)
}

func LoadIni(data interface{}, safe ...bool) (*Parser, error) {
	return vjson.LoadIni(data, safe...)
}
