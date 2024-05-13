package metrics

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Gauge is under construction, open telemetry recently added support synchronous
// gauges to the spec, and a PR to add support to the Go library is open.
// This will be added once that PR is merged and released.
type Gauge interface {
}

type DefaultGauge struct {
	gauge        metric.Float64ObservableGauge
	staticLabels []attribute.KeyValue
	opts         []MeasurementOption
}

func (mf *DefaultMetricsFactory) NewGauge(name string, opts ...MetricOption) (Gauge, error) {
	if g, ok := mf.gauges[name]; ok {
		return g, nil
	}

	opt := metricOpts{}
	for _, o := range opts {
		o(&opt)
	}

	gauge := &DefaultGauge{}

	otelOpts := make([]metric.Float64ObservableGaugeOption, 0)
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

	otelGauge, err := mf.meter.Float64ObservableGauge(name, otelOpts...)
	if err != nil {
		return nil, err
	}

	gauge.gauge = otelGauge
	gauge.opts = make([]MeasurementOption, 0)

	if len(gauge.staticLabels) == 0 {
		gauge.staticLabels = make([]attribute.KeyValue, 0)
	}

	if mf.gauges == nil {
		mf.gauges = make(map[string]Gauge, 1)
	}
	mf.gauges[name] = gauge

	return gauge, nil
}
