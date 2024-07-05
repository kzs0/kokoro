package metrics

type factoryOpts struct {
	staticLabels map[string]string
	factory      Factory
}

type FactoryOption func(*factoryOpts)

// WithStaticLabel allows setting labels that will be set on all metrics
// created with the factory
func WithStaticLabel(label, value string) FactoryOption {
	return func(f *factoryOpts) {
		if f.staticLabels == nil {
			f.staticLabels = make(map[string]string)
		}

		f.staticLabels[label] = value
	}
}

// WithFactory allows providing a custom factory to be used as the DefaultFactory
func WithFactory(factory Factory) FactoryOption {
	return func(f *factoryOpts) {
		f.factory = factory
	}
}

type metricOpts struct {
	desc         string
	unit         string
	staticLabels map[string]string
	labels       map[string]string
	labelNames   []string
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

// WithHistogramBucketsBounds allows to override the default bucket boundaries for a histogram
func WithHistogramBucketsBounds(buckets ...float64) MetricOption {
	return func(opts *metricOpts) {
		opts.buckets = buckets
	}
}

// WithLabelNames sets the labels expected to be provided to the metric.
//
// Subsequent WithLabelNames will overwrite the previous set of names passed in.
// Labels passed in that were not provided as a LabelName will be ignored.
// Labels not passed in that were expected will result in an error being returned. // TODO <- This could also just fill in -?
func WithLabelNames(labels []string) MetricOption {
	return func(opts *metricOpts) {
		opts.labelNames = labels
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
