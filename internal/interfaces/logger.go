package interfaces

import "go.uber.org/zap"

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warning(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)

	Sync() error
	WithField(key string, value interface{}) Logger
	WithFields(fields ...zap.Field) Logger
	WithError(err error) Logger
	Write(p []byte) (n int, err error)
}

type LogField interface {
	apply(field interface{})
}
