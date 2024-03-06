package crc32

import (
	"hash/crc32"

	"github.com/luyingjie/utils/conv"
)

// Encrypt encrypts any type of variable using CRC32 algorithms.
// It uses conv package to convert <v> to its bytes type.
func Encrypt(v interface{}) uint32 {
	return crc32.ChecksumIEEE(conv.Bytes(v))
}
