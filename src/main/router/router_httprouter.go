package router

import (
	"context"
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

func (m *muxHTTPRouter) Add(method, path string, h http.Handler) {
	m.mux.Handle(method, path,
		func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			ctx := r.Context()
			for i := range p {
				ctx = context.WithValue(ctx, ParamKey{p[i].Key}, p[i].Value)
			}
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
}

func (m *muxHTTPRouter) Set404(h http.Handler) {
	m.mux.NotFound = h
}

func (m *muxHTTPRouter) Set405(h http.Handler) {
	m.mux.MethodNotAllowed = h
}

func (m *muxHTTPRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r = r.WithContext(ctx)
	m.mux.ServeHTTP(w, r)
}
