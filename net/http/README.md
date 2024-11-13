使用示例：

```go
package main

import (
	"context"
	"fmt"
	"time"
	"your/project/pkg/httpclient"
)

func main() {
	// 创建客户端实例，设置5秒超时
	client := httpclient.NewHTTPClient(5 * time.Second)

	// 设置请求头
	headers := map[string]string{
		"Authorization": "Bearer your-token",
	}

	// GET请求示例
	resp, err := client.Get(context.Background(), "https://api.example.com/users", headers)
	if err != nil {
		fmt.Printf("GET请求失败: %v\n", err)
		return
	}
	fmt.Printf("GET响应: %s\n", string(resp))

	// POST请求示例
	requestBody := map[string]interface{}{
		"name": "张三",
		"age":  25,
	}
	
	resp, err = client.Post(context.Background(), "https://api.example.com/users", requestBody, headers)
	if err != nil {
		fmt.Printf("POST请求失败: %v\n", err)
		return
	}
	fmt.Printf("POST响应: %s\n", string(resp))
}
```

这个HTTP客户端封装提供以下特性：
1. 支持超时设置
2. 支持上下文（Context）控制
3. 支持自定义请求头
4. 自动处理JSON序列化
5. 支持GET、POST、PUT、DELETE等常用HTTP方法
6. 统一的错误处理
7. 自动关闭响应体

使用这个封装可以大大简化HTTP请求的处理过程。你可以根据实际需求进行进一步的扩展，比如：
- 添加重试机制
- 添加请求和响应的拦截器
- 添加更多的请求方法
- 支持更多的请求体格式
- 添加响应状态码的处理
- 添加请求日志记录


需要注意的是，这个封装默认使用 JSON 作为请求和响应的格式。如果你需要支持其他格式（如 XML、Form 等），可以进行相应的扩展。