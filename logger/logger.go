// Package logger
// @description wraps the zap logger and lumberjack
package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Plugin = zapcore.Core

func NewLogger(plugin zapcore.Core, options ...zap.Option) *zap.Logger {
	return zap.New(plugin, append(DefaultOption(), options...)...)
}

func NewPlugin(writer zapcore.WriteSyncer, enabler zapcore.LevelEnabler) Plugin {
	return zapcore.NewCore(DefaultEncoder(), writer, enabler)
}

func NewStdoutPlugin(enabler zapcore.LevelEnabler) Plugin {
	return NewPlugin(zapcore.Lock(zapcore.AddSync(os.Stdout)), enabler)
}

func NewStderrPlugin(enabler zapcore.LevelEnabler) Plugin {
	return NewPlugin(zapcore.Lock(zapcore.AddSync(os.Stderr)), enabler)
}

// NewFilePlugin creates a file plugin
// After get an io.Closer instance, You need to invoke Close() before
// the process exits to ensure that all the writing can be flushed to disk.
func NewFilePlugin(filePath string, enabler zapcore.LevelEnabler) (Plugin, io.Closer) {
	logger := DefaultLumberjackLogger()
	logger.Filename = filePath
	return NewPlugin(zapcore.AddSync(logger), enabler), logger
}
