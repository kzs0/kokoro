package metrics

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	api "go.opentelemetry.io/otel/sdk/metric"
)

var DefaultFactory Factory

type Metrics struct {
	MetricsPort int    `env:"METRICS_PORT" envDefault:"8000"`
	ServiceName string `env:"SERVICE_NAME" envDefault:"_"`
	Environment string `env:"ENVIRONMENT" envDefault:"dev"`
}

type Factory interface {
	NewCounter(name string, opts ...MetricOption) (Counter, error)
	NewHistogram(name string, opts ...MetricOption) (Histogram, error)
	NewGauge(name string, opts ...MetricOption) (Gauge, error)
}

// Loadable is a behavior where measurement options can be loaded prior to
// measuring with the metric
type Loadable interface {
	// Load will load the Metric with the MeasurementOption provided
	Load(opts ...MeasurementOption)
}

type DefaultMetricsFactory struct {
	config       Metrics
	meter        metric.Meter
	staticLabels map[string]string
	counters     map[string]Counter
	histograms   map[string]Histogram
	gauges       map[string]Gauge
}

func Init(config Metrics, options ...FactoryOption) error {
	opts := factoryOpts{}
	for _, o := range options {
		o(&opts)
	}

	exporter, err := prometheus.New()
	if err != nil {
		return fmt.Errorf("failed to load prometheus exporter: %w", err)
	}

	provider := api.NewMeterProvider(api.WithReader(exporter))
	meter := provider.Meter("github.com/kzs0/kokoro")

	static := map[string]string{
		"service": config.ServiceName,
		"env":     config.Environment,
	}

	for k, v := range opts.staticLabels {
		static[k] = v
	}

	DefaultFactory = &DefaultMetricsFactory{
		meter:        meter,
		counters:     make(map[string]Counter),
		histograms:   make(map[string]Histogram),
		gauges:       make(map[string]Gauge),
		staticLabels: static,
	}

	if opts.factory != nil {
		DefaultFactory = opts.factory
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
			slog.Error("failed to serve/failed while serving metrics",
				slog.String("error", err.Error()), slog.Int("port", config.MetricsPort))

			panic(err)
		}
	}()

	return nil
}
