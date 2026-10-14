package z

import (
	"fmt"

	"go.uber.org/zap"
)

var oldLog *zap.Logger

func check(cfg *LogConfig) {
	if oldLog == nil && cfg != nil {
		oldLog = initZapLogger(cfg, 1)
	}
}

func Debugf(format string, v ...interface{}) {
	oldLog.Debug(fmt.Sprintf(format, v...))
}

func Debug(v ...interface{}) {
	oldLog.Debug(fmt.Sprintln(v...))
}

func Infof(format string, v ...interface{}) {
	oldLog.Info(fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) {
	oldLog.Info(fmt.Sprintln(v...))
}

func Warnf(format string, v ...interface{}) {
	oldLog.Warn(fmt.Sprintf(format, v...))
}

func Warn(v ...interface{}) {
	oldLog.Warn(fmt.Sprintln(v...))
}

func Errorf(format string, v ...interface{}) {
	oldLog.Error(fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	oldLog.Error(fmt.Sprintln(v...))
}

func Fatalf(format string, v ...interface{}) {
	oldLog.Fatal(fmt.Sprintf(format, v...))
}

func Fatal(v ...interface{}) {
	oldLog.Fatal(fmt.Sprintln(v...))
}

func Printf(format string, v ...interface{}) {
	oldLog.Debug(fmt.Sprintf(format, v...))
}

func Println(v ...interface{}) {
	oldLog.Debug(fmt.Sprintln(v...))
}
