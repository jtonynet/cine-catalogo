package logger

import (
	"io"
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// TODO: converts logEventsTotal to decorator PrometheusLog
var (
	logEventsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "log_events_total",
			Help: "Total number of log events.",
		},
		[]string{"level"},
	)
)

type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
	writer  io.Writer
}

func NewLogger(p string) *Logger {
	writer := io.Writer(os.Stdout)
	logger := log.New(writer, p, log.Ldate|log.Ltime)

	return &Logger{
		debug:   log.New(writer, ">> DEBUG ", logger.Flags()),
		info:    log.New(writer, ">> INFO ", logger.Flags()),
		warning: log.New(writer, ">> WARNING ", logger.Flags()),
		error:   log.New(writer, ">> ERROR ", logger.Flags()),
		writer:  writer,
	}
}

func (l *Logger) Debug(v ...interface{}) {
	l.debug.Println(v...)
	logEventsTotal.WithLabelValues("debug").Inc()
}

func (l *Logger) Info(v ...interface{}) {
	l.info.Println(v...)
	logEventsTotal.WithLabelValues("info").Inc()
}

func (l *Logger) Warning(v ...interface{}) {
	l.warning.Println(v...)
	logEventsTotal.WithLabelValues("warning").Inc()
}

func (l *Logger) Error(v ...interface{}) {
	l.error.Println(v...)
	logEventsTotal.WithLabelValues("error").Inc()
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.debug.Printf(format, v...)
	logEventsTotal.WithLabelValues("debug").Inc()
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.info.Printf(format, v...)
	logEventsTotal.WithLabelValues("info").Inc()
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.warning.Printf(format, v...)
	logEventsTotal.WithLabelValues("warning").Inc()
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.error.Printf(format, v...)
	logEventsTotal.WithLabelValues("error").Inc()
}
