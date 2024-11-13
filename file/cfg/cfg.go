// Package gcfg provides reading, caching and managing for configuration.
package cfg

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/luyingjie/utils/text/str"

	"github.com/luyingjie/utils/container/varray"
	"github.com/luyingjie/utils/container/vmap"
	"github.com/luyingjie/utils/encoding/json"
	vfile "github.com/luyingjie/utils/file/file"
	"github.com/luyingjie/utils/file/fsnotify"
	"github.com/luyingjie/utils/file/res"
	"github.com/luyingjie/utils/file/spath"
	"github.com/luyingjie/utils/util/cmdenv"
)

const (
	DEFAULT_CONFIG_FILE = "config.toml" // The default configuration file name.
	CMDENV_KEY          = "utils.cfg"   // Configuration key for command argument or environment.
)

// Configuration struct.
type Config struct {
	name  string           // Default configuration file name.
	paths *varray.StrArray // Searching path array.
	jsons *vmap.StrAnyMap  // The pared JSON objects for configuration files.
	vc    bool             // Whether do violence check in value index searching. It affects the performance when set true(false in default).
}

var (
	resourceTryFiles = []string{"", "/", "config/", "config", "/config", "/config/"}
)

// New returns a new configuration management object.
// The parameter <file> specifies the default configuration file name for reading.
func New(file ...string) *Config {
	name := DEFAULT_CONFIG_FILE
	if len(file) > 0 {
		name = file[0]
	}
	c := &Config{
		name:  name,
		paths: varray.NewStrArray(true),
		jsons: vmap.NewStrAnyMap(true),
	}
	// Customized dir path from env/cmd.
	if envPath := cmdenv.Get(fmt.Sprintf("%s.path", CMDENV_KEY)).String(); envPath != "" {
		if vfile.Exists(envPath) {
			_ = c.SetPath(envPath)
		} else {
			if errorPrint() {
				fmt.Printf("Configuration directory path does not exist: %s", envPath)
			}
		}
	} else {
		// Dir path of working dir.
		_ = c.SetPath(vfile.Pwd())
		// Dir path of binary.
		if selfPath := vfile.SelfDir(); selfPath != "" && vfile.Exists(selfPath) {
			_ = c.AddPath(selfPath)
		}
		// Dir path of main package.
		if mainPath := vfile.MainPkgPath(); mainPath != "" && vfile.Exists(mainPath) {
			_ = c.AddPath(mainPath)
		}
	}
	return c
}

// filePath returns the absolute configuration file path for the given filename by <file>.
func (c *Config) filePath(file ...string) (path string) {
	name := c.name
	if len(file) > 0 {
		name = file[0]
	}
	path = c.FilePath(name)
	if path == "" {
		buffer := bytes.NewBuffer(nil)
		if c.paths.Len() > 0 {
			buffer.WriteString(fmt.Sprintf("[cfg] cannot find config file \"%s\" in following paths:", name))
			c.paths.RLockFunc(func(array []string) {
				index := 1
				for _, v := range array {
					v = str.TrimRight(v, `\/`)
					buffer.WriteString(fmt.Sprintf("\n%d. %s", index, v))
					index++
					buffer.WriteString(fmt.Sprintf("\n%d. %s", index, v+vfile.Separator+"config"))
					index++
				}
			})
		} else {
			buffer.WriteString(fmt.Sprintf("[cfg] cannot find config file \"%s\" with no path set/add", name))
		}
		if errorPrint() {
			fmt.Printf(buffer.String())
		}
	}
	return path
}

// SetPath sets the configuration directory path for file search.
// The parameter <path> can be absolute or relative path,
// but absolute path is strongly recommended.
func (c *Config) SetPath(path string) error {
	var (
		isDir    = false
		realPath = ""
	)
	if file := res.Get(path); file != nil {
		realPath = path
		isDir = file.FileInfo().IsDir()
	} else {
		// Absolute path.
		realPath = vfile.RealPath(path)
		if realPath == "" {
			// Relative path.
			c.paths.RLockFunc(func(array []string) {
				for _, v := range array {
					if path, _ := spath.Search(v, path); path != "" {
						realPath = path
						break
					}
				}
			})
		}
		if realPath != "" {
			isDir = vfile.IsDir(realPath)
		}
	}
	// Path not exist.
	if realPath == "" {
		buffer := bytes.NewBuffer(nil)
		if c.paths.Len() > 0 {
			buffer.WriteString(fmt.Sprintf("[cfg] SetPath failed: cannot find directory \"%s\" in following paths:", path))
			c.paths.RLockFunc(func(array []string) {
				for k, v := range array {
					buffer.WriteString(fmt.Sprintf("\n%d. %s", k+1, v))
				}
			})
		} else {
			buffer.WriteString(fmt.Sprintf(`[cfg] SetPath failed: path "%s" does not exist`, path))
		}
		err := errors.New(buffer.String())
		if errorPrint() {
			fmt.Println(buffer.String())
		}
		return err
	}
	// Should be a directory.
	if !isDir {
		err := fmt.Errorf(`[cfg] SetPath failed: path "%s" should be directory type`, path)
		if errorPrint() {
			fmt.Println(err)
		}
		return err
	}
	// Repeated path check.
	if c.paths.Search(realPath) != -1 {
		return nil
	}
	c.jsons.Clear()
	c.paths.Clear()
	c.paths.Append(realPath)
	return nil
}

