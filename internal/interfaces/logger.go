package interfaces

type Logger interface {
	Debug(msg string, fields ...LogField)
	Info(msg string, fields ...LogField)
	Warning(msg string, fields ...LogField)
	Error(msg string, fields ...LogField)

	Sync() error
	WithField(key string, value interface{}) Logger
	WithFields(fields ...LogField) Logger
	WithError(err error) Logger
	Write(p []byte) (n int, err error)
}

type LogField interface {
	GetKey() string
	GetValue() interface{}
}
