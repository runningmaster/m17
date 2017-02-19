package router

import (
	"context"
	"fmt"
	"net/http"
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
	Add(string, string, http.Handler)

	// Set404 sets custom NotFound handler.
	Set404(http.Handler)

	// Set405 sets custom MethodNotAllowed handler.
	Set405(http.Handler)
}

// New returns router as http.Handler.
func New(ctx context.Context, mux MuxKind) (HTTPRouter, error) {
	if ctx == nil {
		panic("nil context")
	}

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
