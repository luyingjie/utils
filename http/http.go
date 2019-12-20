package http

import (
	"encoding/json"
	"io/ioutil"

	"net/http"
	"strings"
	"utils/error"
)

var ContextType string = "application/json;charset=utf-8"

// Post : Post提交和获取Json数据
func Post(url, data string, request *interface{}) {
	resp, err := http.Post(url, ContextType, strings.NewReader(data))
	if err != nil {
		error.TryError(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		error.TryError(err)
	}

	json.Unmarshal(body, request)
}

// PostToMap : Post提交数据，参数是map，返回interface{}
func PostToMap(url string, data map[string]interface{}, request *interface{}) {
	d, err := json.Marshal(data)
	if err != nil {
		error.TryError(err)
	}
	resp, err := http.Post(url, ContextType, strings.NewReader(string(d)))
	if err != nil {
		error.TryError(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		error.TryError(err)
	}

	json.Unmarshal(body, request)
}

// Get : Get方式提交数据
func Get(url string, request *interface{}) {
	resp, err := http.Get(url)
	if err != nil {
		error.TryError(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		error.TryError(err)
	}

	json.Unmarshal(body, request)
}

// Put : Put方式的数据获取
func Put(url string, data interface{}, request *interface{}) {
	d, err := json.Marshal(data)
	if err != nil {
		error.TryError(err)
	}
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(d)))
	if err != nil {
		error.TryError(err)
	}

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		error.TryError(err)
	}

	json.Unmarshal(body, request)
}

// Delete : Delete的方式获取数据
func Delete(url string, request *interface{}) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		error.TryError(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		error.TryError(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		error.TryError(err)
	}

	json.Unmarshal(body, request)
}
