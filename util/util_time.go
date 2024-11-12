package util

import "time"

// 时间测试
// fmt.Println("2021-07-04T18:28:03Z")
// fmt.Println(time.Now().Add(time.Hour).Format("2006-01-02T15:04:05Z"))
// fmt.Println(TimeToString(time_stamp, "ISO 8601"))
// fmt.Println(time.Now().UTC().Format(time.RFC3339))
// fmt.Println(TimeToString(time_stamp, "RFC 822"))
// fmt.Println(time.Now().UTC().Format(http.TimeFormat))


var layoutMap = map[string]string{
	"RFC 822":  "Mon, 02 Jan 2006 15:04:05 GMT",
	"ISO 8601": "2006-01-02T15:04:05Z",
}

// TimeToString transforms given time to string.
func TimeToString(timeValue time.Time, format string) string {
	return timeValue.UTC().Format(layoutMap[format])
}