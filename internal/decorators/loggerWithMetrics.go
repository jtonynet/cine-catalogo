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
	l.next.Debug(msg, fields...)

	logEventsTotal.WithLabelValues("debug").Inc()
}

func (l *LoggerDecorated) Info(msg string, fields ...interfaces.LogField) {
	l.next.Info(msg, fields...)

	logEventsTotal.WithLabelValues("info").Inc()
}

func (l *LoggerDecorated) Warning(msg string, fields ...interfaces.LogField) {
	l.next.Warning(msg, fields...)

	logEventsTotal.WithLabelValues("warning").Inc()
}

func (l *LoggerDecorated) Error(msg string, fields ...interfaces.LogField) {
	l.next.Error(msg, fields...)

	logEventsTotal.WithLabelValues("error").Inc()
}

func (l *LoggerDecorated) Sync() error {
	return l.next.Sync()
}

func (l *LoggerDecorated) WithField(key string, value interface{}) interfaces.Logger {
	return l.next.WithField(key, value)
}

func (l *LoggerDecorated) WithFields(fields ...interfaces.LogField) interfaces.Logger {
	return l.next.WithFields(fields...)
}

func (l *LoggerDecorated) WithError(err error) interfaces.Logger {
	return l.next.WithError(err)
}

func (l *LoggerDecorated) Write(p []byte) (n int, err error) {
	return l.next.Write(p)
}
