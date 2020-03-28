# easy-logrus
封装logrus组件,简单易用,方便在项目中集成应用
1. 可使用配置文件
2. 支持日志分割（file-rotatelogs）
3. 集成LineHook

### 使用easy-logrus

```
package logger

import (
	"testing"
	"time"

	"github.com/ajdwfnhaps/easy-logrus/conf"
)

func TestLog(t *testing.T) {
	if err := UseLogrus(func(c *conf.LogOption) {

		c.AppName = "Go应用002"
		c.LogFileRotationTime = 60
		c.LogFilePathFormat = ".%Y-%m-%d-%H-%M.log"
		//c.Format = "text"
		//c.Level = WarnLevel

	}); err != nil {
		t.Error(err)
	}

	logger := CreateLogger()

	//增加自定义字段
	logger.WithField("animal", "walrus").Info("A walrus appears")

	logger.Debugf("just for \r\n test %s", " oh my god")
	logger.Infof("info...info...info...F %s", " oh my god")

	logger.Warn("Warning \r\n test")
	logger.Error("abc \r\n test")

}

```


### 使用easy-logrus通过配置文件

```

package logger

import (
	"testing"
	"time"

	"github.com/ajdwfnhaps/easy-logrus/conf"
)

func TestLogConfigFile(t *testing.T) {

	//-----------初始化日志组件
	if err := UseLogrusWithConfig("D:\\Work\\ajdwfnhaps\\Product\\easy-logrus\\logger\\log.toml"); err != nil {
		t.Error(err)
	}

	//可去掉追踪tid
	GlobalLogOption.TIDFunc = nil

	//-----------日志组件使用

	logger := CreateLogger()

	//增加自定义字段
	logger.WithField("animal", "walrus").Info("A walrus appears")

	time.Sleep(2 * time.Second)

	logger.Debugf("just for \r\n test %s", " oh my god")
	logger.Infof("info...info...info...F %s", " oh my god")

	time.Sleep(1 * time.Second)

	logger.Warn("Warning \r\n test")
	logger.Error("abc \r\n test")

	CreateLogger().Info("shit...")

}


```

### 配置文件详解

```

# 日志配置
[log]
# 应用编号，例如：10001
app_no = 10001
# 应用名称
app_name = "Go应用003"
# 日志级别(1:fatal 2:error,3:warn,4:info,5:debug)
log_level = 5
# 日志格式（支持输出格式：text/json）
format = "json"
# 日志输出(支持：stdout/stderr/file/multi)
# output = "multi"
# 指定日志输出的文件路径
output_file = "logs/app"
# 是否禁用自定义时间戳显示
disable_custom_timestamp = false
# 是否禁用行号信息显示(WarnLevel以上才会显示)
disable_line_hook = false
# 设置保留天数
log_file_max_age = 7
# 设置每天分割日志文件
log_file_rotation_time = 86400
# 设置日志文件名规则
log_file_path_format = ".%Y-%m-%d.log"

```

详细可参考[测试用例](logger/logger_test.go)
