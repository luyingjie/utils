package verror

import (
	"testing"
	"utils/log"
)

func TestLog(t *testing.T) {
	_log := log.NewLog(nil, "LOGS_MODULE", log.BitShortFile|log.BitTime)
	_log.ResetFlags(log.BitDate | log.BitLongFile | log.BitLevel)

	//设置日志前缀，主要标记当前日志模块
	_log.SetPrefix("1200:MODULE")

	//添加标记位
	_log.AddFlag(log.BitShortFile | log.BitTime)

	// _log.Stack(" vic Stack! ")

	//设置日志写入文件
	_log.SetLogFile("./logs", "testfile.log")
	_log.Debug("===> vic debug content ~~666")
	_log.Debug("===> vic debug content ~~888")
	_log.Error("===> vic Error!!!! ~~~555~~~")
}

func TestLog2(t *testing.T) {
	_log := log.NewLog(nil, "LOGS_MODULE", log.BitShortFile|log.BitTime)
	_log.ResetFlags(log.BitDate | log.BitLongFile | log.BitLevel)

	//设置日志前缀，主要标记当前日志模块
	// _log.SetPrefix("MODULE")

	//添加标记位
	_log.AddFlag(log.BitShortFile | log.BitTime)
	// _log.Stack(" vic Stack! ")

	//设置日志写入文件
	_log.SetLogFile("./logs", "testfile2.log")
	_log.Debug("===> vic debug content ~~666")
	_log.Debug("===> vic debug content ~~888")
	_log.Error("===> vic Error!!!! ~~~555~~~")
}
