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

func TestLogConfigFile(t *testing.T) {

	//-----------初始化日志组件
	if err := UseLogrusWithConfig("D:\\Work\\ajdwfnhaps\\Product\\easy-logrus\\logger\\log.toml"); err != nil {
		t.Error(err)
	}

	DeviceInit("9100156", "134", "7")

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
