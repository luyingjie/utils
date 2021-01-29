package error

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"utils/os/log"

	"utils/os/error/base"
)

var _log *log.Logger

// 加载框架的日志
func loadLog() {
	_log = log.NewLog(nil, "LOGS_MODULE", 16)
	_log.ResetFlags(log.BitDate | log.BitLongFile | log.BitLevel)

	//添加标记位
	_log.AddFlag(log.BitShortFile | log.BitTime)

	//设置日志写入文件
	fileTime := time.Now().Format("2006-01-02")
	_log.SetLogFile("./logs/system", fileTime+".log")
}

// 写入日志
func writing(err ErrorModel) {
	_log.SetPrefix(strconv.Itoa(err.Code))

	switch err.Level {
	case 0:
		_log.Debug(fmt.Sprintf("%+v", err.Error))
	case 1:
		_log.Info(fmt.Sprintf("%+v", err.Error))
	case 2:
		_log.Warn(fmt.Sprintf("%+v", err.Error))
	case 3:
		_log.Error(fmt.Sprintf("%+v", err.Error))
	}
}

func init() {
	loadLog()
}

// OpenDebug : 打开Debug调试
func OpenDebug() {
	_log.OpenDebug()
}

// CloseDebug : 关闭Debug调试
func CloseDebug() {
	_log.CloseDebug()
}

// New 创建一个错误对象
func New(text string) error {
	return errors.New(text)
}

// Newf 创建一个格式化错误对象
func Newf(format string, args ...interface{}) error {
	if format == "" {
		return nil
	}
	return errors.New(fmt.Sprintf(format, args...))
}

// Try 处理异常。Leve : 0(Debug), 1(Info), 2(Warn), 3(Error), 4(Panic), 5(Fatal)
func Try(code, level int, err error) {
	errModel := ErrorModel{
		Code:  code,
		Level: level,
		// Model: model,
		Error: base.Wrap(err, ""),
	}
	errModel.SetMessage()

	writing(errModel)
	panic(errModel)
}

// TryText 处理string的异常。Level : 0(Debug), 1(Info), 2(Warn), 3(Error), 4(Panic), 5(Fatal)
func TryText(code, level int, str string) {
	errModel := ErrorModel{
		Code:  code,
		Level: level,
		// Model: model,
		Error: base.New(str),
	}
	errModel.SetMessage()

	writing(errModel)
	panic(errModel)
}

// Error 处理异常。Leve : 0(Debug), 1(Info), 2(Warn), 3(Error), 4(Panic), 5(Fatal)
func Error(code, level int, err error) ErrorModel {
	errModel := ErrorModel{
		Code:  code,
		Level: level,
		// Model: model,
		Error: base.Wrap(err, ""),
	}
	errModel.SetMessage()

	writing(errModel)
	return errModel
}

// ErrorText 处理string的异常。Level : 0(Debug), 1(Info), 2(Warn), 3(Error), 4(Panic), 5(Fatal)
func ErrorText(code, level int, str string) ErrorModel {
	errModel := ErrorModel{
		Code:  code,
		Level: level,
		// Model: model,
		Error: base.New(str),
	}
	errModel.SetMessage()

	writing(errModel)
	return errModel
}

// Log 写入日志
func Log(code, level int, err error) {
	errModel := ErrorModel{
		Code:  code,
		Level: level,
		// Model: model,
		Error: base.Wrap(err, ""),
	}
	errModel.SetMessage()

	writing(errModel)
}

// LogText 写入日志
func LogText(code, level int, str string) {
	errModel := ErrorModel{
		Code:  code,
		Level: level,
		// Model: model,
		Error: base.New(str),
	}
	errModel.SetMessage()

	writing(errModel)
}
