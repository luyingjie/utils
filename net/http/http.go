package http

import (
	"crypto/tls"
	"encoding/json"
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

// Proxy Http的反向代理, 使用基础包的ReverseProxy。
func Proxy(_url string, rw http.ResponseWriter, req *http.Request, resFunc func(*http.Response) error, errFunc func(http.ResponseWriter, *http.Request, error)) {
	u, _ := url.Parse(_url)
	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	proxy := httputil.NewSingleHostReverseProxy(u)
	// proxy.Transport = tr
	if resFunc != nil {
		proxy.ModifyResponse = resFunc
	}
	if errFunc != nil {
		proxy.ErrorHandler = errFunc
	}

	proxy.ServeHTTP(rw, req)
}

// 请尽量使用proxy3, 因为项目使用，所以保存了该方法。不会在更新。大部分场景可用，有前端代理，特别是登录和重定向处理要小心使用。
// Proxy2 Http的反向代理， 使用http包自定义逻辑。
func Proxy2(_url string, rw http.ResponseWriter, req *http.Request, resFunc func(*http.Response) error, errFunc func(http.ResponseWriter, *http.Request, error)) {
	_req, err := http.NewRequest(req.Method, _url+req.RequestURI, nil)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}

	_req.Header = req.Header
	_req.Body = req.Body
	// 这里可以考虑是否直接用传进来的req作为请求参数， 还需要测试还支撑。
	res, err := http.DefaultClient.Do(_req)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}

	if resFunc != nil {
		err := resFunc(res)
		if err != nil && errFunc != nil {
			errFunc(rw, req, err)
		}
	} else {
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
}

func TLSProxy2(_url string, rw http.ResponseWriter, req *http.Request, resFunc func(*http.Response) error, errFunc func(http.ResponseWriter, *http.Request, error)) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	_req, err := http.NewRequest(req.Method, _url+req.RequestURI, nil)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}

	_req.Header = req.Header
	_req.Body = req.Body
	// 这里可以考虑是否直接用传进来的req作为请求参数， 还需要测试还支撑。
	res, err := client.Do(_req)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}

	if resFunc != nil {
		err := resFunc(res)
		if err != nil && errFunc != nil {
			errFunc(rw, req, err)
		}
	} else {
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
}

// proxy2的改进，当时项目用到proxy2， 而且优化了参数模式,所以保持proxy2方法。
func Proxy3(_url string, rw http.ResponseWriter, req *http.Request, resFunc func(*http.Response), errFunc func(http.ResponseWriter, *http.Request, error)) {
	jar, err := cookiejar.New(nil)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}
	client := &http.Client{Jar: jar}
	_req, err := http.NewRequest(req.Method, _url+req.RequestURI, nil)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}

	// u, _ := url.Parse(_url)
	// _req.URL = u
	_req.Header = req.Header
	_req.Body = req.Body
	// 这里可以考虑是否直接用传进来的req作为请求参数， 还需要测试还支撑。
	res, err := client.Do(_req)
	if err != nil && errFunc != nil {
		errFunc(rw, req, err)
	}

	// 处理cookiejar
	for _, v := range jar.Cookies(_req.URL) { //res.Request.URL
		http.SetCookie(rw, v)
	}

	// 处理重定向
	// if res.Request.Response != nil && res.Request.Response.Request.Method == "POST" && (res.Request.Response.StatusCode == 301 || res.Request.Response.StatusCode == 302) {
	// 	// Proxy2(_url+res.Request.URL.Path, rw, res.Request, nil, nil)
	// 	// 这个方案不行，因为body的流已经关闭了。
	// 	// rw.WriteHeader(200)
	// 	// defer res.Body.Close()
	// 	// body, err := ioutil.ReadAll(res.Request.Response.Body)
	// 	// if err != nil && errFunc != nil {
	// 	// 	errFunc(rw, req, err)
	// 	// }
	// 	// rw.Write(body)
	// 	return
	// }

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
