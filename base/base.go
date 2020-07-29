package base

import (
	"encoding/base64"
	"strings"
)

// 深克隆
func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		newMap := make(map[string]interface{})
		for k, v := range valueMap {
			newMap[k] = DeepCopy(v)
		}

		return newMap
	} else if valueSlice, ok := value.([]interface{}); ok {
		newSlice := make([]interface{}, len(valueSlice))
		for k, v := range valueSlice {
			newSlice[k] = DeepCopy(v)
		}

		return newSlice
	}

	return value
}

// 好像是获取字符串的制定位置的字符。
func Substr(s string, l int) string {
	if len(s) <= l {
		return s
	}
	ss, sl, rl, rs := "", 0, 0, []rune(s)
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			rl = 1
		} else {
			rl = 2
		}
		if sl+rl > l {
			break
		}
		sl += rl
		ss += string(r)
	}
	return ss
}

// 查找数组中的元素
func InArray(arr []interface{}, a interface{}) bool {
	ist := false
	for _, value := range arr {
		switch value.(type) {
		case int:
			if value.(int) == a.(int) {
				ist = true
				break
			}
		case string:
			if value.(string) == a.(string) {
				ist = true
				break
			}
			// default:
		}
	}
	return ist
}

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
