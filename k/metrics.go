package k

import "github.com/kzs0/kokoro/telemetry/metrics"

func NewCounter(name string, opts ...metrics.MetricOption) (metrics.Counter, error) {
	return metrics.DefaultFactory.NewCounter(name, opts...)
}

func NewHistogram(name string, opts ...metrics.MetricOption) (metrics.Histogram, error) {
	return metrics.DefaultFactory.NewHistogram(name, opts...)
}

func NewGauge(name string, opts ...metrics.MetricOption) (metrics.Gauge, error) {
	return metrics.DefaultFactory.NewGauge(name, opts...)
}
