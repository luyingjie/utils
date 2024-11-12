package file

import (
	"os"
	"time"
)

func MTime(path string) time.Time {
	s, e := os.Stat(path)
	if e != nil {
		return time.Time{}
	}
	return s.ModTime()
}

func MTimestamp(path string) int64 {
	mtime := MTime(path)
	if mtime.IsZero() {
		return -1
	}
	return mtime.Unix()
}

func MTimestampMilli(path string) int64 {
	mtime := MTime(path)
	if mtime.IsZero() {
		return -1
	}
	return mtime.UnixNano() / 1000000
}
