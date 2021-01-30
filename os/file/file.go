package file

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"utils/text/str"

	"utils/os/conv"

	vtype "utils/container/type"
)

var (
	Separator = string(filepath.Separator)

	DefaultPermOpen = os.FileMode(0666)

	DefaultPermCopy = os.FileMode(0777)

	mainPkgPath = vtype.NewString()

	selfPath = ""

	tempDir = "/tmp"
)

func init() {
	if Separator != "/" || !Exists(tempDir) {
		tempDir = os.TempDir()
	}
	selfPath, _ = exec.LookPath(os.Args[0])
	if selfPath != "" {
		selfPath, _ = filepath.Abs(selfPath)
	}
	if selfPath == "" {
		selfPath, _ = filepath.Abs(os.Args[0])
	}
}

func Mkdir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func Create(path string) (*os.File, error) {
	dir := Dir(path)
	if !Exists(dir) {
		if err := Mkdir(dir); err != nil {
			return nil, err
		}
	}
	return os.Create(path)
}

func Open(path string) (*os.File, error) {
	return os.Open(path)
}

func OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(path, flag, perm)
}

func OpenWithFlag(path string, flag int) (*os.File, error) {
	f, err := os.OpenFile(path, flag, DefaultPermOpen)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func OpenWithFlagPerm(path string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func Join(paths ...string) string {
	var s string
	for _, path := range paths {
		if s != "" {
			s += Separator
		}
		s += str.TrimRight(path, Separator)
	}
	return s
}

func Exists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func Pwd() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}

func Chdir(dir string) error {
	return os.Chdir(dir)
}

func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

func Info(path string) (os.FileInfo, error) {
	return Stat(path)
}

func Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func Move(src string, dst string) error {
	return os.Rename(src, dst)
}

func Rename(src string, dst string) error {
	return Move(src, dst)
}

func DirNames(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func Glob(pattern string, onlyNames ...bool) ([]string, error) {
	if list, err := filepath.Glob(pattern); err == nil {
		if len(onlyNames) > 0 && onlyNames[0] && len(list) > 0 {
			array := make([]string, len(list))
			for k, v := range list {
				array[k] = Basename(v)
			}
			return array, nil
		}
		return list, nil
	} else {
		return nil, err
	}
}

func Remove(path string) error {
	return os.RemoveAll(path)
}

func IsReadable(path string) bool {
	result := true
	file, err := os.OpenFile(path, os.O_RDONLY, DefaultPermOpen)
	if err != nil {
		result = false
	}
	file.Close()
	return result
}

func IsWritable(path string) bool {
	result := true
	if IsDir(path) {
		tmpFile := strings.TrimRight(path, Separator) + Separator + conv.String(time.Now().UnixNano())
		if f, err := Create(tmpFile); err != nil || !Exists(tmpFile) {
			result = false
		} else {
			f.Close()
			Remove(tmpFile)
		}
	} else {
		// 如果是文件，那么判断文件是否可打开
		file, err := os.OpenFile(path, os.O_WRONLY, DefaultPermOpen)
		if err != nil {
			result = false
		}
		file.Close()
	}
	return result
}

func Chmod(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

func Abs(path string) string {
	p, _ := filepath.Abs(path)
	return p
}

func RealPath(path string) string {
	p, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	if !Exists(p) {
		return ""
	}
	return p
}

func SelfPath() string {
	return selfPath
}

func SelfName() string {
	return Basename(SelfPath())
}

func SelfDir() string {
	return filepath.Dir(SelfPath())
}

func Basename(path string) string {
	return filepath.Base(path)
}

func Name(path string) string {
	base := filepath.Base(path)
	if i := strings.LastIndexByte(base, '.'); i != -1 {
		return base[:i]
	}
	return base
}

func Dir(path string) string {
	return filepath.Dir(path)
}

func IsEmpty(path string) bool {
	stat, err := Stat(path)
	if err != nil {
		return true
	}
	if stat.IsDir() {
		file, err := os.Open(path)
		if err != nil {
			return true
		}
		defer file.Close()
		names, err := file.Readdirnames(-1)
		if err != nil {
			return true
		}
		return len(names) == 0
	} else {
		return stat.Size() == 0
	}
}

func Ext(path string) string {
	ext := filepath.Ext(path)
	if p := strings.IndexByte(ext, '?'); p != -1 {
		ext = ext[0:p]
	}
	return ext
}

func ExtName(path string) string {
	return strings.TrimLeft(Ext(path), ".")
}

func TempDir(names ...string) string {
	path := tempDir
	for _, name := range names {
		path += Separator + name
	}
	return path
}
