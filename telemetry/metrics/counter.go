package metrics

import (
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"golang.org/x/net/context"
)

type Counter interface {
	// Incr will increment the counter by 1
	Incr(ctx context.Context, opts ...MeasurementOption) error

	// Add will add the given addend to the counter
	Add(ctx context.Context, addend float64, opts ...MeasurementOption) error

	// Curry will curry the Counter with the MeasurementOption provided
	Curry(opts ...MeasurementOption) Counter
}

type DefaultCounter struct {
	counter      metric.Float64Counter
	staticLabels []attribute.KeyValue
	opts         []MeasurementOption
}

func (c *DefaultCounter) Incr(ctx context.Context, opts ...MeasurementOption) error {
	return c.Add(ctx, 1, opts...)
}

func (c *DefaultCounter) Add(ctx context.Context, addend float64, opts ...MeasurementOption) error {
	if addend < 0 {
		return fmt.Errorf("addend cannot be negative")
	}

	opt := metricOpts{}
	for _, o := range opts {
		o(&opt)
	}

	labels := c.staticLabels
	for k, v := range opt.labels {
		labels = append(labels, attribute.Key(k).String(v))
	}

	c.counter.Add(ctx, addend, metric.WithAttributeSet(attribute.NewSet(labels...)))

	return nil
}

func (c *DefaultCounter) Curry(opts ...MeasurementOption) Counter {
	c.opts = append(c.opts, opts...)
	return c
}

// NewCounter will produce a Counter for measuring values that go up
//
// It will create a new counter on first invocation, or return a cached counter
// previously created by name
func (mf *DefaultMetricsFactory) NewCounter(name string, opts ...MetricOption) (Counter, error) {
	if c, ok := mf.counters[name]; ok {
		return c, nil
	}

	opt := metricOpts{}
	for _, o := range opts {
		o(&opt)
	}

	counter := &DefaultCounter{}

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

	if len(counter.staticLabels) == 0 {
		counter.staticLabels = make([]attribute.KeyValue, 0)
	}

	if mf.counters == nil {
		mf.counters = make(map[string]Counter, 1)
	}
	mf.counters[name] = counter

	return counter, nil
}
