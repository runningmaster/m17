package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Make sure the Router conforms with the HTTPRouter interface
var _ HTTPRouter = newMuxHTTPRouter(context.Background())

type muxHTTPRouter struct {
	ctx context.Context
	mux *httprouter.Router
}

func newMuxHTTPRouter(ctx context.Context) HTTPRouter {
	return &muxHTTPRouter{
		ctx,
		httprouter.New(),
	}
}

func (m *muxHTTPRouter) Add(method, path string, h http.Handler) error {
	err := validateAddParams(method, path, h)
	if err != nil {
		return err
	}

	m.mux.Handle(method, path,
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			ctx := r.Context()

			for i := range p {
				ctx = ContextWithParamValue(ctx, p[i].Key, p[i].Value)
			}

			for k := range r.URL.Query() {
				ctx = ContextWithQueryValue(ctx, k, r.URL.Query().Get(k))
			}

			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})

	return nil
}

func (m *muxHTTPRouter) Set404(h http.Handler) error {
	if h == nil {
		return fmt.Errorf("%v handler", h)
	}

	m.mux.NotFound = h
	return nil
}

func (m *muxHTTPRouter) Set405(h http.Handler) error {
	if h == nil {
		return fmt.Errorf("%v handler", h)
	}

	m.mux.MethodNotAllowed = h
	return nil
}

func (m *muxHTTPRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(m.ctx)
	m.mux.ServeHTTP(w, r)
}
