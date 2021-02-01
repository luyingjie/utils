package cfg

import (
	"fmt"
	vmap "utils/container/map"
)

const (
	// Default group name for instance usage.
	DEFAULT_NAME = "default"
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
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	return instances.GetOrSetFuncLock(key, func() interface{} {
		c := New()
		file := fmt.Sprintf(`%s.toml`, key)
		if c.Available(file) {
			c.SetFileName(file)
		}
		return c
	}).(*Config)
}

// InstanceF 类似Instance，逻辑一样，多出一个可以指定文件格式的第二参数，用于非toml格式的读取，比如json。
func InstanceF(name ...string) *Config {
	key := DEFAULT_NAME
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	f := "toml"
	if len(name) > 1 && name[1] != "" {
		f = name[1]
	}
	return instances.GetOrSetFuncLock(key, func() interface{} {
		c := New()
		file := fmt.Sprintf(`%s.`+f, key)
		if c.Available(file) {
			c.SetFileName(file)
		}
		return c
	}).(*Config)
}
