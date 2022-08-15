package logger

import (
	"context"
	"fmt"
	"os"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ctxMarker, ctxLogger, and ctxMarkerKey are for extracting logger from ctx
type ctxMarker struct{}

type ctxLogger struct {
	logger *zap.Logger
	fields []zapcore.Field
}

var ctxMarkerKey = &ctxMarker{}

var cLogger *ctxLogger = &ctxLogger{}

// defaultLogFields defines the tag keys whose values should be logged
var defaultLogFields = []string{
	"trace.traceid",
	"trace.spanid",
	"grpc.request.service",
	"grpc.request.method",
}

// InitLogger - Initializes the default logger.
//
// level indicates the lowest level of logs.
func InitLogger(level zapcore.Level) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.ConsoleSeparator = CONSOLE_SEPARATOR
	fileEncoder := zapcore.NewConsoleEncoder(config)

	logFile, _ := os.OpenFile(LOG_FILENAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, LOG_PERMISSION)
	writer := zapcore.AddSync(logFile)

	defaultLogLevel := level

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)
	cLogger.logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

// WithContext - returns the context logger, or the default logger if context
// logger does not exist.
//
// Default log fields are determined by the entries of defaultLogFields.
func WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return cLogger.logger
	}

	currLogger := cLogger
	zapFields := []zapcore.Field{}

	l, ok := ctx.Value(ctxMarkerKey).(*ctxLogger)
	if ok && l != nil {
		currLogger = l
	}

	tags := grpc_ctxtags.Extract(ctx)

	for _, tagKey := range defaultLogFields {
		if tags.Has(tagKey) {
			tagValue := fmt.Sprint(tags.Values()[tagKey])

			zapFields = append(zapFields, zap.String(tagKey, tagValue))
		}
	}

	zapFields = append(zapFields, cLogger.fields...)

	return currLogger.logger.With(zapFields...)
}

// Sync - flushes any buffered log entries. Should be called before application
// exits.
func Sync() {
	_ = cLogger.logger.Sync()
}

// AddPermanentFields - fields will be permanently logged by the logger before all fields
// specified in defaultLogFields.
func AddPermanentFields(ctx context.Context, fields ...zapcore.Field) {
	l, ok := ctx.Value(ctxMarkerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}

	l.fields = append(l.fields, fields...)
}
