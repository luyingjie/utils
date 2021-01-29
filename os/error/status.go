package error

// ErrorStatus 错误编码
var ErrorStatus map[int]string = map[int]string{
	200:  "成功",
	400:  "错误请求",
	401:  "未认证，请先登录",
	403:  "未授权访问",
	404:  "页面不存在",
	405:  "不支持该方法",
	429:  "请勿频繁请求",
	498:  "客户端取消请求",
	500:  "网络错误，请稍后重试",
	503:  "过载保护,服务暂不可用",
	504:  "服务调用超时",
	600:  "应用程序不存在或已被封禁",
	601:  "签名校验失败",
	602:  "重复请求",
	603:  "验证码错误",
	604:  "资源锁定中，请稍后重试",
	605:  "请求体大小超出限制",
	606:  "系统升级中",
	1000: "参数错误",
	2000: "连接错误",
	3000: "资源错误",
	4000: "权限等错误",
	5000: "系统错误",
}

// 错误模型

// ErrorModel Code 0(Debug), 1(Info), 2(Warn), 3(Error), 4(Panic), 5(Fatal)
type ErrorModel struct {
	Code    int
	Level   int
	Model   string
	Message string
	Error   error
}

// SetMessage 设置模型的预定错误信息
func (m *ErrorModel) SetMessage() {
	if m != nil {
		m.Message = ErrorStatus[m.Code]
	} else {
		m.Message = "未知错误"
	}
}

// AddStatus 外部追加错误信息的方法, 如果又追加就替换原始值。 为了性能使用乐观锁，一般在主程序跑起来的时候会先追加，后面不会在改变。
func AddStatus(status map[int]string) {
	for key, value := range status {
		ErrorStatus[key] = value
	}
}
