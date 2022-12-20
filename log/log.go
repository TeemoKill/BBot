package log

import (
	"path"
	"runtime"
	"time"

	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// Init 使用默认的日志格式配置，会写入到logs文件夹内，日志会保留七天
func Init() (l *logrus.Logger, err error) {
	writer, err := rotatelogs.New(
		path.Join("logs", "%Y-%m-%d.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		logrus.WithError(err).Error("unable to write logs")
		return nil, err
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		PadLevelText:     true,
		QuoteEmptyFields: true,
	})
	logrus.AddHook(lfshook.NewHook(writer, &logrus.TextFormatter{
		FullTimestamp:    true,
		PadLevelText:     true,
		QuoteEmptyFields: true,
		ForceQuote:       true,
	}))
	return logrus.StandardLogger(), err
}

func StandardLogger() *logrus.Logger {
	return logrus.StandardLogger()
}

// ModuleLogger - 提供一个为 Module 使用的 logrus.Entry
// 包含 logrus.Fields
func ModuleLogger(name string) *logrus.Entry {
	return logrus.WithField("module", name)
}

// CurrentModuleLogger provides a logrus.Entry with the current module function name
func CurrentModuleLogger() *logrus.Entry {
	moduleFrame := currentFrame(1)
	return ModuleLogger(moduleFrame.Function)
}

// currentFrame retrieve the current function call frame
func currentFrame(extraSkip int) *runtime.Frame {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2+extraSkip, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return &frame
}
