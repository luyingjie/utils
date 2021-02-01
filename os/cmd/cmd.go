// Package gcmd provides console operations, like options/arguments reading and command running.
package cmd

import (
	"os"

	vvar "utils/container/var"

	"utils/text/regex"
)

var (
	defaultParsedArgs     = make([]string, 0)
	defaultParsedOptions  = make(map[string]string)
	defaultCommandFuncMap = make(map[string]func())
)

// Custom initialization.
func doInit() {
	if len(defaultParsedArgs) > 0 {
		return
	}
	// Parsing os.Args with default algorithm.
	// The option should use '=' to separate its name and value in default.
	for _, arg := range os.Args {
		array, _ := regex.MatchString(`^\-{1,2}([\w\?\.\-]+)={0,1}(.*)$`, arg)
		if len(array) == 3 {
			defaultParsedOptions[array[1]] = array[2]
		} else {
			defaultParsedArgs = append(defaultParsedArgs, arg)
		}
	}
}

// GetOpt returns the option value named <name>.
func GetOpt(name string, def ...string) string {
	doInit()
	if v, ok := defaultParsedOptions[name]; ok {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// GetOptVar returns the option value named <name> as vvar.Var.
func GetOptVar(name string, def ...string) *vvar.Var {
	doInit()
	return vvar.New(GetOpt(name, def...))
}

// GetOptAll returns all parsed options.
func GetOptAll() map[string]string {
	doInit()
	return defaultParsedOptions
}

// ContainsOpt checks whether option named <name> exist in the arguments.
func ContainsOpt(name string, def ...string) bool {
	doInit()
	_, ok := defaultParsedOptions[name]
	return ok
}

// GetArg returns the argument at <index>.
func GetArg(index int, def ...string) string {
	doInit()
	if index < len(defaultParsedArgs) {
		return defaultParsedArgs[index]
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// GetArgVar returns the argument at <index> as vvar.Var.
func GetArgVar(index int, def ...string) *vvar.Var {
	doInit()
	return vvar.New(GetArg(index, def...))
}

// GetArgAll returns all parsed arguments.
func GetArgAll() []string {
	doInit()
	return defaultParsedArgs
}

// BuildOptions builds the options as string.
func BuildOptions(m map[string]string, prefix ...string) string {
	options := ""
	leadStr := "-"
	if len(prefix) > 0 {
		leadStr = prefix[0]
	}
	for k, v := range m {
		if len(options) > 0 {
			options += " "
		}
		options += leadStr + k
		if v != "" {
			options += "=" + v
		}
	}
	return options
}
