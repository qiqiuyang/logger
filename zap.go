package logger

import (
	"github.com/natefinch/lumberjack"
	"os"
	"time"

	"github.com/qiqiuyang/logger/model"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var suffix string

// 自定义初始化zap日志
func Zap(config model.Zap) (logger *zap.Logger) {
	logLevel := paserLoggerLevel(config.Level)
	suffix = config.Suffix
	// logPath := getLogPath(config.FilePath)
	// 限定打印级别 例如设置为info，则只会保存info到Fatal之间的级别日志
	priority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return logLevel <= lev && lev <= zapcore.FatalLevel
	})

	cores := [...]zapcore.Core{
		getEncoderCore(config, priority),
	}
	logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())

	if config.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	zap.ReplaceGlobals(logger)
	return logger
}

// 根据输入日志等级的字符返回对应zap的日志等级定义
func paserLoggerLevel(level string) zapcore.Level {
	logLevel := zap.DebugLevel
	switch level {
	case "debug":
		logLevel = zap.DebugLevel // -1
	case "info":
		logLevel = zap.InfoLevel // 0
	case "warn":
		logLevel = zap.WarnLevel // 1
	case "error":
		logLevel = zap.ErrorLevel // 2
	case "panic":
		logLevel = zap.PanicLevel // 4
	case "fatal":
		logLevel = zap.FatalLevel // 5
	default:
		logLevel = zap.InfoLevel // 0
	}
	return logLevel
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(config model.Zap, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer := getWriteSyncer(config) // 使用file-rotatelogs进行日志分割
	return zapcore.NewCore(getEncoder(config), writer, level)
}

// getEncoder 获取zapcore.Encoder
func getEncoder(config model.Zap) zapcore.Encoder {
	if config.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig(config))
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig(config))
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig(config model.Zap) (config_ zapcore.EncoderConfig) {
	config_ = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      config.CallerKey,
		StacktraceKey:  config.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case config.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config_.EncodeLevel = zapcore.LowercaseLevelEncoder
	case config.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config_.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case config.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config_.EncodeLevel = zapcore.CapitalLevelEncoder
	case config.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config_.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config_.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config_
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000" + suffix))
}

// @function: GetWriteSyncer
// @description: zap logger中加入file-rotatelogs
// @return: zapcore.WriteSyncer, error
func getWriteSyncer(config model.Zap) zapcore.WriteSyncer {
	logFileMode := os.FileMode(0666)

	_, err := os.OpenFile(config.FilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, logFileMode)
	if err != nil {
		if os.IsNotExist(err) {
			_, _ = os.Create(config.FilePath)
		}
	}

	_ = os.Chmod(config.FilePath, logFileMode)

	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.FilePath,   // 日志文件的位置
		MaxSize:    config.MaxSize,    // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: config.MaxBackups, // 保留旧文件的最大个数
		MaxAge:     config.MaxAge,     // 保留旧文件的最大天数
		Compress:   config.Compress,   // 是否压缩/归档旧文件
	}

	if config.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}
