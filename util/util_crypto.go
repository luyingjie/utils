package util

import (
	"encoding/base64"
	"strings"
)

//将字符串加密成MD5
func PassMD5(t string) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(t)), "=")
}

//将一个MD5字串解密。
func ToMD5(s string) string {
	num := 4 - len(s)%4
	for i := 0; i < num; i++ {
		s = s + "="
	}
	decoded, _ := base64.URLEncoding.DecodeString(s)

	return string(decoded)
}
