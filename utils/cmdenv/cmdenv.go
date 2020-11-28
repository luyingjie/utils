package cmdenv

import (
	"os"
	"regexp"
	"strings"

	vvar "utils/container/var"
)

var (
	cmdOptions = make(map[string]string)
)

func init() {
	doInit()
}

func doInit() {
	reg := regexp.MustCompile(`\-\-{0,1}(.+?)=(.+)`)
	for i := 0; i < len(os.Args); i++ {
		result := reg.FindStringSubmatch(os.Args[i])
		if len(result) > 1 {
			cmdOptions[result[1]] = result[2]
		}
	}
}

func Get(key string, def ...interface{}) *vvar.Var {
	value := interface{}(nil)
	if len(def) > 0 {
		value = def[0]
	}
	cmdKey := strings.ToLower(strings.Replace(key, "_", ".", -1))
	if v, ok := cmdOptions[cmdKey]; ok {
		value = v
	} else {
		envKey := strings.ToUpper(strings.Replace(key, ".", "_", -1))
		if v := os.Getenv(envKey); v != "" {
			value = v
		}
	}
	return vvar.New(value)
}
