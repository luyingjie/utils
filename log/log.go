package log

/*
   全局默认提供一个Log对外句柄，可以直接使用API系列调用
   全局日志对象 vLog
*/

import (
	"os"
)

var Log = NewLog(os.Stderr, "", BitDefault)

//获取Log 标记位
func Flags() int {
	return Log.Flags()
}

//设置Log标记位
func ResetFlags(flag int) {
	Log.ResetFlags(flag)
}

//添加flag标记
func AddFlag(flag int) {
	Log.AddFlag(flag)
}

//设置Log 日志头前缀
func SetPrefix(prefix string) {
	Log.SetPrefix(prefix)
}

//设置Log绑定的日志文件
func SetLogFile(fileDir string, fileName string) {
	Log.SetLogFile(fileDir, fileName)
}

//设置关闭debug
func CloseDebug() {
	Log.CloseDebug()
}

//设置打开debug
func OpenDebug() {
	Log.OpenDebug()
}

// ====> Debug <====
func Debugf(format string, v ...interface{}) {
	Log.Debugf(format, v...)
}

func Debug(v ...interface{}) {
	Log.Debug(v...)
}

// ====> Info <====
func Infof(format string, v ...interface{}) {
	Log.Infof(format, v...)
}

func Info(v ...interface{}) {
	Log.Info(v...)
}

// ====> Warn <====
func Warnf(format string, v ...interface{}) {
	Log.Warnf(format, v...)
}

func Warn(v ...interface{}) {
	Log.Warn(v...)
}

// ====> Error <====
func Errorf(format string, v ...interface{}) {
	Log.Errorf(format, v...)
}

func Error(v ...interface{}) {
	Log.Error(v...)
}

// ====> Fatal 需要终止程序 <====
func Fatalf(format string, v ...interface{}) {
	Log.Fatalf(format, v...)
}

func Fatal(v ...interface{}) {
	Log.Fatal(v...)
}

// ====> Panic  <====
func Panicf(format string, v ...interface{}) {
	Log.Panicf(format, v...)
}

func Panic(v ...interface{}) {
	Log.Panic(v...)
}

// ====> Stack  <====
func Stack(v ...interface{}) {
	Log.Stack(v...)
}

func init() {
	//因为Log对象 对所有输出方法做了一层包裹，所以在打印调用函数的时候，比正常的logger对象多一层调用
	//一般的Logger对象 calldDepth=2, Log的calldDepth=3
	Log.CalldDepth = 3
}
