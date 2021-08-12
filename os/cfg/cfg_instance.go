package cfg

import (
	"fmt"
	vmap "utils/container/map"
)

const (
	// Default group name for instance usage.
	DEFAULT_NAME   = "config"
	DEFAULT_FORMAT = "toml"
)

var (
	// Instances map containing configuration instances.
	instances = vmap.NewStrAnyMap(true)
)

// Instance returns an instance of Config with default settings.
// The parameter <name> is the name for the instance. But very note that, if the file "name.toml"
// exists in the configuration directory, it then sets it as the default configuration file. The
// toml file type is the default configuration file type.
func Instance(name ...string) *Config {
	key := DEFAULT_NAME
	format := DEFAULT_FORMAT
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	return instances.GetOrSetFuncLock(key, func() interface{} {
		c := New()
		file := fmt.Sprintf(`%s.%s`, key, format)
		if c.Available(file) {
			c.SetFileName(file)
		}
		return c
	}).(*Config)
}

// InstanceF 自定义获取配置文件。
// 第一个参数为文件名，默认：config。
// 第二个参数为文件格式，兼顾非toml格式的读取，所以默认为json。
// 第三个以后参数为路径，可以指定多个。
func InstanceF(name ...string) *Config {
	key := DEFAULT_NAME
	format := "json"
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	if len(name) > 1 && name[1] != "" {
		format = name[1]
	}
	return instances.GetOrSetFuncLock(key, func() interface{} {
		c := New()
		file := ""
		if name[0] == "" && name[1] == "" && name[2] != "" {
			file = name[2]
		} else {
			if len(name) > 2 && name[2] != "" {
				for _, v := range name[2:] {
					c.AddPath(v)
				}
			}
			file = fmt.Sprintf(`%s.%s`, key, format)
		}

		if c.Available(file) {
			c.SetFileName(file)
		}
		return c
	}).(*Config)
}
