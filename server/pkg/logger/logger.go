package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	sugar *zap.SugaredLogger
	once sync.Once
)

// Init 初始化日志系统（应用启动时调用一次）
func Init(level string, logPath string) {
	once.Do(func() {
		var zapLevel zapcore.Level
		switch level {
		case "debug":
			zapLevel = zapcore.DebugLevel
		case "warn":
			zapLevel = zapcore.WarnLevel
		case "error":
			zapLevel = zapcore.ErrorLevel
		default:
			zapLevel = zapcore.InfoLevel
		}

		// 编码器配置（JSON格式，便于ELK收集）
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// 多输出目标：控制台 + 文件
		core := zapcore.NewTee(
			// 控制台输出（带颜色，开发调试用）
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoderConfig),
				os.Stdout,
				zapLevel,
			),
		)

		// 如果配置了文件路径，添加文件输出
		if logPath != "" {
			file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err == nil {
				fileCore := zapcore.NewCore(
					zapcore.NewJSONEncoder(encoderConfig), // 文件使用JSON格式
					file,
					zapLevel,
				)
				core = zapcore.NewTee(core, fileCore)
			}
		}

		// 创建Logger实例（添加调用者信息和堆栈跟踪）
		log = zap.New(core,
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
		sugar = log.Sugar()
	})
}

// Debug 调试级别日志
func Debug(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

// Info 信息级别日志
func Info(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

// Warn 警告级别日志
func Warn(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

// Error 错误级别日志
func Error(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

// Fatal 致命错误日志（会终止程序）
func Fatal(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}

// With 为日志添加上下文字段（用于链路追踪）
func With(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}

// Sugar 返回SugaredLogger（支持printf风格格式化）
func Sugar() *zap.SugaredLogger {
	return sugar
}

// Logger 返回原始Zap.Logger
func Logger() *zap.Logger {
	return log
}

// Sync 刷新缓冲区（程序退出前应调用）
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
