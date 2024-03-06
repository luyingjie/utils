package file

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/luyingjie/utils/container/array"
)

func Search(name string, prioritySearchPaths ...string) (realPath string, err error) {
	realPath = RealPath(name)
	if realPath != "" {
		return
	}

	array := array.NewStrArray()
	array.Append(prioritySearchPaths...)
	array.Append(Pwd(), SelfDir())
	if path := MainPkgPath(); path != "" {
		array.Append(path)
	}

	array.Unique()

	array.RLockFunc(func(array []string) {
		path := ""
		for _, v := range array {
			path = RealPath(v + Separator + name)
			if path != "" {
				realPath = path
				break
			}
		}
	})

	if realPath == "" {
		buffer := bytes.NewBuffer(nil)
		buffer.WriteString(fmt.Sprintf("cannot find file/folder \"%s\" in following paths:", name))
		array.RLockFunc(func(array []string) {
			for k, v := range array {
				buffer.WriteString(fmt.Sprintf("\n%d. %s", k+1, v))
			}
		})
		err = errors.New(buffer.String())
	}
	return
}
