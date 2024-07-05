package koko

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Attribute func(context.Context) context.Context

func Str(k, s string) Attribute {
	return func(ctx context.Context) context.Context {
		st, ok := getStack(ctx)
		if !ok {
			return ctx
		}

		st.Strs[k] = s

		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.String(k, s))

		return saveStack(ctx, st)
	}
}

func Bool(k string, b bool) Attribute {
	return func(ctx context.Context) context.Context {
		st, ok := getStack(ctx)
		if !ok {
			return ctx
		}

		st.Bools[k] = b

		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.Bool(k, b))

		return saveStack(ctx, st)
	}
}

func intAttr(k string, i int64) Attribute {
	return func(ctx context.Context) context.Context {
		st, ok := getStack(ctx)
		if !ok {
			return ctx
		}

		st.Ints[k] = i

		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.Int64(k, i))

		return saveStack(ctx, st)
	}
}

func Uint8(k string, u uint8) Attribute {
	return intAttr(k, int64(u))
}

func Uint16(k string, u uint16) Attribute {
	return intAttr(k, int64(u))
}

func Uint32(k string, u uint32) Attribute {
	return intAttr(k, int64(u))
}
func Int8(k string, i int8) Attribute {
	return intAttr(k, int64(i))
}

func Int16(k string, i int16) Attribute {
	return intAttr(k, int64(i))
}

func Int32(k string, i int32) Attribute {
	return intAttr(k, int64(i))
}

func Int64(k string, i int64) Attribute {
	return intAttr(k, i)
}

func floatAttr(k string, f float64) Attribute {
	return func(ctx context.Context) context.Context {
		st, ok := getStack(ctx)
		if !ok {
			return ctx
		}

		st.Floats[k] = f

		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.Float64(k, f))

		return saveStack(ctx, st)
	}
}

func Float32(k string, f float32) Attribute {
	return floatAttr(k, float64(f))
}

func Float64(k string, f float64) Attribute {
	return floatAttr(k, f)
}

func Register(ctx context.Context, attrs ...Attribute) context.Context {
	for _, attr := range attrs {
		ctx = attr(ctx)
	}

	return ctx
}
