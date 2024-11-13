package conv

import (
	"time"

	"github.com/luyingjie/utils/internal/util"

	"github.com/luyingjie/utils/container/vtime"
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
	if !util.IsNumeric(s) {
		d, _ := vtime.ParseDuration(s)
		return d
	}
	return time.Duration(Int64(i))
}

func VTime(i interface{}, format ...string) *vtime.Time {
	if i == nil {
		return nil
	}
	// It's already this type.
	if len(format) == 0 {
		if v, ok := i.(*vtime.Time); ok {
			return v
		}
	}
	s := String(i)
	if len(s) == 0 {
		return vtime.New()
	}
	// Priority conversion using given format.
	if len(format) > 0 {
		t, _ := vtime.StrToTimeFormat(s, format[0])
		return t
	}
	if util.IsNumeric(s) {
		return vtime.NewFromTimeStamp(Int64(s))
	} else {
		t, _ := vtime.StrToTime(s)
		return t
	}
}
