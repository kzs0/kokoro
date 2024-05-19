package k

import (
	"context"
)

type stack struct {
	Strs     map[string]string
	Ints     map[string]int64
	Floats   map[string]float64
	Bools    map[string]bool
	LogLevel string
}

type key int

var stackKey key

func initStack(ctx context.Context) context.Context {
	st := stack{
		Strs:     make(map[string]string),
		Ints:     make(map[string]int64),
		Floats:   make(map[string]float64),
		Bools:    make(map[string]bool),
		LogLevel: "DEBUG",
	}

	return context.WithValue(ctx, stackKey, &st)
}

func getStack(ctx context.Context) (stack, bool) {
	st, ok := ctx.Value(stackKey).(stack)
	if !ok {
		return stack{}, false
	}

	return st, true
}

func saveStack(ctx context.Context, st stack) context.Context {
	return context.WithValue(ctx, stackKey, st)
}

func pop(ctx context.Context) (stack, bool) {
	st, ok := ctx.Value(stackKey).(stack)
	return st, ok
}
