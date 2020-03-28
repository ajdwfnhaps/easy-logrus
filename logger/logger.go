package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ajdwfnhaps/easy-logrus/common"
	"github.com/ajdwfnhaps/easy-logrus/conf"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

// 定义键名
const (
	TraceIDKey        = "tid"
	SpanTitleKey      = "span_title"
	SpanFunctionKey   = "span_function"
	VersionKey        = "version"
	AppNameKey        = "name"
	AppNoKey          = "appno"
	TimestampKey      = "timestamp"
	FrmDeviceIDKey    = "frmDeviceID"
	ProductTypeKey    = "productType"
	ProductSubTypeKey = "productSubType"
)

type (
	traceIDContextKey struct{}

	//LogConfigFunc 设置全局配置函数
	LogConfigFunc func(c *conf.LogOption)
)

// Entry 定义统一的日志写入方式
type Entry struct {
	entry *logrus.Entry
}

var (
	//GlobalLogOption 全局日志配置
	GlobalLogOption conf.LogOption
	//日志配置文件路径
	logConfigPath string
)

//引用包时初始化函数
func init() {
	GlobalLogOption = conf.LogOption{
		AppName:                "Go应用001",
		Level:                  conf.InfoLevel,    //日志级别(1:fatal 2:error,3:warn,4:info,5:debug)
		Format:                 "json",            //日志格式（支持输出格式：text/json）
		Output:                 "multi",           //日志输出(支持：stdout/stderr/file/multi)
		OutputFile:             "logs/app",        //指定日志输出的文件路径 logs/app.log
		TIDFunc:                common.NewTraceID, //设置获取跟踪ID的函数
		DisableCustomTimestamp: false,             //是否禁用自定义时间戳显示
		DisableLineHook:        false,             //是否禁用行号信息显示(WarnLevel以上才会显示)
		LogFileMaxAge:          7,                 //设置保留7天
		LogFileRotationTime:    24 * 60 * 60,      //设置每天分割日志文件
		LogFilePathFormat:      ".%Y-%m-%d.log",   //设置日志文件名规则
	}

	//TODO reload with config file
}

//DeviceInit 手动初始化
func DeviceInit(frmDeviceID string, productType string, productSubType string) {
	GlobalLogOption.FrmDeviceID = frmDeviceID
	GlobalLogOption.ProductType = productType
	GlobalLogOption.ProductSubType = productSubType
}

//UseLogrusWithConfig 注册使用logrus日志
func UseLogrusWithConfig(fpath string) error {
	logConfigPath = fpath
	cf := conf.LogConfig{Log: &GlobalLogOption}
	c, err := conf.ParseConfig(fpath, &cf)
	if err != nil {
		return err
	}

	err = Config(c.Log)
	if err != nil {
		return err
	}

	return nil
}

//UseLogrus 注册使用logrus日志
func UseLogrus(optFunc LogConfigFunc) error {

	//执行传入的自定义设置全局配置函数(在此之前已执行包的init方法，初始化全局配置)
	if optFunc != nil {
		optFunc(&GlobalLogOption)
	}

	err := Config(&GlobalLogOption)
	if err != nil {
		return err
	}

	return nil
}

// Config 设定日志输出格式
func Config(opt *conf.LogOption) error {

	// 设定日志输出
	//var file *os.File
	logrus.SetLevel(logrus.Level(opt.Level))

	switch opt.Format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:  "2006-01-02 15:04:05.000",
			DisableTimestamp: !opt.DisableCustomTimestamp,
			FieldMap: logrus.FieldMap{
				//logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "type",
				logrus.FieldKeyMsg:   "message",
			},
		})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat:  "2006-01-02 15:04:05.000",
			DisableTimestamp: !opt.DisableCustomTimestamp,
			QuoteEmptyFields: true,
			FieldMap: logrus.FieldMap{
				//logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "type",
				logrus.FieldKeyMsg:   "message"},
		})
	}

	if opt.Output != "" {
		switch opt.Output {

		case "stderr":
			logrus.SetOutput(os.Stderr)

		case "file", "multi":
			if name := opt.OutputFile; name != "" {
				//_ = os.MkdirAll(filepath.Dir(name), 0777)

				// f, err := os.OpenFile(name+".log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
				// if err != nil {
				// 	return err
				// }

				// if opt.Output == "multi" {
				// 	out := io.MultiWriter(os.Stdout, f)
				// 	logrus.SetOutput(out)
				// } else {
				// 	logrus.SetOutput(f)
				// }

				setRotatelogs(opt, "")

			}

		default:
			logrus.SetOutput(os.Stdout)
		}
	}

	if !GlobalLogOption.DisableLineHook {
		//add LineHook
		logrus.AddHook(LineHook{
			Field: "source",
			Skip:  8,
		})
	}

	return nil
}

func setRotatelogs(opt *conf.LogOption, name string) {

	fileName := opt.OutputFile + opt.LogFilePathFormat
	if len(name) > 0 {
		fileName = opt.OutputFile + "-" + name + opt.LogFilePathFormat
	}

	logWriter, err := rotatelogs.New(
		// "log"+".%Y%m%d%H%M",
		// rotatelogs.WithLinkName("log"),
		fileName,
		rotatelogs.WithLinkName("log"), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(opt.LogFileMaxAge)*24*time.Hour),
		//rotatelogs.WithMaxAge(time.Duration(opt.LogFileMaxAge)),             // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Duration(opt.LogFileRotationTime)), // 日志切割时间间隔
	)

	if err != nil {
		panic(err)
	}

	if opt.Output == "multi" {
		out := io.MultiWriter(os.Stdout, logWriter)
		logrus.SetOutput(out)
	} else {
		logrus.SetOutput(logWriter)
	}

}

