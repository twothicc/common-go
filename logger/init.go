package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// InitLogger - Initializes the default logger.
// isTest indicates whether the environment is test or production.
// logger is configured to display Info level and above for production and
// Debug level and above for test.
func InitLogger(isTest bool) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.ConsoleSeparator = CONSOLE_SEPARATOR
	fileEncoder := zapcore.NewConsoleEncoder(config)

	logFile, _ := os.OpenFile(LOG_FILENAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)

	defaultLogLevel := zapcore.InfoLevel
	if isTest {
		defaultLogLevel = zapcore.DebugLevel
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

// NewLogContext - adds the current zap logger into context.
// Values specified in fields will be propagated down as part of log context.
//
// e.g. Pass in a request id to trace the function calls of a specific request
// call.
func NewLogContext(ctx context.Context, fields ...zap.Field) context.Context {
	return context.WithValue(ctx, LOGGER_KEY, WithContext(ctx).With(fields...))
}

// WithContext - returns the context logger, or the default logger if context
// logger does not exist.
func WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return logger
	}

	if ctxLogger, ok := ctx.Value(LOGGER_KEY).(*zap.Logger); ok {
		return ctxLogger
	} else {
		return logger
	}
}
