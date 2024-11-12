package redis

import (
	"fmt"
	"time"

	"github.com/luyingjie/utils/pkg/container/vmap"
	"github.com/luyingjie/utils/pkg/conv"
	"github.com/luyingjie/utils/pkg/text/regex"
)

const (
	DEFAULT_GROUP_NAME = "default" // Default configuration group name.
	DEFAULT_REDIS_PORT = 6379      // Default redis port configuration if not passed.
)

var (
	// Configuration groups.
	configs = vmap.NewStrAnyMap(true)
)

// SetConfig sets the global configuration for specified group.
// If <name> is not passed, it sets configuration for the default group name.
func SetConfig(config Config, name ...string) {
	group := DEFAULT_GROUP_NAME
	if len(name) > 0 {
		group = name[0]
	}
	configs.Set(group, config)
	instances.Remove(group)
	// 先只输出到命令行
	fmt.Printf(`SetConfig for group "%s": %+v`, group, config)
}

// SetConfigByStr sets the global configuration for specified group with string.
// If <name> is not passed, it sets configuration for the default group name.
func SetConfigByStr(str string, name ...string) error {
	group := DEFAULT_GROUP_NAME
	if len(name) > 0 {
		group = name[0]
	}
	config, err := ConfigFromStr(str)
	if err != nil {
		return err
	}
	configs.Set(group, config)
	instances.Remove(group)
	return nil
}

// GetConfig returns the global configuration with specified group name.
// If <name> is not passed, it returns configuration of the default group name.
func GetConfig(name ...string) (config Config, ok bool) {
	group := DEFAULT_GROUP_NAME
	if len(name) > 0 {
		group = name[0]
	}
	if v := configs.Get(group); v != nil {
		return v.(Config), true
	}
	return Config{}, false
}

// RemoveConfig removes the global configuration with specified group.
// If <name> is not passed, it removes configuration of the default group name.
func RemoveConfig(name ...string) {
	group := DEFAULT_GROUP_NAME
	if len(name) > 0 {
		group = name[0]
	}
	configs.Remove(group)
	instances.Remove(group)

	fmt.Printf(`RemoveConfig: %s`, group)
}

// ConfigFromStr parses and returns config from given str.
// Eg: host:port[,db,pass?maxIdle=x&maxActive=x&idleTimeout=x&maxConnLifetime=x]
func ConfigFromStr(str string) (config Config, err error) {
	array, _ := regex.MatchString(`([^:]+):*(\d*),{0,1}(\d*),{0,1}(.*)\?(.+)`, str)
	if len(array) == 6 {
		parse, _ := str.Parse(array[5])
		config = Config{
			Host: array[1],
			Port: conv.Int(array[2]),
			Db:   conv.Int(array[3]),
			Pass: array[4],
		}
		if config.Port == 0 {
			config.Port = DEFAULT_REDIS_PORT
		}
		if v, ok := parse["maxIdle"]; ok {
			config.MaxIdle = conv.Int(v)
		}
		if v, ok := parse["maxActive"]; ok {
			config.MaxActive = conv.Int(v)
		}
		if v, ok := parse["idleTimeout"]; ok {
			config.IdleTimeout = conv.Duration(v) * time.Second
		}
		if v, ok := parse["maxConnLifetime"]; ok {
			config.MaxConnLifetime = conv.Duration(v) * time.Second
		}
		if v, ok := parse["tls"]; ok {
			config.TLS = conv.Bool(v)
		}
		if v, ok := parse["skipVerify"]; ok {
			config.TLSSkipVerify = conv.Bool(v)
		}
		return
	}
	array, _ = regex.MatchString(`([^:]+):*(\d*),{0,1}(\d*),{0,1}(.*)`, str)
	if len(array) == 5 {
		config = Config{
			Host: array[1],
			Port: conv.Int(array[2]),
			Db:   conv.Int(array[3]),
			Pass: array[4],
		}
		if config.Port == 0 {
			config.Port = DEFAULT_REDIS_PORT
		}
	} else {
		err = fmt.Errorf(`invalid redis configuration: "%s"`, str)
	}
	return
}

// ClearConfig removes all configurations and instances of redis.
func ClearConfig() {
	configs.Clear()
	instances.Clear()
}
