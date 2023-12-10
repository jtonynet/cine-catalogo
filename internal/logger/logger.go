package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/jtonynet/cine-catalogo/internal/interfaces"
)

type LogField struct {
	Key   string
	Value interface{}
}

func LogFieldFactory(key string, value interface{}) *LogField {
	return &LogField{Key: key, Value: value}
}

func (lf *LogField) GetKey() string {
	return lf.Key
}

func (lf *LogField) GetValue() interface{} {
	return lf.Value
}

// Logger encapsulates the zap logger.
type Logger struct {
	zapLogger *zap.Logger
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

func convertLogFieldsToZapFields(fields ...interfaces.LogField) []zap.Field {
	zapFields := make([]zap.Field, len(fields))

	for i, f := range fields {
		zapFields[i] = zap.Any(f.GetKey(), f.GetValue())
	}

	return zapFields
}

func (l *Logger) Info(msg string, fields ...interfaces.LogField) {
	zapFields := convertLogFieldsToZapFields(fields...)
	l.zapLogger.Info(msg, zapFields...)
}

func (l *Logger) Debug(msg string, fields ...interfaces.LogField) {
	zapFields := convertLogFieldsToZapFields(fields...)
	l.zapLogger.Debug(msg, zapFields...)
}

func (l *Logger) Warning(msg string, fields ...interfaces.LogField) {
	zapFields := convertLogFieldsToZapFields(fields...)
	l.zapLogger.Warn(msg, zapFields...)
}

func (l *Logger) Error(msg string, fields ...interfaces.LogField) {
	zapFields := convertLogFieldsToZapFields(fields...)
	l.zapLogger.Error(msg, zapFields...)
}

func (l *Logger) Sync() error {
	return l.zapLogger.Sync()
}

func (l *Logger) WithField(key string, value interface{}) interfaces.Logger {
	return &Logger{
		zapLogger: l.zapLogger.With(zap.Any(key, value)),
	}
}

func (l *Logger) WithFields(fields ...interfaces.LogField) interfaces.Logger {
	zapFields := convertLogFieldsToZapFields(fields...)
	return &Logger{
		zapLogger: l.zapLogger.With(zapFields...),
	}
}

func (l *Logger) WithError(err error) interfaces.Logger {
	return &Logger{
		zapLogger: l.zapLogger.With(zap.Error(err)),
	}
}

func (l *Logger) Write(p []byte) (n int, err error) {
	l.zapLogger.Info(string(p))
	return len(p), nil
}
