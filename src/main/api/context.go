package api

import (
	"context"
	"time"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "api context value " + k.name }

// errorContextKey is a context key. The associated value will be of type error.
var errorContextKey = &contextKey{"error"}

func contextWithError(ctx context.Context, v error) context.Context {
	return context.WithValue(ctx, errorContextKey, v)
}

func errorFromContext(ctx context.Context) error {
	v, _ := ctx.Value(errorContextKey).(error)
	return v
}

// uuidContextKey is a context key. The associated value will be of type string.
var uuidContextKey = &contextKey{"UUID"}

func contextWithUUID(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, uuidContextKey, v)
}

func uuidFromContext(ctx context.Context) string {
	v, _ := ctx.Value(uuidContextKey).(string)
	return v
}

// timeContextKey is a context key. The associated value will be of type time.Time.
var timeContextKey = &contextKey{"time"}

func contextWithTime(ctx context.Context, v time.Time) context.Context {
	return context.WithValue(ctx, timeContextKey, v)
}

func timeFromContext(ctx context.Context) time.Time {
	v, _ := ctx.Value(timeContextKey).(time.Time)
	return v
}

// dataContextKey is a context key. The associated value will be of type []byte.
var dataContextKey = &contextKey{"data"}

func contextWithData(ctx context.Context, v []byte) context.Context {
	return context.WithValue(ctx, dataContextKey, v)
}

func dataFromContext(ctx context.Context) []byte {
	v, _ := ctx.Value(dataContextKey).([]byte)
	return v
}

// hostContextKey is a context key. The associated value will be of type string.
var hostContextKey = &contextKey{"host"}

func contextWithHost(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, hostContextKey, v)
}

func hostFromContext(ctx context.Context) string {
	v, _ := ctx.Value(hostContextKey).(string)
	return v
}

// userContextKey is a context key. The associated value will be of type string.
var userContextKey = &contextKey{"user"}

func contextWithUser(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, userContextKey, v)
}

func userFromContext(ctx context.Context) string {
	v, _ := ctx.Value(userContextKey).(string)
	return v
}

// authContextKey is a context key. The associated value will be of type string.
var authContextKey = &contextKey{"auth"}

func contextWithAuth(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, authContextKey, v)
}

func authFromContext(ctx context.Context) string {
	v, _ := ctx.Value(authContextKey).(string)
	return v
}

// clenContextKey is a context key. The associated value will be of type int64.
var clenContextKey = &contextKey{"content-length"}

func contextWithClen(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, clenContextKey, v)
}

func clenFromContext(ctx context.Context) int64 {
	v, _ := ctx.Value(clenContextKey).(int64)
	return v
}

// sizeContextKey is a context key. The associated value will be of type int64.
var sizeContextKey = &contextKey{"size"}

func contextWithSize(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, sizeContextKey, v)
}

func sizeFromContext(ctx context.Context) int64 {
	v, _ := ctx.Value(sizeContextKey).(int64)
	return v
}

// codeContextKey is a context key. The associated value will be of type int.
var codeContextKey = &contextKey{"code"}

func contextWithCode(ctx context.Context, v int) context.Context {
	return context.WithValue(ctx, codeContextKey, v)
}

func codeFromContext(ctx context.Context) int {
	v, _ := ctx.Value(codeContextKey).(int)
	return v
}
