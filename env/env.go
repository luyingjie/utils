// Package genv provides operations for environment variables of system.
package env

import (
	"os"
	"strings"

	"github.com/luyingjie/utils/container/vvar"
)

// All returns a copy of strings representing the environment,
// in the form "key=value".
func All() []string {
	return os.Environ()
}

// Map returns a copy of strings representing the environment as a map.
func Map() map[string]string {
	m := make(map[string]string)
	i := 0
	for _, s := range os.Environ() {
		i = strings.IndexByte(s, '=')
		m[s[0:i]] = s[i+1:]
	}
	return m
}

// Get returns the value of the environment variable named by the <key>.
// It returns given <def> if the variable does not exist in the environment.
func Get(key string, def ...string) string {
	v, ok := os.LookupEnv(key)
	if !ok && len(def) > 0 {
		return def[0]
	}
	return v
}

// GetVar creates and returns a Var with the value of the environment variable
// named by the <key>. It uses the given <def> if the variable does not exist
// in the environment.
func GetVar(key string, def ...interface{}) *vvar.Var {
	v, ok := os.LookupEnv(key)
	if !ok && len(def) > 0 {
		return vvar.New(def[0])
	}
	return vvar.New(v)
}

// Set sets the value of the environment variable named by the <key>.
// It returns an error, if any.
func Set(key, value string) error {
	return os.Setenv(key, value)
}

// Contains checks whether the environment variable named <key> exists.
func Contains(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

// Build builds a map to a environment variable slice.
func Build(m map[string]string) []string {
	array := make([]string, len(m))
	index := 0
	for k, v := range m {
		array[index] = k + "=" + v
		index++
	}
	return array
}

// Remove deletes one or more environment variables.
func Remove(key ...string) error {
	var err error
	for _, v := range key {
		err = os.Unsetenv(v)
		if err != nil {
			return err
		}
	}
	return nil
}
