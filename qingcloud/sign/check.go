// QingCloud KMS API 反向版本，用于检查加密。

package sign

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
)

func CheckSign(r *http.Request, keyid string, secret string, params string) (string, error) {
	s := Signer{
		AccessKeyID:     keyid,
		SecretAccessKey: secret,
	}
	hf := NewHmacSHA256([]byte(s.SecretAccessKey))

	if params != "" {
		r.URL.RawQuery = params
	}

	l, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return "", err
	}

	l.Del("signature")
	r.URL.RawQuery = l.Encode()

	signedHeaders := SignedHeaders(r)
	canonicalRequest, err := CanonicalRequest(r, signedHeaders)
	if err != nil {
		return "", err
	}

	res, err := hf.Hash(canonicalRequest)
	if err != nil {
		return "", err
	}

	signature := strings.TrimSpace(base64.StdEncoding.EncodeToString(res))
	signature = strings.Replace(signature, " ", "+", -1)
	signature = url.QueryEscape(signature)

	return signature, nil
}
