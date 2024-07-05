package metrics

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Gauge interface {
	Loadable

	// Measure will set the Gauge to the provided value
	Measure(ctx context.Context, value float64, opts ...MeasurementOption) error
}

type DefaultGauge struct {
	gauge        metric.Float64Gauge
	staticLabels []attribute.KeyValue
	opts         []MeasurementOption
	labelNames   map[string]struct{}
}

func (g *DefaultGauge) Measure(ctx context.Context, value float64, opts ...MeasurementOption) error {
	opt := metricOpts{}
	for _, o := range opts {
		o(&opt)
	}

	labels := g.staticLabels
	for k, v := range opt.labels {
		if g.labelNames != nil {
			if _, ok := g.labelNames[k]; ok {
				labels = append(labels, attribute.Key(k).String(v))
			}
		}
	}

	g.gauge.Record(ctx, value, metric.WithAttributeSet(attribute.NewSet(labels...)))

	return nil
}

func (g *DefaultGauge) Load(opts ...MeasurementOption) {
	g.opts = append(g.opts, opts...)
}

// NewGauge will produce a Gauge for setting an instantaneous value
//
// It will create a new gauge on first invocation, or return a cached gauge
// previously created by name
func (mf *DefaultMetricsFactory) NewGauge(name string, opts ...MetricOption) (Gauge, error) {
	if g, ok := mf.gauges[name]; ok {
		return g, nil
	}

	opt := metricOpts{}
	for _, o := range opts {
		o(&opt)
	}

	name = strings.TrimSpace(strings.ReplaceAll(fmt.Sprintf("%s_%s", mf.config.ServiceName, name), "-", "_"))

	gauge := &DefaultGauge{}

	otelOpts := make([]metric.Float64GaugeOption, 0)
	if opt.desc != "" {
		otelOpts = append(otelOpts, metric.WithDescription(opt.desc))
	}
	if opt.unit != "" {
		otelOpts = append(otelOpts, metric.WithUnit(opt.unit))
	}
	if len(opt.staticLabels) > 0 {
		attr := make([]attribute.KeyValue, len(opt.staticLabels))
		for k, v := range opt.staticLabels {
			attr = append(attr, attribute.Key(k).String(v))
		}
		gauge.staticLabels = attr
	}

	otelGauge, err := mf.meter.Float64Gauge(name, otelOpts...)
	if err != nil {
		return nil, err
	}

	gauge.gauge = otelGauge
	gauge.opts = make([]MeasurementOption, 0)

	labelNames := make(map[string]struct{})
	if opt.labelNames != nil {
		for _, label := range opt.labelNames {
			labelNames[label] = struct{}{}
		}
	}

	gauge.labelNames = labelNames

	if len(gauge.staticLabels) == 0 {
		gauge.staticLabels = make([]attribute.KeyValue, 0)
	}

	if mf.gauges == nil {
		mf.gauges = make(map[string]Gauge, 1)
	}
	mf.gauges[name] = gauge

	return gauge, nil
}
