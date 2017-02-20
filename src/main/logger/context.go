package logger

import (
	"context"
)

// contextLoggerKey is a unique type to prevent assignment.
// Its associated value should be a Logger interface.
type contextLoggerKey struct{}

// WithLogger returns a new context based on the provided parent ctx.
func WithLogger(ctx context.Context, v Logger) context.Context {
	ctxKey := contextLoggerKey{}
	return context.WithValue(ctx, ctxKey, v)
}

// ContextLogger returns the Logger value associated with the given key.
func ContextLogger(ctx context.Context) Logger {
	ctxKey := contextLoggerKey{}
	if v, ok := ctx.Value(ctxKey).(Logger); ok {
		return v
	}
	panic("Logger not found")
}
