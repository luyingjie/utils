// QingCloud KMS API Gateway Signature
// based on https://github.com/datastream/aws/blob/master/signv4.go

package sign

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

var (
	HeaderHost          = "host"
	QueryAPIVersionKey  = "version"
	QuerySignVersionKey = "signature_version"
	QueryTimestampKey   = "timestamp"
	QueryAccessKeyIdKey = "access_key_id"
	QuerySignMethodKey  = "signature_method"
	QuerySignatureKey   = "signature"
	QuerySkipped        = []string{"signature"}
)

// Build a CanonicalRequest from a regular request string
//
// CanonicalRequest =
//  HTTPRequestMethod + '\n' +
//  CanonicalURI + '\n' +
//  CanonicalQueryString + '\n' +
//  CanonicalHeaders + '\n' +
//  SignedHeaders + '\n' +
//  HexEncode(Hash(RequestPayload))
func CanonicalRequest(r *http.Request, signedHeaders []string) (string, error) {
	data, err := RequestPayload(r)
	if err != nil {
		return "", err
	}
	hexencode, err := HexEncodeMd5Hash(data)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s", r.Method, CanonicalURI(r), CanonicalQueryString(r), hexencode), nil
}

func CanonicalURI(r *http.Request) string {
	pattens := strings.Split(r.URL.Path, "/")
	var uri []string
	for _, v := range pattens {
		uri = append(uri, escape(v))
	}
	urlpath := strings.Join(uri, "/")
	if len(urlpath) == 0 || urlpath[len(urlpath)-1] != '/' {
		urlpath = urlpath + "/"
	}
	return urlpath
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func CanonicalQueryString(r *http.Request) string {
	var keys []string
	query := r.URL.Query()
	for key := range query {
		if ! contains(QuerySkipped, key) {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	var a []string
	for _, key := range keys {
		k := escape(key)
		sort.Strings(query[key])
		for _, v := range query[key] {
			kv := fmt.Sprintf("%s=%s", k, escape(v))
			a = append(a, kv)
		}
	}
	queryStr := strings.Join(a, "&")
	r.URL.RawQuery = queryStr
	return queryStr
}

func CanonicalHeaders(r *http.Request, signerHeaders []string) string {
	var a []string
	header := make(map[string][]string)
	for k, v := range r.Header {
		header[strings.ToLower(k)] = v
	}
	for _, key := range signerHeaders {
		value := header[key]
		if strings.EqualFold(key, HeaderHost) {
			value = []string{r.Host}
		}
		sort.Strings(value)
		for _, v := range value {
			a = append(a, key+":"+strings.TrimSpace(v))
		}
	}
	return fmt.Sprintf("%s\n", strings.Join(a, "\n"))
}

func SignedHeaders(r *http.Request) []string {
	var a []string
	for key := range r.Header {
		a = append(a, strings.ToLower(key))
	}
	sort.Strings(a)
	return a
}

func RequestPayload(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return []byte(""), nil
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []byte(""), err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	return b, err
}

func HexEncodeMd5Hash(body []byte) (string, error) {
	hash := md5.New()
	if body == nil {
		body = []byte("")
	}
	_, err := hash.Write(body)
	return fmt.Sprintf("%x", hash.Sum(nil)), err
}

type Signer struct {
	AccessKeyID     string
	SecretAccessKey string
}
