package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"net/http"
	"time"
)

var DefaultFactory Factory

type Metrics struct {
	MetricsPort int `env:"METRICS_PORT" envDefault:"8000"`
}

type Factory interface {
	NewCounter(name string, opts ...MetricOption) (Counter, error)
	NewHistogram(name string, opts ...MetricOption) (Histogram, error)
}

type DefaultMetricsFactory struct {
	meter      api.Meter
	counters   map[string]Counter
	histograms map[string]Histogram
	gauges     map[string]Gauge
}

func Init(config Metrics) error {
	exporter, err := prometheus.New()
	if err != nil {
		return fmt.Errorf("failed to load prometheus exporter: %w", err)
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("github.com/kenzo-spaulding/kokoro")

	DefaultFactory = &DefaultMetricsFactory{
		meter:      meter,
		counters:   make(map[string]Counter),
		histograms: make(map[string]Histogram),
		gauges:     make(map[string]Gauge),
	}

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/", promhttp.Handler())
		server := &http.Server{
			Addr:              fmt.Sprintf(":%d", config.MetricsPort),
			Handler:           mux,
			ReadTimeout:       30 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       360 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			MaxHeaderBytes:    1 << 20, // 1 MB
		}

		err := server.ListenAndServe()
		if err != nil {
			log.Panic().Err(err).Int("port", config.MetricsPort).
				Msg("failed to serve/failed while serving metrics")
		}
	}()

	return nil
}
