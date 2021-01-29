package http

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

var ContextType string = "application/json;charset=utf-8"

// Post : Post提交和获取Json数据
func Post(url, data string, request *interface{}) error {
	resp, err := http.Post(url, ContextType, strings.NewReader(data))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, request)
	return nil
}

// PostToMap : Post提交数据，参数是map，返回interface{}
func PostToMap(url string, data map[string]interface{}, request *interface{}) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, ContextType, strings.NewReader(string(d)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, request)
	return nil
}

// Get : Get方式提交数据
func Get(url string, request *interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, request)
	return nil
}

// Get2 : Get的方式获取数据
func Get2(url string, header map[string]string, request *interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// resp.Header.Add()
	for key, value := range header {
		req.Header.Add(key, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, request)
	return nil
}

// Post2 : Post方式的数据获取
func Post2(url, data string, header map[string]string, request *interface{}) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return err
	}

	// resp.Header.Add()
	for key, value := range header {
		req.Header.Add(key, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, request)
	return nil
}

// Put : Put方式的数据获取
func Put(url, data string, header map[string]string, request *interface{}) error {
	req, err := http.NewRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		return err
	}

	// resp.Header.Add()
	for key, value := range header {
		req.Header.Add(key, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, request)
	return nil
}

// Delete : Delete的方式获取数据
func Delete(url, data string, header map[string]string, request *interface{}) error {
	// req, err := http.NewRequest("DELETE", url, nil)
	req, err := http.NewRequest("DELETE", url, strings.NewReader(string(data)))
	if err != nil {
		return err
	}

	// resp.Header.Add()
	for key, value := range header {
		req.Header.Add(key, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, request)
	return nil
}

// TLSGet : Get提交数据，参数是map，返回interface{}
func TLSGet(url string, header map[string]string, request *interface{}) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	// resp.Header.Add()
	for key, value := range header {
		resp.Header.Add(key, value)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, request)
	return nil
}

// TLSPost2 : https的post请求。
func TLSPost(url, data string, header map[string]string, request *interface{}) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return err
	}

	// resp.Header.Add()
	for key, value := range header {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, request)
	return nil
}

// TLSPost : Post提交数据，参数是map，返回interface{} ，使用client.Post 这个方法目前有问题。
func TLSPost2(url, data string, header map[string]string, request *interface{}) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Post(url, ContextType, strings.NewReader(data))
	if err != nil {
		return err
	}

	// resp.Header.Add()
	for key, value := range header {
		resp.Header.Add(key, value)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, request)
	return nil
}

// TLSPut : Put方式的数据获取
func TLSPut(url, data string, header map[string]string, request *interface{}) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		return err
	}

	// resp.Header.Add()
	for key, value := range header {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, request)
	return nil
}

// TLSDelete : Delete的方式获取数据
func TLSDelete(url, data string, header map[string]string, request *interface{}) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("DELETE", url, strings.NewReader(data))
	if err != nil {
		return err
	}

	// resp.Header.Add()
	for key, value := range header {
		req.Header.Add(key, value)
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, request)
	return nil
}
