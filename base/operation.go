package base

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Sha1Stream struct {
	_sha1 hash.Hash
}

func (obj *Sha1Stream) Update(data []byte) {
	if obj._sha1 == nil {
		obj._sha1 = sha1.New()
	}
	obj._sha1.Write(data)
}

func (obj *Sha1Stream) Sum() string {
	return hex.EncodeToString(obj._sha1.Sum([]byte("")))
}

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

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
