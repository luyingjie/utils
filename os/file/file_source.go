package file

import (
	"os"
	"runtime"
	"strings"

	"github.com/luyingjie/utils/text/regex"
)

var (
	goRootForFilter = runtime.GOROOT()
)

func init() {
	if goRootForFilter != "" {
		goRootForFilter = strings.Replace(goRootForFilter, "\\", "/", -1)
	}
}

func MainPkgPath() string {
	if goRootForFilter == "" {
		return ""
	}
	path := mainPkgPath.Val()
	if path != "" {
		return path
	}
	lastFile := ""
	for i := 1; i < 10000; i++ {
		if _, file, _, ok := runtime.Caller(i); ok {
			if goRootForFilter != "" && len(file) >= len(goRootForFilter) && file[0:len(goRootForFilter)] == goRootForFilter {
				continue
			}
			if regex.IsMatchString(`/github.com/[^/]+/gf/`, file) &&
				!regex.IsMatchString(`/github.com/[^/]+/gf/\.example/`, file) {
				continue
			}
			if Ext(file) != ".go" {
				continue
			}
			lastFile = file
			if regex.IsMatchString(`package\s+main`, GetContents(file)) {
				mainPkgPath.Set(Dir(file))
				return Dir(file)
			}
		} else {
			break
		}
	}
	if lastFile != "" {
		for path = Dir(lastFile); len(path) > 1 && Exists(path) && path[len(path)-1] != os.PathSeparator; {
			files, _ := ScanDir(path, "*.go")
			for _, v := range files {
				if regex.IsMatchString(`package\s+main`, GetContents(v)) {
					mainPkgPath.Set(path)
					return path
				}
			}
			path = Dir(path)
		}
	}
	return ""
}
