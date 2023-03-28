package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"strings"
	qcutil "utils/qingcloud"
)

func setResponse(body []byte, request *interface{}, resp *http.Response) {
	err := json.Unmarshal(body, &request)
	if err != nil {
		if len(body) != 0 {
			*request = string(body)
		} else {
			*request = resp.Status
		}
	}
}

var ContextType string = "application/json;charset=utf-8"

func ToForm(params map[string]interface{}, urlencoded bool) string {
	parts := []string{}
	for k, v := range params {
		_v := v.(string)
		if urlencoded {
			_v = qcutil.QueryEscape(_v)
		}
		parts = append(parts, k+"="+_v)
	}
	return strings.Join(parts, "&")
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

	setResponse(body, request, resp)
	return nil
}

// FilePost 文件处理的Post，用于下载
func FilePost(url, data string, header ...map[string]string) (map[string]string, []byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return nil, nil, err
	}

	if len(header) > 0 && header[0] != nil {
		for key, value := range header[0] {
			req.Header.Add(key, value)
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	contentDisposition := res.Header.Get("Content-Disposition")
	contentType := res.Header.Get("Content-Type")
	head := map[string]string{
		"Content-Disposition": contentDisposition,
		"Content-Type":        contentType,
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	return head, body, nil
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

	setResponse(body, request, resp)
	return nil
}

// Get2 : Get的方式获取数据
func Get2(url string, request *interface{}, header ...map[string]string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if len(header) > 0 && header[0] != nil {
		for key, value := range header[0] {
			req.Header.Add(key, value)
		}
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

	setResponse(body, request, res)
	return nil
}

// GetBody 直接返回整个返回体。
func GetBody(url string, header ...map[string]string) (string, http.Header, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, err
	}

	if len(header) > 0 && header[0] != nil {
		for key, value := range header[0] {
			req.Header.Add(key, value)
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil, err
	}

	return string(body), res.Header, nil
}

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

	// json.Unmarshal(body, request)
	setResponse(body, request, resp)
	return nil
}

// Post2 : Post方式的数据获取
func Post2(url, data string, request *interface{}, header ...map[string]string) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return err
	}

	if len(header) > 0 && header[0] != nil {
		for key, value := range header[0] {
			req.Header.Add(key, value)
		}
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

	setResponse(body, request, res)
	return nil
}

// PostBody : 直接返回请求返回的body
func PostBody(url, data string, header ...map[string]string) (string, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return "", err
	}

	if len(header) > 0 && header[0] != nil {
		for key, value := range header[0] {
			req.Header.Add(key, value)
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Put : Put方式的数据获取
func Put(url, data string, request *interface{}, header ...map[string]string) error {
	req, err := http.NewRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		return err
	}

	if len(header) > 0 && header[0] != nil {
		for key, value := range header[0] {
			req.Header.Add(key, value)
		}
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

	setResponse(body, request, res)
	return nil
}

// Delete : Delete的方式获取数据
func Delete(url string, request *interface{}, header ...map[string]string) error {
	// req, err := http.NewRequest("DELETE", url, nil)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	if len(header) > 0 && header[0] != nil {
		for key, value := range header[0] {
			req.Header.Add(key, value)
		}
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

	setResponse(body, request, res)
	return nil
}

func TLSGet(url string, request *interface{}, header ...map[string]string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if len(header) > 0 && header[0] != nil {
		for key, value := range header[0] {
			req.Header.Add(key, value)
		}
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

	setResponse(body, request, res)
	return nil
}

// TLSGet : Get提交数据，参数是map，返回interface{}
func TLSGet2(url string, request *interface{}) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	setResponse(body, request, res)
	return nil
}

// TLSPost2 : https的post请求。
func TLSPost(url, data string, request *interface{}, header ...map[string]string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return err
	}

	if len(header) > 0 && header[0] != nil {
		for key, value := range header[0] {
			req.Header.Add(key, value)
		}
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

	setResponse(body, request, res)
	return nil
}

// TLSPost : Post提交数据。
func TLSPost2(url, data string, request *interface{}) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Post(url, ContextType, strings.NewReader(data))
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	setResponse(body, request, res)
	return nil
}

// TLSPut : Put方式的数据获取
func TLSPut(url, data string, request *interface{}, header ...map[string]string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		return err
	}

	if len(header) > 0 && header[0] != nil {
		for key, value := range header[0] {
			req.Header.Add(key, value)
		}
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

	setResponse(body, request, res)
	return nil
}

// TLSDelete : Delete的方式获取数据
func TLSDelete(url, data string, request *interface{}, header ...map[string]string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("DELETE", url, strings.NewReader(data))
	if err != nil {
		return err
	}

	if len(header) > 0 && header[0] != nil {
		for key, value := range header[0] {
			req.Header.Add(key, value)
		}
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

	setResponse(body, request, res)
	return nil
}

// TLSDelete : Delete的方式获取数据
func TLSDelete2(url string, request *interface{}, header ...map[string]string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	if len(header) > 0 && header[0] != nil {
		for key, value := range header[0] {
			req.Header.Add(key, value)
		}
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

	setResponse(body, request, res)
	return nil
}

// ReverseProxy Http的反向代理, 使用基础包的ReverseProxy。
func ReverseProxy(_url string, rw http.ResponseWriter, req *http.Request, resFunc func(*http.Response), errFunc func(http.ResponseWriter, *http.Request, error)) {
	u, _ := url.Parse(_url)
	proxy := httputil.NewSingleHostReverseProxy(u)

	if resFunc != nil {
		proxy.ModifyResponse = func(res *http.Response) error {
			resFunc(res)
			for key, value := range res.Header {
				for _, v := range value {
					rw.Header().Add(key, v)
				}
			}
			rw.WriteHeader(res.StatusCode)

			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}
			rw.Write(body)
			return nil
		}
	}
	if errFunc != nil {
		proxy.ErrorHandler = errFunc
	}

	proxy.ServeHTTP(rw, req)
}

// Proxy Http的反向代理， 使用http包自定义逻辑。
func Proxy(host string, rw http.ResponseWriter, req *http.Request, resFunc func(*http.Response), errFunc func(http.ResponseWriter, *http.Request, error)) {
	outreq, err := http.NewRequest(req.Method, host+req.RequestURI, nil)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}

	outreq.Header = req.Header
	outreq.Body = req.Body

	client := http.DefaultClient

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	res, err := client.Do(outreq)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}

	if resFunc != nil {
		resFunc(res)
	}

	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}
	// c.Writer.Header().Add("Test", "0")
	rw.WriteHeader(res.StatusCode)
	// io.Copy(c.Writer, res.Body)
	// res.Body.Close()
	// c.Writer.WriteHeaderNow()
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}
	rw.Write(body)
}

