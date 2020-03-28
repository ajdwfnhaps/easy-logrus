package logger

import "fmt"

//Logger 日志默认实现
type Logger struct {
}

func (l *Logger) logf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

//Fatalf Fatalf
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logf(format, args...)
}

//Errorf Errorf
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logf(format, args...)
}

//Warnf Warnf
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logf(format, args...)
}

//Infof Infof
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logf(format, args...)
}

//Printf Printf
func (l *Logger) Printf(format string, args ...interface{}) {
	l.logf(format, args...)
}

//Debugf Debugf
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logf(format, args...)
}

func (l *Logger) log(args ...interface{}) {
	fmt.Println(args...)
}

//Trace Trace
func (l *Logger) Trace(args ...interface{}) {
	l.log(args...)
}

//Debug Debug
func (l *Logger) Debug(args ...interface{}) {
	l.log(args...)
}

//Print Print
func (l *Logger) Print(args ...interface{}) {
	l.log(args...)
}

//Info Info
func (l *Logger) Info(args ...interface{}) {
	l.log(args...)
}

//Warn Warn
func (l *Logger) Warn(args ...interface{}) {
	l.log(args...)
}

//Error Error
func (l *Logger) Error(args ...interface{}) {
	l.log(args...)
}

//Fatal Fatal
func (l *Logger) Fatal(args ...interface{}) {
	l.log(args...)
}
