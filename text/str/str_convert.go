package str

import (
	"regexp"
	"strconv"
)

var (
	octReg = regexp.MustCompile(`\\[0-7]{3}`)
)

func OctStr(str string) string {
	return octReg.ReplaceAllStringFunc(
		str,
		func(s string) string {
			i, _ := strconv.ParseInt(s[1:], 8, 0)
			return string([]byte{byte(i)})
		},
	)
}
