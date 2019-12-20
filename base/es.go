package base

import "strings"

// SetESSelectValue : ES查询条件的判断
func SetESSelectValue(val string) (string, string) {
	// strings.Contains("helloogo", "hello")
	// strings.Replace(license, "\\n", "", -1)

	t := "Wildcard" // QueryString
	is := strings.ContainsAny(val, "AND&OR")
	if is {
		t = "QueryString"
		val = SetESQueryStringValue(val)
	}

	return t, val
}

// SetESQueryStringValue : 设置查询空格
func SetESQueryStringValue(val string) string {
	val = strings.Replace(val, " AND ", "AND", -1)
	val = strings.Replace(val, " OR ", "OR", -1)
	val = strings.Replace(val, " ", "\\\\ ", -1)
	val = strings.Replace(val, "AND", " AND ", -1)
	val = strings.Replace(val, "OR", " OR ", -1)
	return val
}
