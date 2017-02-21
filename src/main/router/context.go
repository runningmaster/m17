package router

import (
	"context"
)

// contextParamKey is a unique type to prevent assignment.
// Its associated value should be a string.
type contextParamKey struct {
	name string
}

// ContextWithParamValue returns a new context based on the provided parent ctx.
func ContextWithParamValue(ctx context.Context, key, val string) context.Context {
	return context.WithValue(ctx, contextParamKey{key}, val)
}

// ParamValueFromContext returns the first value associated with the given key.
func ParamValueFromContext(ctx context.Context, key string) string {
	v, _ := ctx.Value(contextParamKey{key}).(string)
	return v
}

// contextQueryKey is a unique type to prevent assignment.
// Its associated value should be a string.
type contextQueryKey struct {
	name string
}

// ContextWithQueryValue returns a new context based on the provided parent ctx.
func ContextWithQueryValue(ctx context.Context, key, val string) context.Context {
	return context.WithValue(ctx, contextQueryKey{key}, val)
}

// QueryValueFromContext returns the first value associated with the given key.
func QueryValueFromContext(ctx context.Context, key string) string {
	v, _ := ctx.Value(contextQueryKey{key}).(string)
	return v
}
