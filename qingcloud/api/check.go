package api

import (
	"net/http"
	"net/url"
	"time"

	verror "github.com/luyingjie/utils/os/error"
	"github.com/luyingjie/utils/qingcloud/iaas"
	"github.com/luyingjie/utils/qingcloud/sign"
)

func Check(request *http.Request, access_key_id, signature, time_stamp string, data *map[string]interface{}, apiConfig map[string]interface{}) (string, error) {
	// 验证时间戳, 参数中的时间和过期时间。
	_time, err := time.Parse("2006-01-02T15:04:05Z", time_stamp)
	if err != nil {
		return "", verror.New("time_stamp Invalid format")
	}
	if _time.Before(time.Now()) {
		return "", verror.New("request out of date")
	}

	// 这里发现空就造成流程走不下去的问题，这里的判断有问题，目测是没有意义的，先去掉，留存一下逻辑。 Luke 2022.5.16
	// if data != nil && *data != nil {
	params := map[string]interface{}{
		"action":        "DescribeAccessKeys",
		"access_keys.n": access_key_id,
	}

	_resp, err := iaas.Send("GET", params, apiConfig)
	if err != nil {
		return "", err
	}

	if _resp == nil {
		return "", nil
	}

	resp_data := _resp.(map[string]interface{})
	if resp_data["ret_code"].(float64) == 0 && resp_data["total_count"].(float64) == 1 {
		_data := resp_data["access_key_set"].([]interface{})[0].(map[string]interface{})
		if _data["status"].(string) != "active" {
			return "", verror.New("user unavailable")
		}
		*data = _data
	} else {
		return "", verror.New("access_key_id error")
	}
	// }

	var req http.Request
	req.Method = request.Method
	req.Body = request.Body
	var url url.URL
	url.Path = request.URL.Path
	url.RawQuery = request.URL.RawQuery
	req.URL = &url

	signature, err1 := sign.CheckSign(&req, access_key_id, (*data)["secret_access_key"].(string), "")
	if err1 != nil {
		return "", err
	}

	return signature, nil
}
