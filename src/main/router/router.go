package router

import (
	"context"
	"fmt"
	"net/http"
)

type kindMux int

const (
	kindBone kindMux = iota
	kindHTTPRouter
	kindVestigo
)

var methodMap = map[string]struct{}{
	"GET":     struct{}{},
	"POST":    struct{}{},
	"PUT":     struct{}{},
	"DELETE":  struct{}{},
	"HEAD":    struct{}{},
	"PATCH":   struct{}{},
	"OPTIONS": struct{}{},
}

// String is satisfy fmt.Stringer interface.
func (kind kindMux) String() string {
	switch kind {
	case kindBone:
		return "bone"
	case kindHTTPRouter:
		return "httprouter"
	case kindVestigo:
		return "vestigo"
	default:
		return ("unknown multipexer")
	}
}

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

// New returns interface Router.
func New(ctx context.Context, options ...func(*Option) error) (Router, error) {
	err := defaultOption.override(options...)
	if err != nil {
		return nil, err
	}

	kind := defaultOption.kind
	switch kind {
	case kindBone:
		return newMuxBone(ctx), nil
	case kindHTTPRouter:
		return newMuxHTTPRouter(ctx), nil
	case kindVestigo:
		return newMuxVestigo(ctx), nil
	default:
		return nil, fmt.Errorf("%s not implemented", kind)
	}
}

func validateAddParams(method string, path string, h http.Handler) error {
	if _, ok := methodMap[method]; !ok {
		return fmt.Errorf("%v is invalid method", method)
	}
	if path[0] != '/' {
		return fmt.Errorf("%v is invalid path", path)
	}
	if h == nil {
		return fmt.Errorf("%v is invalid handler", h)
	}

	return nil
}
