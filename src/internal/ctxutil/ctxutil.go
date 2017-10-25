package ctxutil

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// key is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type key struct {
	name string
}

// String supports Stringer interface.
func (k *key) String() string {
	return "api context value " + k.name
}

// keyCode is a context key. The associated value will be of type int.
var keyCode = &key{"Code"}

// WithCode returns context.Context value with v.
func WithCode(ctx context.Context, v int) context.Context {
	return context.WithValue(ctx, keyCode, v)
}

// CodeFrom returns v from context.Context value placed in it by WithCode().
func CodeFrom(ctx context.Context) int {
	v, _ := ctx.Value(keyCode).(int)
	return v
}

// keyError is a context key. The associated value will be of type error.
var keyError = &key{"Error"}

// WithError returns context.Context with err and code.
func WithError(ctx context.Context, err error, codes ...int) context.Context {
	if err == nil {
		return ctx
	}

	code := CodeFrom(ctx)
	if code == 0 {
		code = http.StatusInternalServerError
	}
	for i := range codes {
		code = codes[i] // last is true, amen
	}

	err = fmt.Errorf("%d %s: %s", code, http.StatusText(code), err.Error())
	ctx = context.WithValue(ctx, keyError, err)
	return WithCode(ctx, code)
}

// ErrorFrom returns v from context.Context value placed in it by WithError().
func ErrorFrom(ctx context.Context) error {
	v, _ := ctx.Value(keyError).(error)
	return v
}

// keyUUID is a context key. The associated value will be of type string.
var keyUUID = &key{"UUID"}

// WithUUID returns context.Context value with v.
func WithUUID(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, keyUUID, v)
}

// UUIDFrom returns v from context.Context value placed in it by WithUUID().
func UUIDFrom(ctx context.Context) string {
	v, _ := ctx.Value(keyUUID).(string)
	return v
}

// keyTime is a context key. The associated value will be of type time.Time.
var keyTime = &key{"Time"}

// WithTime returns context.Context value with v.
func WithTime(ctx context.Context, v time.Time) context.Context {
	return context.WithValue(ctx, keyTime, v)
}

// TimeFrom returns v from context.Context value placed in it by WithTime().
func TimeFrom(ctx context.Context) time.Time {
	v, _ := ctx.Value(keyTime).(time.Time)
	return v
}

// keyHost a context key. The associated value will be of type string.
var keyHost = &key{"Host"}

// WithHost returns context.Context value with v.
func WithHost(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, keyHost, v)
}

// HostFrom returns v from context.Context value placed in it by WithHost().
func HostFrom(ctx context.Context) string {
	v, _ := ctx.Value(keyHost).(string)
	return v
}

// keyUser is a context key. The associated value will be of type string.
var keyUser = &key{"User"}

// WithUser returns context.Context value with v.
func WithUser(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, keyUser, v)
}

// UserFrom returns v from context.Context value placed in it by WithUser().
func UserFrom(ctx context.Context) string {
	v, _ := ctx.Value(keyUser).(string)
	return v
}

// keyAuth is a context key. The associated value will be of type string.
var keyAuth = &key{"Auth"}

// WithAuth returns context.Context value with v.
func WithAuth(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, keyAuth, v)
}

// AuthFrom returns v from context.Context value placed in it by WithAuth().
func AuthFrom(ctx context.Context) string {
	v, _ := ctx.Value(keyAuth).(string)
	return v
}

// keyCLen is a context key. The associated value will be of type int64.
var keyCLen = &key{"Content-Length"}

// WithCLen returns context.Context value with v.
func WithCLen(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, keyCLen, v)
}

// CLenFrom returns v from context.Context value placed in it by WithCLen().
func CLenFrom(ctx context.Context) int64 {
	v, _ := ctx.Value(keyCLen).(int64)
	return v
}

// keySize is a context key. The associated value will be of type int64.
var keySize = &key{"Size"}

// WithSize returns context.Context value with v.
func WithSize(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, keySize, v)
}

// SizeFrom returns v from context.Context value placed in it by WithSize().
func SizeFrom(ctx context.Context) int64 {
	v, _ := ctx.Value(keySize).(int64)
	return v
}

// keyBody is a context key. The associated value will be of type []byte.
var keyBody = &key{"Body"}

// WithBody returns context.Context value with v.
func WithBody(ctx context.Context, v []byte) context.Context {
	return context.WithValue(ctx, keyBody, v)
}

// BodyFrom returns v from context.Context value placed in it by WithBody().
func BodyFrom(ctx context.Context) []byte {
	v, _ := ctx.Value(keyBody).([]byte)
	return v
}

// keyResult is a context key. The associated value will be of type interface{}.
var keyResult = &key{"Data"}

// WithResult returns context.Context value with v.
func WithResult(ctx context.Context, v interface{}) context.Context {
	return context.WithValue(ctx, keyResult, v)
}

// ResultFrom returns v from context.Context value placed in it by WithResult().
func ResultFrom(ctx context.Context) interface{} {
	return ctx.Value(keyResult)
}
