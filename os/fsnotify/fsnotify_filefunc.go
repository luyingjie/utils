package fsnotify

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func fileDir(path string) string {
	return filepath.Dir(path)
}

func fileRealPath(path string) string {
	p, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	if !fileExists(p) {
		return ""
	}
	return p
}

func fileExists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

func fileIsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func fileAllDirs(path string) (list []string) {
	list = []string{path}
	file, err := os.Open(path)
	if err != nil {
		return list
	}
	defer file.Close()
	names, err := file.Readdirnames(-1)
	if err != nil {
		return list
	}
	for _, name := range names {
		path := fmt.Sprintf("%s%s%s", path, string(filepath.Separator), name)
		if fileIsDir(path) {
			if array := fileAllDirs(path); len(array) > 0 {
				list = append(list, array...)
			}
		}
	}
	return
}

func fileScanDir(path string, pattern string, recursive ...bool) ([]string, error) {
	list, err := doFileScanDir(path, pattern, recursive...)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

func doFileScanDir(path string, pattern string, recursive ...bool) ([]string, error) {
	list := ([]string)(nil)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	names, err := file.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	filePath := ""
	for _, name := range names {
		filePath = fmt.Sprintf("%s%s%s", path, string(filepath.Separator), name)
		if fileIsDir(filePath) && len(recursive) > 0 && recursive[0] {
			array, _ := doFileScanDir(filePath, pattern, true)
			if len(array) > 0 {
				list = append(list, array...)
			}
		}
		for _, p := range strings.Split(pattern, ",") {
			if match, err := filepath.Match(strings.TrimSpace(p), name); err == nil && match {
				list = append(list, filePath)
			}
		}
	}
	return list, nil
}
