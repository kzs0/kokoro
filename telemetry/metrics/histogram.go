package metrics

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Histogram interface {
	// Record will observe the measurement
	Record(ctx context.Context, measurement float64, opts ...MeasurementOption) error

	// Curry will curry the Histogram with the MeasurementOption provided
	Curry(opts ...MeasurementOption) Histogram
}

type DefaultHistogram struct {
	histogram    metric.Float64Histogram
	staticLabels []attribute.KeyValue
	opts         []MeasurementOption
}

func (h *DefaultHistogram) Record(ctx context.Context, measurement float64, opts ...MeasurementOption) error {
	if measurement < 0 {
		return fmt.Errorf("measurement cannot be negative")
	}

	opt := metricOpts{}
	for _, o := range opts {
		o(&opt)
	}

	labels := h.staticLabels
	for k, v := range opt.labels {
		labels = append(labels, attribute.Key(k).String(v))
	}

	h.histogram.Record(ctx, measurement, metric.WithAttributeSet(attribute.NewSet(labels...)))

	return nil
}

func (h *DefaultHistogram) Curry(opts ...MeasurementOption) Histogram {
	h.opts = append(h.opts, opts...)
	return h
}

// NewHistogram will produce a Histogram for observing values
//
// It will create a new histogram on first invocation, or return a cached histogram
func (mf *DefaultMetricsFactory) NewHistogram(name string, opts ...MetricOption) (Histogram, error) {
	if h, ok := mf.histograms[name]; ok {
		return h, nil
	}

	opt := metricOpts{}
	for _, o := range opts {
		o(&opt)
	}

	histogram := &DefaultHistogram{}

	otelOpts := make([]metric.Float64HistogramOption, 0)
	if opt.desc != "" {
		otelOpts = append(otelOpts, metric.WithDescription(opt.desc))
	}
	if opt.unit != "" {
		otelOpts = append(otelOpts, metric.WithUnit(opt.unit))
	}
	if len(opt.buckets) > 0 {
		otelOpts = append(otelOpts, metric.WithExplicitBucketBoundaries(opt.buckets...))
	}
	if len(opt.staticLabels) > 0 {
		attr := make([]attribute.KeyValue, len(opt.staticLabels))
		for k, v := range opt.staticLabels {
			attr = append(attr, attribute.Key(k).String(v))
		}
		histogram.staticLabels = attr
	}

	otelHistogram, err := mf.meter.Float64Histogram(name, otelOpts...)
	if err != nil {
		return nil, err
	}

	histogram.histogram = otelHistogram
	histogram.opts = make([]MeasurementOption, 0)

	if len(histogram.staticLabels) == 0 {
		histogram.staticLabels = make([]attribute.KeyValue, 0)
	}

	if mf.histograms == nil {
		mf.histograms = make(map[string]Histogram, 1)
	}
	mf.histograms[name] = histogram

	return histogram, nil
}
