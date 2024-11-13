package res

import (
	"archive/zip"
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/luyingjie/utils/conv"
	"github.com/luyingjie/utils/encoding/base64"
	"github.com/luyingjie/utils/encoding/compress"
	"github.com/luyingjie/utils/text/str"

	"github.com/luyingjie/utils/file/file"
)

const (
	PACKAGE_TEMPLATE = `
	package %s

	import "github.com/luyingjie/utils/os/res"

	func init() {
		if err := res.Add("%s"); err != nil {
			panic("add binary content to resource manager failed: " + err.Error())
		}
	}
`
)

// Pack packs the path specified by <srcPaths> into bytes.
// The unnecessary parameter <keyPrefix> indicates the prefix for each file
// packed into the result bytes.
//
// Note that parameter <srcPaths> supports multiple paths join with ','.
func Pack(srcPaths string, keyPrefix ...string) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	headerPrefix := ""
	if len(keyPrefix) > 0 && keyPrefix[0] != "" {
		headerPrefix = keyPrefix[0]
	}
	err := zipPathWriter(srcPaths, buffer, headerPrefix)
	if err != nil {
		return nil, err
	}
	// Gzip the data bytes to reduce the size.
	return compress.Gzip(buffer.Bytes(), 9)
}

// PackToFile packs the path specified by <srcPaths> to target file <dstPath>.
// The unnecessary parameter <keyPrefix> indicates the prefix for each file
// packed into the result bytes.
//
// Note that parameter <srcPaths> supports multiple paths join with ','.
func PackToFile(srcPaths, dstPath string, keyPrefix ...string) error {
	data, err := Pack(srcPaths, keyPrefix...)
	if err != nil {
		return err
	}
	return file.PutBytes(dstPath, data)
}

// PackToGoFile packs the path specified by <srcPaths> to target go file <goFilePath>
// with given package name <pkgName>.
//
// The unnecessary parameter <keyPrefix> indicates the prefix for each file
// packed into the result bytes.
//
// Note that parameter <srcPaths> supports multiple paths join with ','.
func PackToGoFile(srcPath, goFilePath, pkgName string, keyPrefix ...string) error {
	data, err := Pack(srcPath, keyPrefix...)
	if err != nil {
		return err
	}
	return file.PutContents(
		goFilePath,
		fmt.Sprintf(str.TrimLeft(PACKAGE_TEMPLATE), pkgName, base64.EncodeToString(data)),
	)
}

// Unpack unpacks the content specified by <path> to []*File.
func Unpack(path string) ([]*File, error) {
	realPath, err := file.Search(path)
	if err != nil {
		return nil, err
	}
	return UnpackContent(file.GetContents(realPath))
}

// UnpackContent unpacks the content to []*File.
func UnpackContent(content string) ([]*File, error) {
	var data []byte
	var err error
	if isHexStr(content) {
		// It here keeps compatible with old version packing string using hex string.
		// TODO remove this support in the future.
		data, err = compress.UnGzip(hexStrToBytes(content))
		if err != nil {
			return nil, err
		}
	} else if isBase64(content) {
		// New version packing string using base64.
		b, err := base64.DecodeString(content)
		if err != nil {
			return nil, err
		}
		data, err = compress.UnGzip(b)
		if err != nil {
			return nil, err
		}
	} else {
		data, err = compress.UnGzip(conv.UnsafeStrToBytes(content))
		if err != nil {
			return nil, err
		}
	}
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	array := make([]*File, len(reader.File))
	for i, file := range reader.File {
		array[i] = &File{file: file}
	}
	return array, nil
}

// isBase64 checks and returns whether given content <s> is base64 string.
// It returns true if <s> is base64 string, or false if not.
func isBase64(s string) bool {
	var r bool
	for i := 0; i < len(s); i++ {
		r = (s[i] >= '0' && s[i] <= '9') ||
			(s[i] >= 'a' && s[i] <= 'z') ||
			(s[i] >= 'A' && s[i] <= 'Z') ||
			(s[i] == '+' || s[i] == '-') ||
			(s[i] == '_' || s[i] == '/') || s[i] == '='
		if !r {
			return false
		}
	}
	return true
}

// isHexStr checks and returns whether given content <s> is hex string.
// It returns true if <s> is hex string, or false if not.
func isHexStr(s string) bool {
	var r bool
	for i := 0; i < len(s); i++ {
		r = (s[i] >= '0' && s[i] <= '9') ||
			(s[i] >= 'a' && s[i] <= 'f') ||
			(s[i] >= 'A' && s[i] <= 'F')
		if !r {
			return false
		}
	}
	return true
}

// hexStrToBytes converts hex string content to []byte.
func hexStrToBytes(s string) []byte {
	src := conv.UnsafeStrToBytes(s)
	dst := make([]byte, hex.DecodedLen(len(src)))
	hex.Decode(dst, src)
	return dst
}
