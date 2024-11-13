// Package gparser provides convenient API for accessing/converting variable and JSON/XML/YAML/TOML.
package parser

import (
	vjson "github.com/luyingjie/utils/encoding/json"
)

// Parser is actually alias of json.Json.
type Parser = vjson.Json
