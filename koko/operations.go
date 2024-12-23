package koko

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"time"

	"github.com/kzs0/kokoro/telemetry/logs"
	"github.com/kzs0/kokoro/telemetry/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracerName string = "kzs0/kokoro"

type recorder struct {
	operation string
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
		successes, err := Counter(fmt.Sprintf("%s_success", r.operation))
		if err != nil {
			return err
		}

		err = successes.Incr(ctx)
		if err != nil {
			return err
		}
	} else {
		failures, err := Counter(fmt.Sprintf("%s_failures", r.operation))
		if err != nil {
			return err
		}

		err = failures.Incr(ctx)
		if err != nil {
			return err
		}
	}

	ops, err := Counter(fmt.Sprintf("%s_count", r.operation))
	if err != nil {
		return err
	}

	err = ops.Incr(ctx)
	if err != nil {
		return err
	}

	timer, err := Histogram(fmt.Sprintf("%s_millis", r.operation))
	err = timer.Record(ctx, float64(dur.Milliseconds()))
	if err != nil {
		return err
	}

	return nil
}

func newRecorder(op string, opts ...metrics.MetricOption) (*recorder, error) {
	successes, err := Counter(fmt.Sprintf("%s_success", op), opts...)
	if err != nil {
		return nil, err
	}

	failures, err := Counter(fmt.Sprintf("%s_failures", op), opts...)
	if err != nil {
		return nil, err
	}

	count, err := Counter(fmt.Sprintf("%s_count", op), opts...)
	if err != nil {
		return nil, err
	}

	timer, err := Histogram(fmt.Sprintf("%s_millis", op), opts...)
	if err != nil {
		return nil, err
	}

	return &recorder{
		operation: op,
		successes: successes,
		failures:  failures,
		count:     count,
		timer:     timer,
	}, nil
}

type Done func(*context.Context, *error)

type DoneNoErr func(*context.Context)

// Operation will bootstrap a short lived code path and report traces, metrics,
// and logs automatically.
//
// An operation is assumed to have some failure condition due to side effects.
func Operation(ctx context.Context, operation string, opts ...metrics.MetricOption) (context.Context, Done) {
	ctx = initStack(ctx)
	start := time.Now()

	tracer := otel.Tracer(tracerName)
	ctx, _ = tracer.Start(ctx, operation)

	r, err := newRecorder(operation, opts...)
	if err != nil {
		slog.Warn("failed to create metrics", slog.String("error", err.Error()))
		return ctx, func(ctx *context.Context, err *error) {}
	}

	done := func(ctx *context.Context, err *error) {
		stop := time.Since(start)

		st, ok := pop(*ctx)
		if !ok {
			return
		}

		if err == nil {
			var perr error
			err = &perr
		}

		var level slog.Level
		level, lerr := logs.ParseLevel(st.LogLevel)
		if lerr != nil {
			slog.Debug("failed to parse log level, using default",
				slog.String("log_level", strings.ToUpper(st.LogLevel)))
			level = slog.LevelDebug
		}

		if *err != nil && slog.LevelWarn > level {
			level = slog.LevelWarn
		}

		span := trace.SpanFromContext(*ctx)
		span.SetStatus(codes.Error, "error encountered")

		if *err == nil {
			// OK > Error so this will overwrite the previous status
			span.SetStatus(codes.Ok, "success")
		}

		attrs := []slog.Attr{
			slog.Duration("duration", time.Since(start)),
			slog.String("operation", operation),
		}

		*ctx = Register(*ctx, Int64("duration",
			int64(time.Since(start).Milliseconds())), Str("operation", operation))

		for k, f := range st.Floats {
			attrs = append(attrs, slog.Float64(k, f))
			r.AddLabels(metrics.WithLabel(k, fmt.Sprint(f)))
		}
		for k, i := range st.Ints {
			attrs = append(attrs, slog.Int64(k, i))
			r.AddLabels(metrics.WithLabel(k, fmt.Sprint(i)))
		}
		for k, s := range st.Strs {
			attrs = append(attrs, slog.String(k, s))
			r.AddLabels(metrics.WithLabel(k, s))
		}
		for k, b := range st.Bools {
			attrs = append(attrs, slog.Bool(k, b))
			r.AddLabels(metrics.WithLabel(k, fmt.Sprint(b)))
		}

		if *err != nil {
			attrs = append(attrs, slog.String("error", (*err).Error()))
			span.RecordError(*err)
		}

		slog.LogAttrs(*ctx, level, operation, attrs...)
		span.End()

		rerr := r.Record(*ctx, stop, *err == nil)
		if rerr != nil {
			slog.Debug("failed to record metrics for operation",
				slog.String("operation", operation))
		}
	}

	return ctx, done
}

func getCallerName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "span"
	}

	funcDetails := runtime.FuncForPC(pc)
	if funcDetails == nil {
		return "span"
	}

	return funcDetails.Name()
}

// Pure will initiate a new span that cannot encounter an error during
// operation
func Pure(ctx context.Context) (context.Context, DoneNoErr) {
	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, getCallerName())

	done := func(ctx *context.Context) {
		span.SetStatus(codes.Ok, "success")
		span.End()
	}

	return ctx, done
}

// Impure will initiate a new span that can encounter an error during
// operation
func Impure(ctx context.Context) (context.Context, Done) {
	tracer := otel.Tracer(tracerName)
	ctx, span := tracer.Start(ctx, getCallerName())

	done := func(ctx *context.Context, err *error) {
		if *err == nil {
			span.SetStatus(codes.Ok, "success")
		} else {
			span.SetStatus(codes.Error, "error encountered")
			span.RecordError(*err)
		}
		span.End()
	}

	return ctx, done
}
