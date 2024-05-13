package metrics

type metricOpts struct {
	desc         string
	unit         string
	staticLabels map[string]string
	labels       map[string]string
	buckets      []float64
	factory      Factory
}

type MetricOption func(*metricOpts)
type MeasurementOption func(*metricOpts)

// WithDescription set the description of the metric.
func WithDescription(desc string) MetricOption {
	return func(opts *metricOpts) {
		opts.desc = desc
	}
}

// WithUnit sets units of the measurement.
//
// The unit u should be defined using the appropriate [UCUM](https://ucum.org) case-sensitive code.
func WithUnit(unit string) MetricOption {
	return func(opts *metricOpts) {
		opts.unit = unit
	}
}

// WithStaticLabels set static labels which will always export by the metric
func WithStaticLabels(labels map[string]string) MetricOption {
	return func(opts *metricOpts) {
		opts.staticLabels = labels
	}
}

func WithHistogramBucketsBounds(buckets ...float64) MetricOption {
	return func(opts *metricOpts) {
		opts.buckets = buckets
	}
}

// WithLabel applies a label to the measurement being requested
//
// If multiple WithLabel are applied with the same key, the last entry will be respected
func WithLabel(k, v string) MeasurementOption {
	return func(opts *metricOpts) {
		if opts.labels == nil {
			opts.labels = make(map[string]string)
		}

		opts.labels[k] = v
	}
}
