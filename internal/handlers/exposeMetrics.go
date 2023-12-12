package handlers

import (
	"time"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/jtonynet/cine-catalogo/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
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

func ExposeMetrics(r *gin.Engine, cfg config.API) {
	metricsLoop()

	// Use OpenTelemetry middleware
	r.Use(otelgin.Middleware(cfg.Name))

	// Adding router metrics with ginprom for all requests routes, find
	// another lib with more uses or better way to get this metrics
	// https://github.com/Depado/ginprom/tree/master
	p := ginprom.New(
		ginprom.Engine(r),
		ginprom.Path("/metrics"),
	)
	r.Use(p.Instrument())
}
