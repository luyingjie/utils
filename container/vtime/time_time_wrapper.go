package vtime

import (
	"time"
)

type TimeWrapper struct {
	time.Time
}

func (t TimeWrapper) String() string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
