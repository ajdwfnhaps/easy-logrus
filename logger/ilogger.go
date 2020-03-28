package logger

//ILogger 日志接口
type ILogger interface {
	//CreateLogger() ILogger

	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Debugf(format string, args ...interface{})

	Trace(args ...interface{})
	Debug(args ...interface{})
	Print(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}
