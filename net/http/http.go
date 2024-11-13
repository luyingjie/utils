package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// HTTPClient 定义HTTP客户端结构体
type HTTPClient struct {
	client  *http.Client
	timeout time.Duration
}

// NewHTTPClient 创建新的HTTP客户端实例
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// Request 通用请求方法
func (c *HTTPClient) Request(ctx context.Context, method, url string, body interface{}, headers map[string]string) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}

	// 设置默认请求头
	req.Header.Set("Content-Type", "application/json")

	// 设置自定义请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Get 发送GET请求
func (c *HTTPClient) Get(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	return c.Request(ctx, http.MethodGet, url, nil, headers)
}

// Post 发送POST请求
func (c *HTTPClient) Post(ctx context.Context, url string, body interface{}, headers map[string]string) ([]byte, error) {
	return c.Request(ctx, http.MethodPost, url, body, headers)
}

// Put 发送PUT请求
func (c *HTTPClient) Put(ctx context.Context, url string, body interface{}, headers map[string]string) ([]byte, error) {
	return c.Request(ctx, http.MethodPut, url, body, headers)
}

// Delete 发送DELETE请求
func (c *HTTPClient) Delete(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	return c.Request(ctx, http.MethodDelete, url, nil, headers)
}
