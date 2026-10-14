package z

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

//var Hook func(zapcore.Entry) error

func LoadLogger(fn func(conf *LogConfig)) {
	cfg := GetLogConfig()
	if fn != nil {
		fn(cfg)
	}
	// 应用日志开关配置
	logEnabled.Store(cfg.Enable)
	if cfg.LogDir != "" {
		if cfg.LogName != "" {
			cfg.Path = filepath.Join(cfg.LogDir, cfg.LogName)
		} else {
			cfg.Path = filepath.Join(cfg.LogDir, "app.log")
		}

		cfg.ErrorPath = filepath.Join(cfg.LogDir, "error.log")
	}
	logger := initZapLogger(cfg, cfg.AddCallerSkip)
	if cfg.TagName != "" {
		zap.ReplaceGlobals(logger.Named(cfg.TagName))
	} else {
		zap.ReplaceGlobals(logger.Named("glog"))
	}

	//fmt.Printf("加载日志 fn:%+v cfg:%+v\n", fn, cfg)
	check(cfg)
}

// initZapLogger 初始化日志
func initZapLogger(cfg *LogConfig, addCallerSkip int) *zap.Logger {
	// 设置日志级别
	var level = getZapLevel(cfg.Level)
	var consoleEncoder, fileEncoder zapcore.Encoder
	env := strings.ToLower(os.Getenv("APP_ENV"))
	switch env {
	case "prod", "production":
		consoleEncoder = zapcore.NewJSONEncoder(jsonEncoderConfig())
		fileEncoder = consoleEncoder
	default:
		consoleEncoder = zapcore.NewConsoleEncoder(consoleEncoderConfig(true))
		fileEncoder = zapcore.NewConsoleEncoder(consoleEncoderConfig(false))
	}
	// 创建核心列表
	var cores []zapcore.Core
	cores = append(cores, createMainCore(consoleEncoder, level))
	if cfg.Path != "" {
		fileCore := createNormalCore(cfg, fileEncoder, level)
		cores = append(cores, fileCore)
	}
	if cfg.SeparateErrorLog && cfg.ErrorPath != "" {
		errorCore := createErrorCore(cfg, fileEncoder, zap.ErrorLevel)
		cores = append(cores, errorCore)
	}
	//zap.AddCallerSkip(1)
	// 创建 tee 核心（多路复用）
	core := zapcore.NewTee(cores...)
	// 包装日志开关
	core = &switchCore{Core: core, enabled: &logEnabled}
	// 创建 Logger
	//ZapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0), zap.AddStacktrace(zap.ErrorLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(addCallerSkip), zap.AddStacktrace(zap.ErrorLevel),
		zap.Hooks(func(entry zapcore.Entry) error {
			if cfg.Hook != nil {
				return cfg.Hook(entry)
			}
			return nil
		}))
	//SugaredLogger = ZapLogger.Sugar()
	//logger.Info("logger init success")
	//zap.ReplaceGlobals(logger.Named("glog"))
	return logger
}

func createMainCore(encoder zapcore.Encoder, level zapcore.Level) zapcore.Core {
	consoleWriter := zapcore.Lock(zapcore.AddSync(os.Stdout))
	return zapcore.NewCore(encoder, consoleWriter, level)
}

func createNormalCore(cfg *LogConfig, encoder zapcore.Encoder, level zapcore.Level) zapcore.Core {
	fileWriter := zapcore.Lock(zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Path,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}))
	return zapcore.NewCore(encoder, fileWriter, level)
}

func createErrorCore(cfg *LogConfig, encoder zapcore.Encoder, level zapcore.Level) zapcore.Core {
	fileWriter := zapcore.Lock(zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.ErrorPath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}))
	return zapcore.NewCore(encoder, fileWriter, level)
}

// getZapLevel 获取 Zap 日志级别
func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// customTimeEncoder 自定义时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	loc, err := time.LoadLocation("Asia/Shanghai") // 等价于 UTC+8
	if err != nil {
		loc = time.FixedZone("CST", 8*3600) // 东八区
	}
	beijingTime := t.In(loc)
	enc.AppendString(beijingTime.Format("2006-01-02 15:04:05.000"))
}

// consoleEncoderConfig 控制台编码器配置
func consoleEncoderConfig(color bool) zapcore.EncoderConfig {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:   "time",
		LevelKey:  "level",
		NameKey:   "logger",
		CallerKey: "caller",
		//FunctionKey: "function", // 添加函数名
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径
	}

	if color {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	return encoderConfig
}

// jsonEncoderConfig JSON编码器配置
func jsonEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写level
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601格式
		EncodeDuration: zapcore.MillisDurationEncoder, // 毫秒
		EncodeCaller:   zapcore.ShortCallerEncoder,    // 短路径
	}
}

// isTerminal 判断是否为终端
func isTerminal1() bool {
	fileInfo, _ := os.Stdout.Stat()
	result := fileInfo.Mode() & os.ModeCharDevice
	return result != 0
}

func L() *zap.Logger {
	//if GetLogConfig().Logger == nil {
	//	fmt.Println("Logger is nill")
	//}
	l := zap.L()
	if l.Name() == "" {
		LoadLogger(nil)
	}
	return zap.L()
}
