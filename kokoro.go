package kokoro

import (
	"context"
	"github.com/caarlos0/env/v11"
	"github.com/kenzo-spaulding/kokoro/telemetry/logs"
	"github.com/kenzo-spaulding/kokoro/telemetry/metrics"
	"github.com/kenzo-spaulding/kokoro/telemetry/traces"
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

func Init(opts ...Option) (Done, error) {
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
			return nil, wrapErr(ErrEnvLoadFailed, err)
		}
	}

	if opt.ctx != nil {
		ctx = opt.ctx
	}

	ctx, cancel := context.WithCancel(ctx)

	err := logs.Init(config.Logs)
	if err != nil {
		cancel()
		return nil, wrapErr(ErrInitializationFailed, err)
	}

	err = metrics.Init(config.Metrics)
	if err != nil {
		cancel()
		return nil, wrapErr(ErrInitializationFailed, err)
	}

	err = traces.Init(ctx, config.Traces)
	if err != nil {
		cancel()
		return nil, wrapErr(ErrInitializationFailed, err)
	}

	done := func() {
		cancel()
	}

	return done, nil
}
