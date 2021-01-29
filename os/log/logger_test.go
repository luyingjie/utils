package log

import (
	"testing"
)

func TestLog(t *testing.T) {

	//测试 默认debug输出
	Debug("vic debug content1")
	Debug("vic debug content2")

	Debugf(" vic debug a = %d\n", 10)

	//设置log标记位，加上长文件名称 和 微秒 标记
	ResetFlags(BitDate | BitLongFile | BitLevel)
	Info("vic info content")

	//设置日志前缀，主要标记当前日志模块
	SetPrefix("MODULE")
	Error("vic error content")

	//添加标记位
	AddFlag(BitShortFile | BitTime)
	Stack(" vic Stack! ")

	//设置日志写入文件
	SetLogFile("./log", "testfile.log")
	Debug("===> vic debug content ~~666")
	Debug("===> vic debug content ~~888")
	Error("===> vic Error!!!! ~~~555~~~")

	//关闭debug调试
	CloseDebug()
	Debug("===> 我不应该出现~！")
	Debug("===> 我不应该出现~！")
	Error("===> vic Error  after debug close !!!!")

}
