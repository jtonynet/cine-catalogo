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

func (l *LoggerDecorated) Debug(v ...interface{}) {
	l.next.Debug(v...)

	logEventsTotal.WithLabelValues("debug").Inc()
}

func (l *LoggerDecorated) Info(v ...interface{}) {
	l.next.Info(v...)

	logEventsTotal.WithLabelValues("info").Inc()
}

func (l *LoggerDecorated) Warning(v ...interface{}) {
	l.next.Warning(v...)

	logEventsTotal.WithLabelValues("warning").Inc()
}

func (l *LoggerDecorated) Error(v ...interface{}) {
	l.next.Error(v...)

	logEventsTotal.WithLabelValues("error").Inc()
}

func (l *LoggerDecorated) Debugf(format string, v ...interface{}) {
	l.next.Debugf(format, v...)

	logEventsTotal.WithLabelValues("debug").Inc()
}

func (l *LoggerDecorated) Infof(format string, v ...interface{}) {
	l.next.Infof(format, v...)

	logEventsTotal.WithLabelValues("info").Inc()
}

func (l *LoggerDecorated) Warningf(format string, v ...interface{}) {
	l.next.Warningf(format, v...)

	logEventsTotal.WithLabelValues("warning").Inc()
}

func (l *LoggerDecorated) Errorf(format string, v ...interface{}) {
	l.next.Errorf(format, v...)

	logEventsTotal.WithLabelValues("error").Inc()
}
