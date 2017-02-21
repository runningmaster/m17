package router

import (
	"context"
	"fmt"
	"net/http"

	"main/logger"
	"main/option"
)

// MuxKind is kind of multipexers for HTTP routing.
type MuxKind int

const (
	// MuxBone is
	// github.com/go-zoo/bone
	MuxBone MuxKind = iota

	// MuxHTTPRouter is
	// github.com/julienschmidt/httprouter
	MuxHTTPRouter

	// MuxVestigo is
	// github.com/husobee/vestigo
	MuxVestigo
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

var log logger.Logger

// String is satisfy fmt.Stringer interface.
func (m MuxKind) String() string {
	switch m {
	case MuxBone:
		return "bone"
	case MuxHTTPRouter:
		return "httprouter"
	case MuxVestigo:
		return "vestigo"
	default:
		return ("unknown multipexer")
	}
}

// HTTPRouter is interface for implementing in API pakage.
type HTTPRouter interface {
	http.Handler

	// Add adds a method/handler combination to the router.
	Add(string, string, http.Handler) error

	// Set404 sets custom NotFound handler.
	Set404(http.Handler) error

	// Set405 sets custom MethodNotAllowed handler.
	Set405(http.Handler) error
}

// New returns router as http.Handler.
func New(ctx context.Context, options ...option.Fn) (HTTPRouter, error) {
	opt := &optionReceiver{}
	err := opt.Receive(options...)
	if err != nil {
		return nil, err
	}

	mux := opt.muxKind
	switch mux {
	case MuxBone:
		return newMuxBone(ctx), nil
	case MuxHTTPRouter:
		return newMuxHTTPRouter(ctx), nil
	case MuxVestigo:
		return newMuxVestigo(ctx), nil
	default:
		return nil, fmt.Errorf("%s not implemented", mux)
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
