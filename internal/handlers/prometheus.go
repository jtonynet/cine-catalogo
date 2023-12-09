package handlers

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	startTime = time.Now()

	processUptime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "process_uptime_seconds",
			Help: "Total uptime of the process in seconds",
		},
		[]string{},
	)
)

func metricsLoop() {
	go func() {
		for {
			duration := time.Since(startTime).Seconds()
			processUptime.WithLabelValues().Set(duration)

			time.Sleep(1 * time.Second)
		}
	}()
}

func ExposeMetrics() http.Handler {
	log.Info("handlers: register prometheus metrics route")

	metricsLoop()

	return promhttp.Handler()
}
