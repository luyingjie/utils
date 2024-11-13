package vtime

import (
	"bytes"
	"strconv"
	"time"
)

type Time struct {
	TimeWrapper
}

func New(t ...time.Time) *Time {
	if len(t) > 0 {
		return NewFromTime(t[0])
	}
	return &Time{
		TimeWrapper{time.Time{}},
	}
}

func Now() *Time {
	return &Time{
		TimeWrapper{time.Now()},
	}
}

func NewFromTime(t time.Time) *Time {
	return &Time{
		TimeWrapper{t},
	}
}

func NewFromStr(str string) *Time {
	if t, err := StrToTime(str); err == nil {
		return t
	}
	return nil
}

func NewFromStrFormat(str string, format string) *Time {
	if t, err := StrToTimeFormat(str, format); err == nil {
		return t
	}
	return nil
}

func NewFromStrLayout(str string, layout string) *Time {
	if t, err := StrToTimeLayout(str, layout); err == nil {
		return t
	}
	return nil
}

func NewFromTimeStamp(timestamp int64) *Time {
	if timestamp == 0 {
		return &Time{}
	}
	var sec, nano int64
	if timestamp > 1e9 {
		for timestamp < 1e18 {
			timestamp *= 10
		}
		sec = timestamp / 1e9
		nano = timestamp % 1e9
	} else {
		sec = timestamp
	}
	return &Time{
		TimeWrapper{time.Unix(sec, nano)},
	}
}

func (t *Time) Timestamp() int64 {
	return t.UnixNano() / 1e9
}

func (t *Time) TimestampMilli() int64 {
	return t.UnixNano() / 1e6
}

func (t *Time) TimestampMicro() int64 {
	return t.UnixNano() / 1e3
}

func (t *Time) TimestampNano() int64 {
	return t.UnixNano()
}

func (t *Time) TimestampStr() string {
	return strconv.FormatInt(t.Timestamp(), 10)
}

func (t *Time) TimestampMilliStr() string {
	return strconv.FormatInt(t.TimestampMilli(), 10)
}

func (t *Time) TimestampMicroStr() string {
	return strconv.FormatInt(t.TimestampMicro(), 10)
}

func (t *Time) TimestampNanoStr() string {
	return strconv.FormatInt(t.TimestampNano(), 10)
}

func (t *Time) Second() int {
	return t.Time.Second()
}

func (t *Time) Millisecond() int {
	return t.Time.Nanosecond() / 1e6
}

func (t *Time) Microsecond() int {
	return t.Time.Nanosecond() / 1e3
}

func (t *Time) Nanosecond() int {
	return t.Time.Nanosecond()
}

func (t *Time) String() string {
	if t == nil {
		return ""
	}
	if t.IsZero() {
		return ""
	}
	return t.Format("Y-m-d H:i:s")
}

func (t *Time) Clone() *Time {
	return New(t.Time)
}

func (t *Time) Add(d time.Duration) *Time {
	t.Time = t.Time.Add(d)
	return t
}

func (t *Time) AddStr(duration string) error {
	if d, err := time.ParseDuration(duration); err != nil {
		return err
	} else {
		t.Time = t.Time.Add(d)
	}
	return nil
}

func (t *Time) ToLocation(location *time.Location) *Time {
	t.Time = t.Time.In(location)
	return t
}

func (t *Time) ToZone(zone string) (*Time, error) {
	if l, err := time.LoadLocation(zone); err == nil {
		t.Time = t.Time.In(l)
		return t, nil
	} else {
		return nil, err
	}
}

func (t *Time) UTC() *Time {
	t.Time = t.Time.UTC()
	return t
}

func (t *Time) ISO8601() string {
	return t.Layout("2006-01-02T15:04:05-07:00")
}

func (t *Time) RFC822() string {
	return t.Layout("Mon, 02 Jan 06 15:04 MST")
}

func (t *Time) Local() *Time {
	t.Time = t.Time.Local()
	return t
}

func (t *Time) AddDate(years int, months int, days int) *Time {
	t.Time = t.Time.AddDate(years, months, days)
	return t
}

func (t *Time) Round(d time.Duration) *Time {
	t.Time = t.Time.Round(d)
	return t
}

func (t *Time) Truncate(d time.Duration) *Time {
	t.Time = t.Time.Truncate(d)
	return t
}

func (t *Time) Equal(u *Time) bool {
	return t.Time.Equal(u.Time)
}

func (t *Time) Before(u *Time) bool {
	return t.Time.Before(u.Time)
}

func (t *Time) After(u *Time) bool {
	return t.Time.After(u.Time)
}

func (t *Time) Sub(u *Time) time.Duration {
	return t.Time.Sub(u.Time)
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		t.Time = time.Time{}
		return nil
	}
	newTime, err := StrToTime(string(bytes.Trim(b, `"`)))
	if err != nil {
		return err
	}
	t.Time = newTime.Time
	return nil
}
