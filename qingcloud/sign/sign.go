// QingCloud KMS API Gateway Signature
// based on https://github.com/datastream/aws/blob/master/signv4.go

package sign

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/luyingjie/utils/conv"
)

func SignRequest(r *http.Request, keyid string, secret string, params ...map[string]interface{}) error {
	s := Signer{
		AccessKeyID:     keyid,
		SecretAccessKey: secret,
	}
	hf := NewHmacSHA256([]byte(s.SecretAccessKey))
	var err error
	if len(params) > 0 && params[0] != nil {
		_, err = s.Sign(r, hf, params[0])
	}
	_, err = s.Sign(r, hf)
	return err
}

func (s *Signer) Sign(r *http.Request, hf HmacFunc, params ...map[string]interface{}) (string, error) {
	q := r.URL.Query()
	if len(params) > 0 && params[0] != nil {
		for key, value := range params[0] {
			_v := ""
			if val, ok := value.(int); ok {
				_v = conv.String(val)
			} else if val, ok := value.(string); ok {
				_v = val
			} else if val, ok := value.(bool); ok {
				_v = conv.String(val)
			}
			q.Add(key, _v)
		}
	}
	// if q.Get("version") == "" {
	// 	q.Add("version", "1")
	// }
	if q.Get("signature_version") == "" {
		q.Add("signature_version", "1")
	}
	if q.Get("signature_method") == "" {
		q.Add("signature_method", hf.GetName())
	}
	if q.Get("access_key_id") == "" {
		q.Add("access_key_id", s.AccessKeyID)
	}
	if q.Get("timestamp") == "" {
		t := time.Now().Add(time.Hour)
		q.Add("timestamp", t.Format("2006-01-02T15:04:05Z"))
	}
	if q.Get("time_stamp") == "" {
		t := time.Now().Add(time.Hour)
		q.Add("time_stamp", t.Format("2006-01-02T15:04:05Z")) // time.Now().UTC().Format(time.RFC3339)
	}
	r.URL.RawQuery = q.Encode()

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

	q = r.URL.Query()
	q.Set("signature", signature)
	r.URL.RawQuery = q.Encode()

	return signature, nil
}
