package router

import (
	"context"
	"fmt"
	"net/http"
)

// Router is interface for implementing in API pakage.
type Router interface {
	http.Handler

	// Add adds a method/handler combination to the router.
	Add(string, string, http.Handler) error

	// Set404 sets custom NotFound handler.
	Set404(http.Handler) error

	// Set405 sets custom MethodNotAllowed handler.
	Set405(http.Handler) error
}

// NewMuxBone returns interface Router based on "bone" multiplexer.
func NewMuxBone(ctx context.Context) Router {
	return newMuxBone(ctx)
}

// NewMuxHTTPRouter returns interface Router based on "httprouter" multiplexer.
func NewMuxHTTPRouter(ctx context.Context) Router {
	return newMuxHTTPRouter(ctx)
}

// NewMuxVestigo returns interface Router based on "vestigo" multiplexer.
func NewMuxVestigo(ctx context.Context) Router {
	return newMuxVestigo(ctx)
}

var (
	methodMap = map[string]struct{}{
		"GET":     struct{}{},
		"POST":    struct{}{},
		"PUT":     struct{}{},
		"DELETE":  struct{}{},
		"HEAD":    struct{}{},
		"PATCH":   struct{}{},
		"OPTIONS": struct{}{},
	}
	formatErrInvalidValue = "%%v is invalid value"
)

func validateAddParams(method string, path string, h http.Handler) error {
	if _, ok := methodMap[method]; !ok {
		return fmt.Errorf(formatErrInvalidValue, method)
	}
	if path[0] != '/' {
		return fmt.Errorf(formatErrInvalidValue, path)
	}
	if h == nil {
		return fmt.Errorf(formatErrInvalidValue, h)
	}
	return nil
}

// ctxParam is a unique type to prevent assignment.
// Its associated value should be a string.
type ctxParam struct {
	name string
}

// ctxWithParamValue returns a new context based on the provided parent ctx.
func ctxWithParamValue(ctx context.Context, key, val string) context.Context {
	return context.WithValue(ctx, ctxParam{key}, val)
}

// ParamValueFrom returns the first value associated with the given key.
func ParamValueFrom(ctx context.Context, key string) string {
	v, _ := ctx.Value(ctxParam{key}).(string)
	return v
}

// ctxQuery is a unique type to prevent assignment.
// Its associated value should be a string.
type ctxQuery struct {
	name string
}

// ctxWithQueryValue returns a new context based on the provided parent ctx.
func ctxWithQueryValue(ctx context.Context, key, val string) context.Context {
	return context.WithValue(ctx, ctxQuery{key}, val)
}

// QueryValueFrom returns the first value associated with the given key.
func QueryValueFrom(ctx context.Context, key string) string {
	v, _ := ctx.Value(ctxQuery{key}).(string)
	return v
}
