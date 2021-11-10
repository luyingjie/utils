package qingcloud

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

func Get_iaas_authorization(secret_access_key, method, uri, params string) string {
	method = strings.ToTitle(method)
	string_to_sign := fmt.Sprintf("%s\n%s\n%s", method, uri, params)
	h := hmacSha256(string_to_sign, secret_access_key)
	signature := strings.TrimSpace(base64.StdEncoding.EncodeToString(h))
	signature = strings.Replace(signature, " ", "+", -1)
	signature = url.QueryEscape(signature)
	return signature
}

func Get_api_authorization(secret_access_key, method, uri, data, params string) string {
	md5Params := ""
	if data != "" {
		md5Params = fmt.Sprintf("%x", md5.Sum([]byte(data)))
	}
	method = strings.ToTitle(method)
	string_to_sign := fmt.Sprintf("%s\n%s\n%s\n%s", method, uri, params, md5Params)
	h := hmacSha256(string_to_sign, secret_access_key)
	signature := strings.TrimSpace(base64.StdEncoding.EncodeToString(h))
	signature = strings.Replace(signature, " ", "+", -1)
	signature = url.QueryEscape(signature)
	return signature
}

func hmacSha256(data string, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	// return hex.EncodeToString(h.Sum(nil))
	return h.Sum(nil)
}

func QueryEscape(value string) string {
	value = strings.TrimSpace(value)
	value = url.QueryEscape(value)
	value = strings.Replace(value, "+", "%20", -1)
	return value
}