// ProxyRedirect Http的反向代理， 使用http包自定义逻辑。
func ProxyRedirect(host string, rw http.ResponseWriter, req *http.Request, resFunc func(*http.Response), errFunc func(http.ResponseWriter, *http.Request, error)) {
	outreq, err := http.NewRequest(req.Method, host+req.RequestURI, nil)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}

	outreq.Header = req.Header
	outreq.Body = req.Body

	client := http.DefaultClient

	res, err := client.Do(outreq)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}

	if resFunc != nil {
		resFunc(res)
	}

	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}
	// c.Writer.Header().Add("Test", "0")
	rw.WriteHeader(res.StatusCode)
	// io.Copy(c.Writer, res.Body)
	// res.Body.Close()
	// c.Writer.WriteHeaderNow()
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}
	rw.Write(body)
}

func ProxyCookieRedirect(host string, rw http.ResponseWriter, req *http.Request, resFunc func(*http.Response), errFunc func(http.ResponseWriter, *http.Request, error)) {
	jar, err := cookiejar.New(nil)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}
	client := &http.Client{Jar: jar}
	outreq, err := http.NewRequest(req.Method, host+req.RequestURI, nil)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}

	outreq.Header = req.Header
	outreq.Body = req.Body

	res, err := client.Do(outreq)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}

	// 处理cookiejar
	for _, v := range jar.Cookies(outreq.URL) { //res.Request.URL
		http.SetCookie(rw, v)
	}

	if resFunc != nil {
		resFunc(res)
	}

	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}
	// c.Writer.Header().Add("Test", "0")
	rw.WriteHeader(res.StatusCode)
	// io.Copy(c.Writer, res.Body)
	// res.Body.Close()
	// c.Writer.WriteHeaderNow()
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}
	rw.Write(body)
}

// 测试用
func ProxyTest(host string, rw http.ResponseWriter, req *http.Request, resFunc func(*http.Response), errFunc func(http.ResponseWriter, *http.Request, error), isRedirect ...bool) {
	outreq, err := http.NewRequest(req.Method, host+req.RequestURI, nil)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}

	outreq.Header = req.Header
	outreq.Body = req.Body

	client := http.DefaultClient
	_isRedirect := false
	if len(isRedirect) > 0 {
		_isRedirect = isRedirect[0]
	}

	if !_isRedirect {
		fmt.Println("执行1")
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			fmt.Println("执行2")
			return http.ErrUseLastResponse
		}
	}

	res, err := client.Do(outreq)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}

	if resFunc != nil {
		resFunc(res)
	}

	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}
	// c.Writer.Header().Add("Test", "0")
	rw.WriteHeader(res.StatusCode)
	// io.Copy(c.Writer, res.Body)
	// res.Body.Close()
	// c.Writer.WriteHeaderNow()
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}
	rw.Write(body)
}
