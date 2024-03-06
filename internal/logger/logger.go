package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MustInit initializes zap logger and returns it
func MustInit(level string) *zap.Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	}

	logger := zap.Must(config.Build())
	return logger
}
