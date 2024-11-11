package str

import "strings"

func Contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

func ContainsI(str, substr string) bool {
	return PosI(str, substr) != -1
}

func ContainsAny(s, chars string) bool {
	return strings.ContainsAny(s, chars)
}
