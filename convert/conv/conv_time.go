package conv

import (
	"time"

	"utils/utils"

	mtime "utils/os/time"
)

func Time(i interface{}, format ...string) time.Time {
	// It's already this type.
	if len(format) == 0 {
		if v, ok := i.(time.Time); ok {
			return v
		}
	}
	if t := VTime(i, format...); t != nil {
		return t.Time
	}
	return time.Time{}
}

func Duration(i interface{}) time.Duration {
	// It's already this type.
	if v, ok := i.(time.Duration); ok {
		return v
	}
	s := String(i)
	if !utils.IsNumeric(s) {
		d, _ := mtime.ParseDuration(s)
		return d
	}
	return time.Duration(Int64(i))
}

func VTime(i interface{}, format ...string) *mtime.Time {
	if i == nil {
		return nil
	}
	// It's already this type.
	if len(format) == 0 {
		if v, ok := i.(*mtime.Time); ok {
			return v
		}
	}
	s := String(i)
	if len(s) == 0 {
		return mtime.New()
	}
	// Priority conversion using given format.
	if len(format) > 0 {
		t, _ := mtime.StrToTimeFormat(s, format[0])
		return t
	}
	if utils.IsNumeric(s) {
		return mtime.NewFromTimeStamp(Int64(s))
	} else {
		t, _ := mtime.StrToTime(s)
		return t
	}
}
