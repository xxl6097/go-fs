package z

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	Level            string `yaml:"level"`            // 日志级别
	TagName          string `yaml:"tagName"`            // 日志级别
	LogDir           string `yaml:"logDir"`           // 日志目录
	LogName          string `yaml:"logName"`           // 日志目录
	Path             string `yaml:"path"`             // 主日志路径
	ErrorPath        string `yaml:"errorPath"`        // 错误日志路径（可选）
	MaxSize          int    `yaml:"maxSize"`          // 文件最大大小(MB)
	MaxBackups       int    `yaml:"maxBackups"`       // 最大备份数
	MaxAge           int    `yaml:"maxAge"`           // 最大保存天数
	Compress         bool   `yaml:"compress"`         // 是否压缩
	SeparateErrorLog bool   `yaml:"separateErrorLog"` // 是否分离错误日志
	Enable           bool   `yaml:"enable"`           // 日志开关，默认开启
	AddCallerSkip    int    `yaml:"addCallerSkip"`
	Hook             func(zapcore.Entry) error
	Logger           *zap.Logger
}

//// DefaultConfig 获取日志配置
//func DefaultConfig(fn func(*LogConfig)) *LogConfig {
//	cfg := &LogConfig{
//		Level: "info",
//	}
//	if fn != nil {
//		fn(cfg)
//	}
//	return cfg
//}
//
//func FileConfig(fn func(*LogConfig)) *LogConfig {
//	return DefaultConfig(func(conf *LogConfig) {
//		conf.Level = "info"
//		conf.Path = "./logs/app.log"
//		conf.ErrorPath = "./logs/error.log"
//		conf.MaxSize = 100
//		conf.MaxBackups = 30
//		conf.MaxAge = 7
//		conf.Compress = true
//		conf.SeparateErrorLog = true
//		if fn != nil {
//			fn(conf)
//		}
//	})
//}

// 私有全局实例
var instance *LogConfig
var once sync.Once

// GetLogConfig 获取单例
func GetLogConfig() *LogConfig {
	once.Do(func() {
		// 只执行一次
		instance = &LogConfig{
			Level:            "debug",
			Path:             "./logs/app.log",
			ErrorPath:        "./logs/error.log",
			MaxSize:          100,
			MaxBackups:       30,
			MaxAge:           7,
			Compress:         true,
			SeparateErrorLog: true,
			Enable:           true,
			Hook:             nil,
			AddCallerSkip:    0,
			Logger:           nil,
		}
	})
	return instance
}
