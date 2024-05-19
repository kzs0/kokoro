module github.com/kzs0/kokoro

go 1.22.3

require (
	github.com/caarlos0/env/v11 v11.0.0
	github.com/prometheus/client_golang v1.19.0
	github.com/rs/zerolog v1.32.0
	go.opentelemetry.io/otel v1.26.0
	go.opentelemetry.io/otel/exporters/prometheus v0.48.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.26.0
	go.opentelemetry.io/otel/metric v1.26.1-0.20240518090000-14441aefdf88
	go.opentelemetry.io/otel/sdk v1.26.0
	go.opentelemetry.io/otel/sdk/metric v1.26.0
	go.opentelemetry.io/otel/trace v1.26.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
