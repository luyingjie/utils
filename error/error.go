package error

import (
	"errors"
	"time"

	"utils/log"
)

var _log *log.Logger

func init() {
	_log = log.NewLog(nil, "LOGS_MODULE", 16)
	_log.ResetFlags(log.BitDate | log.BitLongFile | log.BitLevel)

	//添加标记位
	_log.AddFlag(log.BitShortFile | log.BitTime)
	// _log.Stack(" vic Stack! ")

	//设置日志写入文件
	fileTime := time.Now().Format("2006-01-02")
	_log.SetLogFile("./logs/system", fileTime+".log")
}

func TryWarning(err error) {
	if err != nil {
		_log.Warn(err)
		panic(err)
	}
}

func TryWarningString(err string) {
	if err != "" {
		TryWarning(errors.New(err))
	}
}

//统一处理报错
func TryError(err error) {
	if err != nil {
		_log.Error(err)
		panic(err)
	}
}

func TryErrorString(err string) {
	if err != "" {
		TryError(errors.New(err))
	}
}

// 参数错误可以直接传入模型和key，自动分类和返回相对字段文本。
func TryParameterError(err error) {
	if err != nil {
		_log.Warn(err)
		panic(err)
	}
}

func TryParameterErrorString(err string) {
	TryParameterError(errors.New(err))
}

func TryLog(err string) {
	if err != "" {
		_log.Error(err)
	}
}
