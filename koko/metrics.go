package koko

import "github.com/kzs0/kokoro/telemetry/metrics"

func Counter(name string, opts ...metrics.MetricOption) (metrics.Counter, error) {
	return metrics.DefaultFactory.NewCounter(name, opts...)
}

func Histogram(name string, opts ...metrics.MetricOption) (metrics.Histogram, error) {
	return metrics.DefaultFactory.NewHistogram(name, opts...)
}

func Gauge(name string, opts ...metrics.MetricOption) (metrics.Gauge, error) {
	return metrics.DefaultFactory.NewGauge(name, opts...)
}
