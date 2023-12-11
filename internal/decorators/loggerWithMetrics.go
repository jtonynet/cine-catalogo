package decorators

import (
	"github.com/jtonynet/cine-catalogo/internal/interfaces"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

func (l *LoggerDecorated) Debug(msg string, fields ...interfaces.LogField) {
	logEventsTotal.WithLabelValues("debug").Inc()
	l.next.Debug(msg, fields...)
}

func (l *LoggerDecorated) Info(msg string, fields ...interfaces.LogField) {
	logEventsTotal.WithLabelValues("info").Inc()
	l.next.Info(msg, fields...)
}

func (l *LoggerDecorated) Warning(msg string, fields ...interfaces.LogField) {
	logEventsTotal.WithLabelValues("warning").Inc()
	l.next.Warning(msg, fields...)
}

func (l *LoggerDecorated) Error(msg string, fields ...interfaces.LogField) {
	logEventsTotal.WithLabelValues("error").Inc()
	l.next.Error(msg, fields...)
}

func (l *LoggerDecorated) Sync() error {
	return l.next.Sync()
}

func (l *LoggerDecorated) WithField(key string, value interface{}) interfaces.Logger {
	decoratedLogger := &LoggerDecorated{next: l.next.WithField(key, value)}
	return decoratedLogger
}

func (l *LoggerDecorated) WithFields(fields ...interfaces.LogField) interfaces.Logger {
	decoratedLogger := &LoggerDecorated{next: l.next.WithFields(fields...)}
	return decoratedLogger
}

func (l *LoggerDecorated) WithError(err error) interfaces.Logger {
	decoratedLogger := &LoggerDecorated{next: l.next.WithError(err)}
	return decoratedLogger
}

func (l *LoggerDecorated) Write(p []byte) (n int, err error) {
	return l.next.Write(p)
}
