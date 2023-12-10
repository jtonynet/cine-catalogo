package decorators

import (
	"github.com/jtonynet/cine-catalogo/internal/interfaces"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

type LoggerDecorated struct {
	next interfaces.Logger
}

var (
	logEventsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "log_events_total",
			Help: "Total number of log events.",
		},
		[]string{"level"},
	)
)

func NewLoggerWithMetrics(next interfaces.Logger) interfaces.Logger {
	return &LoggerDecorated{
		next: next,
	}
}

func (l *LoggerDecorated) Debug(msg string, fields ...zap.Field) {
	l.next.Debug(msg, fields...)

	logEventsTotal.WithLabelValues("debug").Inc()
}

func (l *LoggerDecorated) Info(msg string, fields ...zap.Field) {
	l.next.Info(msg, fields...)

	logEventsTotal.WithLabelValues("info").Inc()
}

func (l *LoggerDecorated) Warning(msg string, fields ...zap.Field) {
	l.next.Warning(msg, fields...)

	logEventsTotal.WithLabelValues("warning").Inc()
}

func (l *LoggerDecorated) Error(msg string, fields ...zap.Field) {
	l.next.Error(msg, fields...)

	logEventsTotal.WithLabelValues("error").Inc()
}

// Sync flushes any buffered log entries.
func (l *LoggerDecorated) Sync() error {
	return l.next.Sync()
}

// WithField returns a new Logger with an additional field.
func (l *LoggerDecorated) WithField(key string, value interface{}) interfaces.Logger {
	return l.next.WithField(key, value)
}

// WithFields returns a new Logger with additional fields.
func (l *LoggerDecorated) WithFields(fields ...zap.Field) interfaces.Logger {
	return l.next.WithFields(fields...)
}

// WithError returns a new Logger with an additional error field.
func (l *LoggerDecorated) WithError(err error) interfaces.Logger {
	return l.next.WithError(err)
}

// To implement io.Writer interface for the Writer field.
func (l *LoggerDecorated) Write(p []byte) (n int, err error) {
	return l.next.Write(p)
}
