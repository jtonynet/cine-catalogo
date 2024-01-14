# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]
### Added
- Added [Testfy](github.com/stretchr/testify)
- Added "Happy Path" Integration Succesful Routes tests
- Added debug and tests section on readme.md doc


### Fixed



---

## [0.1.1] - 2024-01-14
### Added
- Added [Testfy](github.com/stretchr/testify)
- Added "Happy Path" Integration Succesful Routes tests
- Added debug and tests section on readme.md doc

---

## [0.1.0] - 2024-01-06
### Added

- Added [Zap Logger](https://github.com/uber-go/zap) for log management.
- Added [gopsutil](https://github.com/shirou/gopsutil) for capturing system data for metric collection.
- Introducing observability with metrics using [Prometheus](https://github.com/prometheus/client_golang).
- Log Decorator for log metric management.
- Added [Prometheus Plugin](https://github.com/go-gorm/prometheus) to Gorm for database metrics.
- Included [Grafana Dashboard Script](./scripts/grafana-dashboards/dash-catalogo-api.json) in the project with panel rows: Basic, [RED](https://www.weave.works/blog/the-red-method-key-metrics-for-microservices-architecture/), and USE (only superficially).
- Installed the OpenTelemetry library [otelgin](https://github.com/open-telemetry/opentelemetry-go-contrib).
- `exposeMetrics` handler with custom metrics in a "micrometer-like" fashion for dashboard consumption in [Grafana](https://grafana.com/).
  - process_uptime_seconds
  - system_cpu_usage
  - process_cpu_usage
  - system_load_average_1m
  - system_cpu_cores
  - memory_alloc_bytes
  - memory_total_alloc_bytes
  - memory_used_bytes
  - memory_limit_bytes
  - gin_gonic_requests_total
  - gin_gonic_request_duration
  - gin_gonic_requests_seconds_max

---

## [0.0.1] - 2023-12-06
### Added

- Basic CRUD for Addresses, Cinema Rooms, Movies, and Posters added in [HATEOAS](https://en.wikipedia.org/wiki/HATEOAS) using [HAL](https://github.com/toedter/hal-explorer).
- Functional [Swagger](https://github.com/swaggo/gin-swagger) for the added routes.
- Added [Viper](https://github.com/spf13/viper) dotEnv for managing environment variables.

---

## [0.0.0] - 2023-09-23
### Added

- We start the project [Project](https://github.com/users/jtonynet/projects/2) with initial commit. Base docs: Rich Readme, ADR [0002: Gin, Gorm and Postegres in three tier architecture](./assets/architecture/decisions/0002-gin-gorm-and-postgres-in-three-tier-architecture.md), [Miro Event Storming](https://miro.com/app/board/uXjVNRofMoA=/) and [Mermaid Diagrams](https://github.com/jtonynet/cine-catalogo/tree/main#diagrams)

[0.1.1]: https://github.com/jtonynet/cine-catalogo/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/jtonynet/cine-catalogo/compare/v0.0.1...v0.1.0
[0.0.1]: https://github.com/jtonynet/cine-catalogo/compare/v0.0.0...v0.0.1
[0.0.0]: https://github.com/jtonynet/cine-catalogo/releases/tag/v0.0.0
