package http

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"utils/error"
)

const ContextType string = "application/json;charset=utf-8"

// Post : Post提交和获取Json数据
func Post(url, data string, request *interface{}, header map[string]string) {
	resp, err := http.Post(url, ContextType, strings.NewReader(data))
	if err != nil {
		error.Try(2000, 3, "utils/net/http/http/Post/Post", err)
	}

	// resp.Header.Add()
	for key, value := range header {
		resp.Header.Add(key, value)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		error.Try(5000, 3, "utils/net/http/http/Post/ReadAll", err)
	}

	json.Unmarshal(body, request)
}

// PostToMap : Post提交数据，参数是map，返回interface{}
func PostToMap(url string, data map[string]interface{}, request *interface{}, header map[string]string) {
	d, err := json.Marshal(data)
	if err != nil {
		error.Try(5000, 3, "utils/net/http/http/PostToMap/Marshal", err)
	}
	resp, err := http.Post(url, ContextType, strings.NewReader(string(d)))
	if err != nil {
		error.Try(2000, 3, "utils/net/http/http/PostToMap/Post", err)
	}

	// resp.Header.Add()
	for key, value := range header {
		resp.Header.Add(key, value)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		error.Try(5000, 3, "utils/net/http/http/PostToMap/ReadAll", err)
	}

	json.Unmarshal(body, request)
}

// Get : Get方式提交数据
func Get(url string, request *interface{}, header map[string]string) {
	resp, err := http.Get(url)
	if err != nil {
		error.Try(2000, 3, "utils/net/http/http/Get", err)
	}

	// resp.Header.Add()
	for key, value := range header {
		resp.Header.Add(key, value)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		error.Try(5000, 3, "utils/net/http/http/Get/ReadAll", err)
	}

	json.Unmarshal(body, request)
}

// Put : Put方式的数据获取
func Put(url string, data interface{}, request *interface{}, header map[string]string) {
	d, err := json.Marshal(data)
	if err != nil {
		error.Try(5000, 3, "utils/net/http/http/Put/Marshal", err)
	}
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(d)))
	if err != nil {
		error.Try(2000, 3, "utils/net/http/http/Put/NewRequest", err)
	}

	// resp.Header.Add()
	for key, value := range header {
		req.Header.Add(key, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		error.Try(2000, 3, "utils/net/http/http/Put/Do", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		error.Try(5000, 3, "utils/net/http/http/Put/ReadAll", err)
	}

	json.Unmarshal(body, request)
}

// Delete : Delete的方式获取数据
func Delete(url, data string, request *interface{}, header map[string]string) {
	// req, err := http.NewRequest("DELETE", url, nil)
	req, err := http.NewRequest("DELETE", url, strings.NewReader(string(data)))
	if err != nil {
		error.Try(2000, 3, "utils/net/http/http/Delete/NewRequest", err)
	}

	// resp.Header.Add()
	for key, value := range header {
		req.Header.Add(key, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		error.Try(2000, 3, "utils/net/http/http/Delete/Do", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		error.Try(5000, 3, "utils/net/http/http/Delete/ReadAll", err)
	}

	json.Unmarshal(body, request)
}

// TLSGet : Get提交数据，参数是map，返回interface{}
func TLSGet(url string, request *interface{}, header map[string]string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		error.Try(2000, 3, "utils/net/http/http/TLSGet/Get", err)
	}

	// resp.Header.Add()
	for key, value := range header {
		resp.Header.Add(key, value)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		error.Try(5000, 3, "utils/net/http/http/TLSGet/ReadAll", err)
	}

	json.Unmarshal(body, request)
}

// TLSPost2 : https的post请求。
func TLSPost(url, data string, request *interface{}, header map[string]string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		error.Try(2000, 3, "utils/net/http/http/TLSPost/NewRequest", err)
	}

	// resp.Header.Add()
	for key, value := range header {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		error.Try(2000, 3, "utils/net/http/http/TLSPost/Do", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		error.Try(5000, 3, "utils/net/http/http/TLSPost/ReadAll", err)
	}

	json.Unmarshal(body, request)
}

// TLSPost : Post提交数据，参数是map，返回interface{} ，使用client.Post 这个方法目前有问题。
func TLSPost2(url, data string, request *interface{}, header map[string]string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Post(url, ContextType, strings.NewReader(data))
	if err != nil {
		error.Try(2000, 3, "utils/net/http/http/TLSPost2/Post", err)
	}

	// resp.Header.Add()
	for key, value := range header {
		resp.Header.Add(key, value)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		error.Try(5000, 3, "utils/net/http/http/TLSPost2/ReadAll", err)
	}

	json.Unmarshal(body, request)
}

// TLSPut : Put方式的数据获取
func TLSPut(url, data string, request *interface{}, header map[string]string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		error.Try(2000, 3, "utils/net/http/http/TLSPut/NewRequest", err)
	}

	// resp.Header.Add()
	for key, value := range header {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		error.Try(2000, 3, "utils/net/http/http/TLSPut/Do", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		error.Try(5000, 3, "utils/net/http/http/TLSPut/ReadAll", err)
	}

	json.Unmarshal(body, request)
}

// TLSDelete : Delete的方式获取数据
func TLSDelete(url, data string, request *interface{}, header map[string]string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("DELETE", url, strings.NewReader(data))
	if err != nil {
		error.Try(2000, 3, "utils/net/http/http/TLSDelete/NewRequest", err)
	}

	// resp.Header.Add()
	for key, value := range header {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		error.Try(2000, 3, "utils/net/http/http/TLSDelete/Do", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		error.Try(5000, 3, "utils/net/http/http/TLSDelete/ReadAll", err)
	}

	json.Unmarshal(body, request)
}
