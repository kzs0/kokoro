package metrics

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Counter interface {
	Loadable

	// Incr will increment the counter by 1
	Incr(ctx context.Context, opts ...MeasurementOption) error

	// Add will add the given addend to the counter
	Add(ctx context.Context, addend float64, opts ...MeasurementOption) error
}

type defaultCounter struct {
	counter      metric.Float64Counter
	staticLabels []attribute.KeyValue
	opts         []MeasurementOption
	labelNames   map[string]struct{}
}

func (c *defaultCounter) Incr(ctx context.Context, opts ...MeasurementOption) error {
	return c.Add(ctx, 1, opts...)
}

func (c *defaultCounter) Add(ctx context.Context, addend float64, opts ...MeasurementOption) error {
	if addend < 0 {
		return fmt.Errorf("addend cannot be negative")
	}

	opt := metricOpts{}
	for _, o := range opts {
		o(&opt)
	}
	for _, o := range c.opts {
		o(&opt)
	}

	labels := c.staticLabels
	for k, v := range opt.labels {
		if c.labelNames != nil {
			if _, ok := c.labelNames[k]; ok {
				labels = append(labels, attribute.Key(k).String(v))
			}
		}
	}

	c.counter.Add(ctx, addend, metric.WithAttributeSet(attribute.NewSet(labels...)))

	return nil
}

func (c *defaultCounter) Load(opts ...MeasurementOption) {
	c.opts = append(c.opts, opts...)
}

// NewCounter will produce a Counter for measuring values that go up
//
// It will create a new counter on first invocation, or return a cached counter
// previously created by name
func (mf *defaultMetricsFactory) NewCounter(name string, opts ...MetricOption) (Counter, error) {
	if c, ok := mf.counters[name]; ok {
		return c, nil
	}

	opt := metricOpts{}
	for _, o := range opts {
		o(&opt)
	}

	if mf.config.ServiceName != "_" {
		name = strings.TrimSpace(strings.ReplaceAll(fmt.Sprintf("%s_%s", mf.config.ServiceName, name), "-", "_"))
	}

	counter := &defaultCounter{}

	otelOpts := make([]metric.Float64CounterOption, 0)
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
		counter.staticLabels = attr
	}

	otelCounter, err := mf.meter.Float64Counter(name, otelOpts...)
	if err != nil {
		return nil, err
	}

	counter.counter = otelCounter
	counter.opts = make([]MeasurementOption, 0)

	labelNames := make(map[string]struct{})
	if opt.labelNames != nil {
		for _, label := range opt.labelNames {
			labelNames[label] = struct{}{}
		}
	}

	counter.labelNames = labelNames

	if len(counter.staticLabels) == 0 {
		counter.staticLabels = make([]attribute.KeyValue, 0)
	}

	if mf.counters == nil {
		mf.counters = make(map[string]Counter, 1)
	}
	mf.counters[name] = counter

	return counter, nil
}
