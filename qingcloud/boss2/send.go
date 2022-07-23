package boss2

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"utils/conv"
	vhttp "utils/net/http"
	verror "utils/os/error"
	qcutil "utils/qingcloud"
	"utils/util"
)

// Send 发送请求到boss2
// conf 包含配置：console_key_id,console_secrect_key,console_uri,host,port
func Send(params map[string]interface{}, conf map[string]interface{}, uriKey ...string) (interface{}, error) {
	_uriKey := conf["console_uri"].(string)
	if len(uriKey) > 0 && uriKey[0] != "" {
		_uriKey = conf[uriKey[0]].(string)
	}
	_, _, data, err := Signature(_uriKey, conf["console_key_id"].(string), conf["console_secrect_key"].(string), params)
	if err != nil {
		return nil, err
	}

	port := conf["port"].(string)
	url := fmt.Sprintf(conf["protocol"].(string)+"://%s:%s%s", conf["host"].(string), port, _uriKey)

	dataStr := vhttp.ToForm(data, true)

	headers := map[string]string{}
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["Accept"] = "*/*"
	headers["Connection"] = "Keep-Alive"
	headers["Content-Length"] = conv.String(len(dataStr))

	var resp interface{}
	if conf["protocol"].(string) == "https" {
		vhttp.TLSPost(url, dataStr, &resp, headers)
	} else {
		vhttp.Post2(url, dataStr, &resp, headers)
	}

	return resp, nil
}

func Signature(uri, ak, sk string, params map[string]interface{}) (string, string, map[string]interface{}, error) {
	bData, err := json.Marshal(params)
	if err != nil {
		return "", "", nil, verror.New("parameter parsing error")
	}

	return SignatureStr(uri, ak, sk, params["action"].(string), string(bData))
}

func SignatureStr(uri, ak, sk, action, params string) (string, string, map[string]interface{}, error) {
	_params := map[string]interface{}{
		"access_key_id": ak,
		"action":        action,
		"params":        params,
		"time_stamp":    util.TimeToString(time.Now(), "ISO 8601"),
	}

	// 参数是固定的，其实都不需要排序。
	keys := []string{}
	for key := range _params {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	parts := []string{}
	for _, key := range keys {
		v := _params[key]
		if v != nil {
			_v := ""
			switch reflect.TypeOf(v).String() {
			case "string":
				_v = qcutil.QueryEscape(v.(string))
				parts = append(parts, key+"="+_v)
			}
		}
	}
	urlParams := strings.Join(parts, "&")
	signature := qcutil.Get_boss2_authorization(sk, "POST", uri, urlParams)
	_params["signature"] = signature

	return urlParams, signature, _params, nil
}
