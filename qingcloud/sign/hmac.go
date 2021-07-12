package sign

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
)

type HmacFunc interface {
	GetName() string
	Hash(data string) ([]byte, error)
}

type hmacFuncSHA256Impl struct {
	key []byte
}

func (h *hmacFuncSHA256Impl) GetName() string {
	return "HmacSHA256"
}

func (h *hmacFuncSHA256Impl) Hash(data string) ([]byte, error) {
	hh := hmac.New(sha256.New, h.key)
	if _, err := hh.Write([]byte(data)); err != nil {
		return nil, err
	}

	return hh.Sum(nil), nil
}

type hmacFuncSHA1Impl struct {
	key []byte
}

func (h *hmacFuncSHA1Impl) GetName() string {
	return "HmacSHA1"
}

func (h *hmacFuncSHA1Impl) Hash(data string) ([]byte, error) {
	hh := hmac.New(sha1.New, h.key)
	if _, err := hh.Write([]byte(data)); err != nil {
		return nil, err
	}

	return hh.Sum(nil), nil
}

func NewHmacSHA256(key []byte) HmacFunc {
	return &hmacFuncSHA256Impl{key: key}
}

func NewHmacSHA1(key []byte) HmacFunc {
	return &hmacFuncSHA1Impl{key: key}
}

