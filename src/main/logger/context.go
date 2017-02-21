package logger

import (
	"context"
)

// contextLoggerKey is a unique type to prevent assignment.
// Its associated value should be a Logger interface.
type contextLoggerKey struct{}

// ContextWithLogger returns a new context based on the provided parent ctx.
func ContextWithLogger(ctx context.Context, v Logger) context.Context {
	return context.WithValue(ctx, contextLoggerKey{}, v)
}

// FromContext returns the Logger value associated with the given key.
func FromContext(ctx context.Context) Logger {
	if v, ok := ctx.Value(contextLoggerKey{}).(Logger); ok {
		return v
	}
	panic("Logger not found")
}
