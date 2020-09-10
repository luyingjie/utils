package error

import (
	"errors"
	"time"

	"utils/log"
)

// 错误码大类
// 5000		系统错误
// 4000		权限等错误
// 3000		资源错误
// 2000		连接错误
// 1000		参数错误

// 0(Debug), 1(Info), 2(Warn), 3(Error), 4(Panic), 5(Fatal)
type ErrorModel struct {
	Code  int
	Leve  int
	Model string
	Error error
}

var _log *log.Logger

func init() {
	_log = log.NewLog(nil, "LOGS_MODULE", 16)
	_log.ResetFlags(log.BitDate | log.BitLongFile | log.BitLevel)

	//添加标记位
	_log.AddFlag(log.BitShortFile | log.BitTime)

	//设置日志写入文件
	fileTime := time.Now().Format("2006-01-02")
	_log.SetLogFile("./logs/system", fileTime+".log")
}

// OpenDebug : 打开Debug调试
func OpenDebug() {
	_log.OpenDebug()
}

// CloseDebug : 关闭Debug调试
func CloseDebug() {
	_log.CloseDebug()
}

// Try : 处理异常
func Try(code, leve int, model string, err error) {
	errModel := ErrorModel{
		Code:  code,
		Leve:  leve,
		Model: model,
		Error: err,
	}

	go writing(errModel)
	panic(errModel)
}

// Trys : 处理string的异常
func Trys(code, leve int, model, str string) {
	errModel := ErrorModel{
		Code:  code,
		Leve:  leve,
		Model: model,
		Error: errors.New(str),
	}

	go writing(errModel)
	panic(errModel)
}

// Log : 写入日志
func Log(code, leve int, model string, err error) {
	errModel := ErrorModel{
		Code:  code,
		Leve:  leve,
		Model: model,
		Error: err,
	}

	go writing(errModel)
}

// Leve : 0(Debug), 1(Info), 2(Warn), 3(Error), 4(Panic), 5(Fatal)
func writing(err ErrorModel) {
	_log.SetPrefix(string(err.Code) + ":" + err.Model)

	switch err.Leve {
	case 0:
		_log.Debug(err.Error)
	case 1:
		_log.Info(err.Error)
	case 2:
		_log.Warn(err.Error)
	case 3:
		_log.Error(err.Error)
	}
}
