package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// DefaultOption returns the default zap logger options
func DefaultOption() []zap.Option {
	var stackTraceLevel zap.LevelEnablerFunc = func(level zapcore.Level) bool {
		return level >= zapcore.DPanicLevel
	}
	return []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(stackTraceLevel),
	}
}

// DefaultEncoder returns the default encoder, whish is utf-8
func DefaultEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(DefaultEncoderConfig())
}

// DefaultEncoderConfig returns the default encoder config
func DefaultEncoderConfig() zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	return config
}

// DefaultLumberjackLogger returns the default lumberjack logger, which is used to split logs
func DefaultLumberjackLogger() *lumberjack.Logger {
	return &lumberjack.Logger{
		MaxSize:   200, // MB
		LocalTime: true,
		Compress:  true,
	}
}
