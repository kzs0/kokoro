package kokoro

import (
	"context"
	"errors"

	"github.com/kzs0/kokoro/env"
	"github.com/kzs0/kokoro/telemetry/logs"
	"github.com/kzs0/kokoro/telemetry/metrics"
	"github.com/kzs0/kokoro/telemetry/traces"
)

type options struct {
	ctx    context.Context
	config Config
}

type Option func(*options)
type Done func()

func WithConfig(config Config) Option {
	return func(o *options) {
		o.config = config
	}
}

func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func Init(opts ...Option) (context.Context, Done, error) {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}

	config := opt.config
	def := Config{}
	ctx := context.Background()

	if opt.config == def {
		err := env.Parse(&config)
		if err != nil {
			return ctx, nil, errors.Join(ErrEnvLoadFailed, err)
		}
	}

	if opt.ctx != nil {
		ctx = opt.ctx
	}

	ctx, cancel := context.WithCancel(ctx)

	err := logs.Init(config.Logs)
	if err != nil {
		cancel()
		return ctx, nil, errors.Join(ErrInitializationFailed, err)
	}

	err = metrics.Init(config.Metrics)
	if err != nil {
		cancel()
		return ctx, nil, errors.Join(ErrInitializationFailed, err)
	}

	err = traces.Init(ctx, config.Traces)
	if err != nil {
		cancel()
		return ctx, nil, errors.Join(ErrInitializationFailed, err)
	}

	done := func() {
		cancel()
	}

	return ctx, done, nil
}
