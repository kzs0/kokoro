package k

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kzs0/kokoro/telemetry/metrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type recorder struct {
	successes metrics.Counter
	failures  metrics.Counter
	count     metrics.Counter
	timer     metrics.Histogram
}

func (r *recorder) AddLabels(opts ...metrics.MeasurementOption) {
	r.successes.Load(opts...)
	r.failures.Load(opts...)
	r.count.Load(opts...)
	r.timer.Load(opts...)
}

func (r *recorder) Record(ctx context.Context, dur time.Duration, success bool) error {
	if success {
		successes, err := NewCounter(fmt.Sprintf("%s_success"))
		if err != nil {
			return err
		}

		err = successes.Incr(ctx)
		if err != nil {
			return err
		}

	} else {
		failures, err := NewCounter(fmt.Sprintf("%s_failures"))
		if err != nil {
			return err
		}

		err = failures.Incr(ctx)
		if err != nil {
			return err
		}
	}

	ops, err := NewCounter(fmt.Sprintf("%s_count"))
	if err != nil {
		return err
	}

	err = ops.Incr(ctx)
	if err != nil {
		return err
	}

	timer, err := NewHistogram(fmt.Sprintf("%s_millis"))
	err = timer.Record(ctx, float64(dur.Milliseconds()))
	if err != nil {
		return err
	}

	return nil
}

func newRecorder() (*recorder, error) {
	successes, err := NewCounter(fmt.Sprintf("%s_success"))
	if err != nil {
		return nil, err
	}

	failures, err := NewCounter(fmt.Sprintf("%s_failures"))
	if err != nil {
		return nil, err
	}

	count, err := NewCounter(fmt.Sprintf("%s_count"))
	if err != nil {
		return nil, err
	}

	timer, err := NewHistogram(fmt.Sprintf("%s_millis"))
	if err != nil {
		return nil, err
	}

	return &recorder{
		successes: successes,
		failures:  failures,
		count:     count,
		timer:     timer,
	}, nil
}

type Done func(*context.Context, *error)

func Operation(ctx context.Context, operation string) (context.Context, Done) {
	ctx = initStack(ctx)
	start := time.Now()

	tracer := otel.Tracer("kzs0/kokoro")
	ctx, _ = tracer.Start(ctx, operation)

	r, err := newRecorder()
	if err != nil {
		log.Debug().Err(err).Msg("failed to create metrics")
		return ctx, func(ctx *context.Context, err *error) {}
	}

	done := func(ctx *context.Context, err *error) {
		stop := time.Since(dur)

		st, ok := pop(*ctx)
		if !ok {
			return
		}

		var level zerolog.Level
		level, lerr := zerolog.ParseLevel(strings.ToLower(st.LogLevel))
		if lerr != nil {
			log.Debug().Str("log_level", strings.ToLower(st.LogLevel)).
				Msg("failed to parse log level, using defaults")
			level = zerolog.DebugLevel
		}

		if *err != nil && zerolog.WarnLevel > level {
			level = zerolog.WarnLevel
		}

		span := trace.SpanFromContext(*ctx)
		span.SetStatus(codes.Ok, "success")

		logs := log.WithLevel(level).
			Dur("duration", time.Since(start)).
			Str("operation", operation)

		for k, f := range st.Floats {
			logs = logs.Float64(k, f)
			r.AddLabels(metrics.WithLabel(k, fmt.Sprint(f)))
		}
		for k, i := range st.Ints {
			logs = logs.Int64(k, i)
			r.AddLabels(metrics.WithLabel(k, fmt.Sprintf(i)))
		}
		for k, s := range st.Strs {
			logs = logs.Str(k, s)
			r.AddLabels(metrics.WithLabel(k, s))
		}
		for k, b := range st.Bools {
			logs = logs.Bool(k, b)
			r.AddLabels(metrics.WithLabel(k, fmt.Sprintf(b)))
		}

		if *err != nil {
			logs = logs.Err(*err)
			span.SetStatus(codes.Error, "error encountered")
		}

		logs.Msg(operation)
		span.End()

		rerr := r.Record(*ctx, stop, *err == nil)
		if rerr != nil {
			log.Debug.Msg("failed to record metrics for operation")
		}
	}

	return ctx, done
}