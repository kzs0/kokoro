package metrics

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Histogram interface {
	Loadable

	// Record will observe the measurement
	Record(ctx context.Context, measurement float64, opts ...MeasurementOption) error
}

type defaultHistogram struct {
	histogram    metric.Float64Histogram
	staticLabels []attribute.KeyValue
	opts         []MeasurementOption
	labelNames   map[string]struct{}
}

func (h *defaultHistogram) Record(ctx context.Context, measurement float64, opts ...MeasurementOption) error {
	if measurement < 0 {
		return fmt.Errorf("measurement cannot be negative")
	}

	opt := metricOpts{}
	for _, o := range opts {
		o(&opt)
	}
	for _, o := range h.opts {
		o(&opt)
	}

	labels := h.staticLabels
	for k, v := range opt.labels {
		if h.labelNames != nil {
			if _, ok := h.labelNames[k]; ok {
				labels = append(labels, attribute.Key(k).String(v))
			}
		}
	}

	h.histogram.Record(ctx, measurement, metric.WithAttributeSet(attribute.NewSet(labels...)))

	return nil
}

func (h *defaultHistogram) Load(opts ...MeasurementOption) {
	h.opts = append(h.opts, opts...)
}

// NewHistogram will produce a Histogram for observing values
//
// It will create a new histogram on first invocation, or return a cached histogram
func (mf *defaultMetricsFactory) NewHistogram(name string, opts ...MetricOption) (Histogram, error) {
	if h, ok := mf.histograms[name]; ok {
		return h, nil
	}

	opt := metricOpts{}
	for _, o := range opts {
		o(&opt)
	}

	if mf.config.ServiceName != "_" {
		name = strings.TrimSpace(strings.ReplaceAll(fmt.Sprintf("%s_%s", mf.config.ServiceName, name), "-", "_"))
	}

	histogram := &defaultHistogram{}

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

	labelNames := make(map[string]struct{})
	if opt.labelNames != nil {
		for _, label := range opt.labelNames {
			labelNames[label] = struct{}{}
		}
	}

	histogram.labelNames = labelNames

	if len(histogram.staticLabels) == 0 {
		histogram.staticLabels = make([]attribute.KeyValue, 0)
	}

	if mf.histograms == nil {
		mf.histograms = make(map[string]Histogram, 1)
	}
	mf.histograms[name] = histogram

	return histogram, nil
}