// SetViolenceCheck sets whether to perform hierarchical conflict checking.
// This feature needs to be enabled when there is a level symbol in the key name.
// It is off in default.
//
// Note that, turning on this feature is quite expensive, and it is not recommended
// to allow separators in the key names. It is best to avoid this on the application side.
func (c *Config) SetViolenceCheck(check bool) {
	c.vc = check
	c.Clear()
}

// AddPath adds a absolute or relative path to the search paths.
func (c *Config) AddPath(path string) error {
	var (
		isDir    = false
		realPath = ""
	)
	// It firstly checks the resource manager,
	// and then checks the filesystem for the path.
	if file := res.Get(path); file != nil {
		realPath = path
		isDir = file.FileInfo().IsDir()
	} else {
		// Absolute path.
		realPath = vfile.RealPath(path)
		if realPath == "" {
			// Relative path.
			c.paths.RLockFunc(func(array []string) {
				for _, v := range array {
					if path, _ := spath.Search(v, path); path != "" {
						realPath = path
						break
					}
				}
			})
		}
		if realPath != "" {
			isDir = vfile.IsDir(realPath)
		}
	}
	if realPath == "" {
		buffer := bytes.NewBuffer(nil)
		if c.paths.Len() > 0 {
			buffer.WriteString(fmt.Sprintf("[cfg] AddPath failed: cannot find directory \"%s\" in following paths:", path))
			c.paths.RLockFunc(func(array []string) {
				for k, v := range array {
					buffer.WriteString(fmt.Sprintf("\n%d. %s", k+1, v))
				}
			})
		} else {
			buffer.WriteString(fmt.Sprintf(`[cfg] AddPath failed: path "%s" does not exist`, path))
		}
		err := errors.New(buffer.String())
		if errorPrint() {
			fmt.Println(err)
		}
		return err
	}
	if !isDir {
		err := fmt.Errorf(`[cfg] AddPath failed: path "%s" should be directory type`, path)
		if errorPrint() {
			fmt.Println(err)
		}
		return err
	}
	// Repeated path check.
	if c.paths.Search(realPath) != -1 {
		return nil
	}
	c.paths.Append(realPath)
	//log.Debug("[gcfg] AddPath:", realPath)
	return nil
}

// GetFilePath returns the absolute path of the specified configuration file.
// If <file> is not passed, it returns the configuration file path of the default name.
// If the specified configuration file does not exist,
// an empty string is returned.
func (c *Config) FilePath(file ...string) (path string) {
	name := c.name
	if len(file) > 0 {
		name = file[0]
	}
	// Searching resource manager.
	if !res.IsEmpty() {
		for _, v := range resourceTryFiles {
			if file := res.Get(v + name); file != nil {
				path = file.Name()
				return
			}
		}
		c.paths.RLockFunc(func(array []string) {
			for _, prefix := range array {
				for _, v := range resourceTryFiles {
					if file := res.Get(prefix + v + name); file != nil {
						path = file.Name()
						return
					}
				}
			}
		})
	}
	// Searching the file system.
	c.paths.RLockFunc(func(array []string) {
		for _, prefix := range array {
			prefix = str.TrimRight(prefix, `\/`)
			if path, _ = spath.Search(prefix, name); path != "" {
				return
			}
			if path, _ = spath.Search(prefix+vfile.Separator+"config", name); path != "" {
				return
			}
		}
	})
	return
}

// SetFileName sets the default configuration file name.
func (c *Config) SetFileName(name string) *Config {
	c.name = name
	return c
}

// GetFileName returns the default configuration file name.
func (c *Config) GetFileName() string {
	return c.name
}

// Available checks and returns whether configuration of given <file> is available.
func (c *Config) Available(file ...string) bool {
	var name string
	if len(file) > 0 && file[0] != "" {
		name = file[0]
	} else {
		name = c.name
	}
	if c.FilePath(name) != "" {
		return true
	}
	if GetContent(name) != "" {
		return true
	}
	return false
}

// getJson returns a *json.Json object for the specified <file> content.
// It would print error if file reading fails. It return nil if any error occurs.
func (c *Config) getJson(file ...string) *json.Json {
	var name string
	if len(file) > 0 && file[0] != "" {
		name = file[0]
	} else {
		name = c.name
	}
	r := c.jsons.GetOrSetFuncLock(name, func() interface{} {
		var (
			content  = ""
			filePath = ""
		)
		if content = GetContent(name); content == "" {
			filePath = c.filePath(name)
			if filePath == "" {
				return nil
			}
			if file := res.Get(filePath); file != nil {
				content = string(file.Content())
			} else {
				content = vfile.GetContents(filePath)
			}
		}
		// Note that the underlying configuration json object operations are concurrent safe.
		if j, err := json.LoadContent(content, true); err == nil {
			j.SetViolenceCheck(c.vc)
			// Add monitor for this configuration file,
			// any changes of this file will refresh its cache in Config object.
			if filePath != "" && !res.Contains(filePath) {
				_, err = fsnotify.Add(filePath, func(event *fsnotify.Event) {
					c.jsons.Remove(name)
				})
				if err != nil && errorPrint() {
					fmt.Println(err)
				}
			}
			return j
		} else {
			if errorPrint() {
				if filePath != "" {
					// log.Criticalf(`[gcfg] Load config file "%s" failed: %s`, filePath, err.Error())
					fmt.Println(err)
				} else {
					// log.Criticalf(`[gcfg] Load configuration failed: %s`, err.Error())
					fmt.Println(err)
				}
			}
		}
		return nil
	})
	if r != nil {
		return r.(*json.Json)
	}
	return nil
}
