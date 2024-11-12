// 时间类型和基础操作

package time

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	iutil "github.com/luyingjie/utils/internal/util"
	baseError "github.com/luyingjie/utils/os/error/base"

	"github.com/luyingjie/utils/text/regex"
)

const (
	D  = 24 * time.Hour
	H  = time.Hour
	M  = time.Minute
	S  = time.Second
	MS = time.Millisecond
	US = time.Microsecond
	NS = time.Nanosecond

	// Regular expression1(datetime separator supports '-', '/', '.').
	// Eg:
	// "2017-12-14 04:51:34 +0805 LMT",
	// "2017-12-14 04:51:34 +0805 LMT",
	// "2006-01-02T15:04:05Z07:00",
	// "2014-01-17T01:19:15+08:00",
	// "2018-02-09T20:46:17.897Z",
	// "2018-02-09 20:46:17.897",
	// "2018-02-09T20:46:17Z",
	// "2018-02-09 20:46:17",
	// "2018/10/31 - 16:38:46"
	// "2018-02-09",
	// "2018.02.09",
	TIME_REAGEX_PATTERN1 = `(\d{4}[-/\.]\d{2}[-/\.]\d{2})[:\sT-]*(\d{0,2}:{0,1}\d{0,2}:{0,1}\d{0,2}){0,1}\.{0,1}(\d{0,9})([\sZ]{0,1})([\+-]{0,1})([:\d]*)`

	// Regular expression2(datetime separator supports '-', '/', '.').
	// Eg:
	// 01-Nov-2018 11:50:28
	// 01/Nov/2018 11:50:28
	// 01.Nov.2018 11:50:28
	// 01.Nov.2018:11:50:28
	TIME_REAGEX_PATTERN2 = `(\d{1,2}[-/\.][A-Za-z]{3,}[-/\.]\d{4})[:\sT-]*(\d{0,2}:{0,1}\d{0,2}:{0,1}\d{0,2}){0,1}\.{0,1}(\d{0,9})([\sZ]{0,1})([\+-]{0,1})([:\d]*)`
)

var (
	timeRegex1, _ = regexp.Compile(TIME_REAGEX_PATTERN1)
	timeRegex2, _ = regexp.Compile(TIME_REAGEX_PATTERN2)

	monthMap = map[string]int{
		"jan":       1,
		"feb":       2,
		"mar":       3,
		"apr":       4,
		"may":       5,
		"jun":       6,
		"jul":       7,
		"aug":       8,
		"sep":       9,
		"sept":      9,
		"oct":       10,
		"nov":       11,
		"dec":       12,
		"january":   1,
		"february":  2,
		"march":     3,
		"april":     4,
		"june":      6,
		"july":      7,
		"august":    8,
		"september": 9,
		"october":   10,
		"november":  11,
		"december":  12,
	}
)

func SetTimeZone(zone string) error {
	location, err := time.LoadLocation(zone)
	if err == nil {
		time.Local = location
	}
	return err
}

func Timestamp() int64 {
	return Now().Timestamp()
}

func TimestampMilli() int64 {
	return Now().TimestampMilli()
}

func TimestampMicro() int64 {
	return Now().TimestampMicro()
}

func TimestampNano() int64 {
	return Now().TimestampNano()
}

func TimestampStr() string {
	return Now().TimestampStr()
}

func TimestampMilliStr() string {
	return Now().TimestampMilliStr()
}

func TimestampMicroStr() string {
	return Now().TimestampMicroStr()
}

func TimestampNanoStr() string {
	return Now().TimestampNanoStr()
}

func Second() int64 {
	return Timestamp()
}

func Millisecond() int64 {
	return TimestampMilli()
}

func Microsecond() int64 {
	return TimestampMicro()
}

func Nanosecond() int64 {
	return TimestampNano()
}

// Date returns current date in string like "2006-01-02".
func Date() string {
	return time.Now().Format("2006-01-02")
}

// Datetime returns current datetime in string like "2006-01-02 15:04:05".
func Datetime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// ISO8601 returns current datetime in ISO8601 format like "2006-01-02T15:04:05-07:00".
func ISO8601() string {
	return time.Now().Format("2006-01-02T15:04:05-07:00")
}

// ISO8601 returns current datetime in RFC822 format like "Mon, 02 Jan 06 15:04 MST".
func RFC822() string {
	return time.Now().Format("Mon, 02 Jan 06 15:04 MST")
}

func parseDateStr(s string) (year, month, day int) {
	array := strings.Split(s, "-")
	if len(array) < 3 {
		array = strings.Split(s, "/")
	}
	if len(array) < 3 {
		array = strings.Split(s, ".")
	}

	if len(array) < 3 {
		return
	}

	if iutil.IsNumeric(array[1]) {
		year, _ = strconv.Atoi(array[0])
		month, _ = strconv.Atoi(array[1])
		day, _ = strconv.Atoi(array[2])
	} else {
		if v, ok := monthMap[strings.ToLower(array[1])]; ok {
			month = v
		} else {
			return
		}
		year, _ = strconv.Atoi(array[2])
		day, _ = strconv.Atoi(array[0])
	}
	return
}

