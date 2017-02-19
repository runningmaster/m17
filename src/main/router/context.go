package router

import (
	"context"
)

// contextParamKey is a unique type to prevent assignment.
// Its associated value should be a string.
type contextParamKey struct {
	name string
}

// WithParamValue returns a new context based on the provided parent ctx.
func WithParamValue(ctx context.Context, key, val string) context.Context {
	return context.WithValue(ctx, contextParamKey{key}, val)
}

// ContextParamValue returns the first value associated with the given key.
func ContextParamValue(ctx context.Context, key string) string {
	ctxKey := contextParamKey{key}
	if v, ok := ctx.Value(ctxKey).(string); ok {
		return v
	}
	return ""
}

// contextQueryKey is a unique type to prevent assignment.
// Its associated value should be a string.
type contextQueryKey struct {
	name string
}

// WithQueryValue returns a new context based on the provided parent ctx.
func WithQueryValue(ctx context.Context, key, val string) context.Context {
	return context.WithValue(ctx, contextQueryKey{key}, val)
}

// ContextQueryValue returns the first value associated with the given key.
func ContextQueryValue(ctx context.Context, key string) string {
	ctxKey := contextQueryKey{key}
	if v, ok := ctx.Value(ctxKey).(string); ok {
		return v
	}
	return ""
}
