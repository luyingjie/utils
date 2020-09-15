package file

import (
	"strings"

	"utils/container/array"
)

func fileSortFunc(path1, path2 string) int {
	isDirPath1 := IsDir(path1)
	isDirPath2 := IsDir(path2)
	if isDirPath1 && !isDirPath2 {
		return -1
	}
	if !isDirPath1 && isDirPath2 {
		return 1
	}
	if n := strings.Compare(path1, path2); n != 0 {
		return n
	} else {
		return -1
	}
}

func SortFiles(files []string) []string {
	array := array.NewSortedStrArrayComparator(fileSortFunc)
	array.Add(files...)
	return array.Slice()
}