func StrToTime(str string, format ...string) (*Time, error) {
	if len(format) > 0 {
		return StrToTimeFormat(str, format[0])
	}
	var (
		year, month, day     int
		hour, min, sec, nsec int
		match                []string
		local                = time.Local
	)
	if match = timeRegex1.FindStringSubmatch(str); len(match) > 0 && match[1] != "" {
		for k, v := range match {
			match[k] = strings.TrimSpace(v)
		}
		year, month, day = parseDateStr(match[1])
	} else if match = timeRegex2.FindStringSubmatch(str); len(match) > 0 && match[1] != "" {
		for k, v := range match {
			match[k] = strings.TrimSpace(v)
		}
		year, month, day = parseDateStr(match[1])
	} else {
		return nil, errors.New("unsupported time format")
	}

	// Time
	if len(match[2]) > 0 {
		s := strings.Replace(match[2], ":", "", -1)
		if len(s) < 6 {
			s += strings.Repeat("0", 6-len(s))
		}
		hour, _ = strconv.Atoi(s[0:2])
		min, _ = strconv.Atoi(s[2:4])
		sec, _ = strconv.Atoi(s[4:6])
	}

	if len(match[3]) > 0 {
		nsec, _ = strconv.Atoi(match[3])
		for i := 0; i < 9-len(match[3]); i++ {
			nsec *= 10
		}
	}

	if match[4] != "" && match[6] == "" {
		match[6] = "000000"
	}

	if match[6] != "" {
		zone := strings.Replace(match[6], ":", "", -1)
		zone = strings.TrimLeft(zone, "+-")
		if len(zone) <= 6 {
			zone += strings.Repeat("0", 6-len(zone))
			h, _ := strconv.Atoi(zone[0:2])
			m, _ := strconv.Atoi(zone[2:4])
			s, _ := strconv.Atoi(zone[4:6])
			if h > 24 || m > 59 || s > 59 {
				return nil, baseError.Newf("invalid zone string: %s", match[6])
			}

			_, localOffset := time.Now().Zone()
			if (h*3600 + m*60 + s) != localOffset {
				local = time.UTC
				// UTC conversion.
				operation := match[5]
				if operation != "+" && operation != "-" {
					operation = "-"
				}
				switch operation {
				case "+":
					if h > 0 {
						hour -= h
					}
					if m > 0 {
						min -= m
					}
					if s > 0 {
						sec -= s
					}
				case "-":
					if h > 0 {
						hour += h
					}
					if m > 0 {
						min += m
					}
					if s > 0 {
						sec += s
					}
				}
			}
		}
	}
	if year <= 0 || month <= 0 || day <= 0 || hour < 0 || min < 0 || sec < 0 || nsec < 0 {
		return nil, errors.New("invalid time string:" + str)
	}

	return NewFromTime(time.Date(year, time.Month(month), day, hour, min, sec, nsec, local)), nil
}

func ConvertZone(strTime string, toZone string, fromZone ...string) (*Time, error) {
	t, err := StrToTime(strTime)
	if err != nil {
		return nil, err
	}
	if len(fromZone) > 0 {
		if l, err := time.LoadLocation(fromZone[0]); err != nil {
			return nil, err
		} else {
			t.Time = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Time.Second(), t.Time.Nanosecond(), l)
		}
	}
	if l, err := time.LoadLocation(toZone); err != nil {
		return nil, err
	} else {
		return t.ToLocation(l), nil
	}
}

func StrToTimeFormat(str string, format string) (*Time, error) {
	return StrToTimeLayout(str, formatToStdLayout(format))
}

func StrToTimeLayout(str string, layout string) (*Time, error) {
	if t, err := time.ParseInLocation(layout, str, time.Local); err == nil {
		return NewFromTime(t), nil
	} else {
		return nil, err
	}
}

func ParseTimeFromContent(content string, format ...string) *Time {
	if len(format) > 0 {
		if match, err := regex.MatchString(formatToRegexPattern(format[0]), content); err == nil && len(match) > 0 {
			return NewFromStrFormat(match[0], format[0])
		}
	} else {
		if match := timeRegex1.FindStringSubmatch(content); len(match) >= 1 {
			return NewFromStr(strings.Trim(match[0], "./_- \n\r"))
		} else if match := timeRegex2.FindStringSubmatch(content); len(match) >= 1 {
			return NewFromStr(strings.Trim(match[0], "./_- \n\r"))
		}
	}
	return nil
}

func ParseDuration(s string) (time.Duration, error) {
	if iutil.IsNumeric(s) {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return time.Duration(v), nil
	}
	match, err := regex.MatchString(`^([\-\d]+)[dD](.*)$`, s)
	if err != nil {
		return 0, err
	}
	if len(match) == 3 {
		v, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0, err
		}
		return time.ParseDuration(fmt.Sprintf(`%dh%s`, v*24, match[2]))
	}
	return time.ParseDuration(s)
}

func FuncCost(f func()) int64 {
	t := TimestampNano()
	f()
	return TimestampNano() - t
}
