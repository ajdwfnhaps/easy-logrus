package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/ajdwfnhaps/easy-logrus/common"
)

//LogConfig 统一配置类
type LogConfig struct {
	Log *LogOption `toml:"log"`
}

type (
	// TraceIDFunc 定义获取跟踪ID的函数
	TraceIDFunc func() string
	// LogLevel type
	LogLevel uint32
)

//LogOption 日志配置
type LogOption struct {
	AppNo int `toml:"app_no"`
	//应用名称
	AppName string `toml:"app_name"`
	//日志级别(1:fatal 2:error,3:warn,4:info,5:debug,6:trace)
	Level LogLevel `toml:"log_level"`
	//日志格式（支持输出格式：text/json）
	Format string `toml:"format"`
	//日志输出(支持：stdout/stderr/file)
	Output string `toml:"output"`
	//指定日志输出的文件路径 logs/app.log
	OutputFile string `toml:"output_file"`
	//获取跟踪ID的函数
	TIDFunc TraceIDFunc
	//是否禁用自定义时间戳显示
	DisableCustomTimestamp bool `toml:"disable_custom_timestamp"`
	//是否禁用行号信息显示(WarnLevel以上才会显示)
	DisableLineHook bool `toml:"disable_line_hook"`
	//设置日志文件清理前的最长保存时间 天数
	LogFileMaxAge int `toml:"log_file_max_age"`
	//设置日志分割的时间
	LogFileRotationTime int `toml:"log_file_rotation_time"`
	//设置日志文件名规则
	LogFilePathFormat string `toml:"log_file_path_format"`

	//来源设备ID
	FrmDeviceID string
	//来源设备所属产品类型
	ProductType string
	//来源设备所属产品子类
	ProductSubType string
}

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel LogLevel = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

//ParseConfig 解析配置文件
func ParseConfig(fpath string, c *LogConfig) (*LogConfig, error) {
	if confFileExists, err := common.PathExists(fpath); confFileExists != true {
		fmt.Println("Config File ", fpath, " is not exist!!")
		return nil, err
	}

	_, err := toml.DecodeFile(fpath, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
