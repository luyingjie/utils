package file

import (
	"utils/text/str"
)

func ReplaceFile(search, replace, path string) error {
	return PutContents(path, str.Replace(GetContents(path), search, replace))
}

func ReplaceFileFunc(f func(path, content string) string, path string) error {
	data := GetContents(path)
	result := f(path, data)
	if len(data) != len(result) && data != result {
		return PutContents(path, result)
	}
	return nil
}

func ReplaceDir(search, replace, path, pattern string, recursive ...bool) error {
	files, err := ScanDirFile(path, pattern, recursive...)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err = ReplaceFile(search, replace, file); err != nil {
			return err
		}
	}
	return err
}

func ReplaceDirFunc(f func(path, content string) string, path, pattern string, recursive ...bool) error {
	files, err := ScanDirFile(path, pattern, recursive...)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err = ReplaceFileFunc(f, file); err != nil {
			return err
		}
	}
	return err
}
