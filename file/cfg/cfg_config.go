package cfg

import (
	"github.com/luyingjie/utils/container/vmap"
)

var (
	// Customized configuration content.
	configs = vmap.NewStrStrMap(true)
)

// SetContent sets customized configuration content for specified <file>.
// The <file> is unnecessary param, default is DEFAULT_CONFIG_FILE.
func SetContent(content string, file ...string) {
	name := DEFAULT_CONFIG_FILE
	if len(file) > 0 {
		name = file[0]
	}
	// Clear file cache for instances which cached <name>.
	instances.LockFunc(func(m map[string]interface{}) {
		if configs.Contains(name) {
			for _, v := range m {
				v.(*Config).jsons.Remove(name)
			}
		}
		configs.Set(name, content)
	})
}

// GetContent returns customized configuration content for specified <file>.
// The <file> is unnecessary param, default is DEFAULT_CONFIG_FILE.
func GetContent(file ...string) string {
	name := DEFAULT_CONFIG_FILE
	if len(file) > 0 {
		name = file[0]
	}
	return configs.Get(name)
}

// RemoveContent removes the global configuration with specified <file>.
// If <name> is not passed, it removes configuration of the default group name.
func RemoveContent(file ...string) {
	name := DEFAULT_CONFIG_FILE
	if len(file) > 0 {
		name = file[0]
	}
	// Clear file cache for instances which cached <name>.
	instances.LockFunc(func(m map[string]interface{}) {
		if configs.Contains(name) {
			for _, v := range m {
				v.(*Config).jsons.Remove(name)
			}
			configs.Remove(name)
		}
	})

	// intlog.Printf(`RemoveContent: %s`, name)
}

// ClearContent removes all global configuration contents.
func ClearContent() {
	configs.Clear()
	// Clear cache for all instances.
	instances.LockFunc(func(m map[string]interface{}) {
		for _, v := range m {
			v.(*Config).jsons.Clear()
		}
	})

	// intlog.Print(`RemoveConfig`)
}