// CreateLogger 创建logger
func CreateLogger() *Entry {
	return CreateLoggerWithContext(nil)
}

// CreateLoggerWithContext 创建logger
func CreateLoggerWithContext(ctx context.Context) *Entry {
	if ctx == nil {
		ctx = context.Background()
	}

	fields := map[string]interface{}{
		AppNameKey: GlobalLogOption.AppName,
		AppNoKey:   GlobalLogOption.AppNo,
	}

	if len(GlobalLogOption.FrmDeviceID) > 0 {
		fields[FrmDeviceIDKey] = GlobalLogOption.FrmDeviceID
	}

	if len(GlobalLogOption.ProductType) > 0 {
		fields[ProductTypeKey] = GlobalLogOption.ProductType
	}

	if len(GlobalLogOption.ProductSubType) > 0 {
		fields[ProductSubTypeKey] = GlobalLogOption.ProductSubType
	}

	if GlobalLogOption.TIDFunc != nil {
		fields[TraceIDKey] = FromTraceIDContext(ctx)
	}

	return newEntry(logrus.WithFields(fields))
}

func newEntry(entry *logrus.Entry) *Entry {
	// if GlobalLogOption.DisableCustomTimestamp {
	// 	return &Entry{entry: entry}
	// }

	// return &Entry{entry: entry.WithField(
	// 	TimestampKey, time.Now().UTC().Unix())}

	return &Entry{entry: entry}
}

// FromTraceIDContext 从上下文中获取跟踪ID
func FromTraceIDContext(ctx context.Context) string {
	v := ctx.Value(traceIDContextKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return GlobalLogOption.TIDFunc()
}

//WithCustomTimestamp 追加自定义时间戳列
func (e *Entry) WithCustomTimestamp() *Entry {
	if GlobalLogOption.DisableCustomTimestamp {
		return e
	}

	return e.WithFields(map[string]interface{}{
		TimestampKey: time.Now().UTC().Unix()})
}

func (e *Entry) innerEntry() *logrus.Entry {
	return e.WithCustomTimestamp().entry
}

func (e *Entry) checkAndDelete(fields map[string]interface{}, keys ...string) {
	for _, key := range keys {
		if _, ok := fields[key]; ok {
			delete(fields, key)
		}
	}
}

// WithFields 结构化字段写入
func (e *Entry) WithFields(fields map[string]interface{}) *Entry {
	e.checkAndDelete(fields,
		TraceIDKey,
	)
	return newEntry(e.entry.WithFields(fields))
}

// WithField 结构化字段写入
func (e *Entry) WithField(key string, value interface{}) *Entry {
	return e.WithFields(map[string]interface{}{key: value})
}

// Fatalf 重大错误日志
func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.innerEntry().Fatalf(format, args...)
}

// Errorf 错误日志
func (e *Entry) Errorf(format string, args ...interface{}) {
	e.innerEntry().Errorf(format, args...)
}

// Warnf 警告日志
func (e *Entry) Warnf(format string, args ...interface{}) {
	e.innerEntry().Warnf(format, args...)
}

// Infof 消息日志
func (e *Entry) Infof(format string, args ...interface{}) {
	e.innerEntry().Infof(format, args...)
}

// Printf 消息日志
func (e *Entry) Printf(format string, args ...interface{}) {
	e.innerEntry().Printf(format, args...)
}

// Debugf 写入调试日志
func (e *Entry) Debugf(format string, args ...interface{}) {
	e.innerEntry().Debugf(format, args...)
}

// Trace Trace日志
func (e *Entry) Trace(args ...interface{}) {
	e.innerEntry().Log(logrus.TraceLevel, args...)
}

//Debug Debug日志
func (e *Entry) Debug(args ...interface{}) {
	e.innerEntry().Log(logrus.DebugLevel, args...)
}

//Print 消息日志
func (e *Entry) Print(args ...interface{}) {
	e.innerEntry().Info(args...)
}

//Info 消息日志
func (e *Entry) Info(args ...interface{}) {
	e.innerEntry().Log(logrus.InfoLevel, args...)
}

//Warn 警告日志
func (e *Entry) Warn(args ...interface{}) {
	e.innerEntry().Log(logrus.WarnLevel, args...)
}

//Warning 警告日志
func (e *Entry) Warning(args ...interface{}) {
	e.innerEntry().Warn(args...)
}

//Error Error日志
func (e *Entry) Error(args ...interface{}) {
	e.innerEntry().Log(logrus.ErrorLevel, args...)
}

//Fatal Fatal日志
func (e *Entry) Fatal(args ...interface{}) {
	e.innerEntry().Log(logrus.FatalLevel, args...)
	e.entry.Logger.Exit(1)
}

//Panic Panic日志
func (e *Entry) Panic(args ...interface{}) {
	e.innerEntry().Log(logrus.PanicLevel, args...)
	panic(fmt.Sprint(args...))
}
