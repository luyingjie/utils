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

// Proxy Http的反向代理
func ProxyOld(_url string, rw http.ResponseWriter, req *http.Request) {
	u, _ := url.Parse(_url)
	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	proxy := httputil.NewSingleHostReverseProxy(u)
	// proxy.Transport = tr
	proxy.ServeHTTP(rw, req)
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

			// 这个方法只是一部分，父方法里面还要调用，所以不能close。
			// defer res.Body.Close()
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
// 这里实测一个问题：isRedirect不通传参进来后日志和执行都正常，但是client.CheckRedirect这段会有干扰，可能是这个方法或者写法问题，待求证。
// isRedirect一直是true的时候正常，先传false在传true会无效，这个是按照业务流程和表象推断的测试用例。
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
