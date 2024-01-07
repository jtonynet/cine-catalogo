# 3. Prometheus Libs To Instrumentation

Date: 2024-01-06

## Status

Accepted

## Context

It was decided that the project will be designed with observability in mind, connecting to Grafana through Prometheus with visualization areas for `API BASIC`, [API RED](https://www.weave.works/blog/the-red-method-key-metrics-for-microservices-architecture/), and a minimum of `API USE` (only basic panels, nothing too in-depth). The goal of this document is to select a Prometheus client library that provides a coherent set of metrics (such as gin_gonic_http_requests) for the project.

Among the discussed options are:

- [Ginprom](github.com/Depado/ginprom)
  - Few stars (134), last update two months ago
  - Provides the metric gin_gonic_requests_total per route and code
  - Provides the metric gin_gonic_request_duration_sum per route but IGNORES the code

```go
import "github.com/Depado/ginprom"

func ExposeMetrics(r *gin.Engine, cfg config.API) {
	metricsLoop()

	// Use OpenTelemetry middleware
	r.Use(otelgin.Middleware(cfg.Name))

	// Adding router metrics with ginprom for all requests routes, find
	// another lib with more uses or a better way to get these metrics
	// https://github.com/Depado/ginprom/tree/master
	p := ginprom.New(
		ginprom.Engine(r),
		ginprom.Path("/metrics"),
	)
	r.Use(p.Instrument())
}

```

<br/>

- [go-gin-prometheus](github.com/zsais/go-gin-prometheus)
  - Reasonable number of stars for a Go project (420), but last update two years ago
  - Provides metric gin_requests_total with route and code but HIGH CARDINALITY, requires bugfix
  - No metric gin_request_duration_sum
  - Has a horrible workaround to reduce cardinality
```go
import ginprometheus "github.com/zsais/go-gin-prometheus"

func ExposeMetrics(r *gin.Engine, cfg config.API) {
	metricsLoop()

	// Use OpenTelemetry middleware
	r.Use(otelgin.Middleware(cfg.Name))

	p := ginprometheus.NewPrometheus("gin")

	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		url := c.Request.URL.Path
		for _, p := range c.Params {
			if p.Key == "name" {
				url = strings.Replace(url, p.Value, ":name", 1)
				break
			}
		}
		return url
	}

	p.Use(r)
}
```

- [client_golang](https://github.com/prometheus/client_golang) (Standard Library)
  - A large number of stars (5k) and forks (1.2)
  - Everything needs to be configured in it
  - The previous alternatives internally use this option
```go
import promhttp "github.com/prometheus/client_golang/prometheus/promhttp"

func ExposeMetrics(r *gin.Engine, cfg config.API) {
	metricsLoop()

	// Use OpenTelemetry middleware
	r.Use(otelgin.Middleware(cfg.Name))

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
```

## Decision

After extensive research, it has been concluded that using the standard library [client_golang](https://github.com/prometheus/client_golang) is the best decision in this case. The alternatives are wrappers that add metrics, but nothing close to what is seen in the Java/Spring Boot world with [Micrometer](https://micrometer.io/). With a focus on the low footprint of the Go world, the trade-off is to create the necessary metrics in a customized way, generating a higher mental load for developers but maintaining the "Go Way."

Using the following links and repositories as support for this decision:

- https://dev.to/eminetto/using-prometheus-to-collect-metrics-from-golang-applications-35gc
- https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/
- https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin/blob/main/chapter10/main.go#L47 from book [Building Distributed Applications in Gin](https://www.amazon.com.br/Building-Distributed-Applications-Gin-hands-ebook/dp/B091G3DBRT)

## Consequences

Go is lightweight and modularized; consequently, facilities similar to those in other languages and frameworks, such as a library similar to Micrometer, are rare. We accept this trade-off and move towards what the community advocates in practice. As observed during the study for this decision, wrapper libraries for the standard Prometheus client library did not gain traction, so it does not make sense to use them. This increases the mental load but provides control over our metrics.
