package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jtonynet/cine-catalogo/internal/interfaces"
)

// Logger encapsulates the zap logger.
type Logger struct {
	zapLogger *zap.Logger
}

type LogField struct {
	key   string
	value interface{}
}

func NewLogger() (*Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.OutputPaths = []string{"stdout"}
	config.DisableCaller = true

	zapLogger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("error creating logger: %v", err)
	}

	return &Logger{
		zapLogger: zapLogger,
	}, nil
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, fields...)
}

// Info logs an info message.
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, fields...)
}

// Warning logs a warning message.
func (l *Logger) Warning(msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, fields...)
}

// Error logs an error message.
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zapLogger.Error(msg, fields...)
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() error {
	return l.zapLogger.Sync()
}

// WithField returns a new Logger with an additional field.
func (l *Logger) WithField(key string, value interface{}) interfaces.Logger {
	return &Logger{
		zapLogger: l.zapLogger.With(zap.Any(key, value)),
	}
}

// WithFields returns a new Logger with additional fields.
func (l *Logger) WithFields(fields ...zap.Field) interfaces.Logger {
	return &Logger{
		zapLogger: l.zapLogger.With(fields...),
	}
}

// WithError returns a new Logger with an additional error field.
func (l *Logger) WithError(err error) interfaces.Logger {
	return &Logger{
		zapLogger: l.zapLogger.With(zap.Error(err)),
	}
}

// To implement io.Writer interface for the Writer field.
func (l *Logger) Write(p []byte) (n int, err error) {
	l.zapLogger.Info(string(p))
	return len(p), nil
}

func (zf *LogField) apply(logger Logger) {
	logger.WithField(zf.key, zf.value)
}

func NewZapLogField(key string, value interface{}) *LogField {
	return &LogField{
		key:   key,
		value: value,
	}
}
