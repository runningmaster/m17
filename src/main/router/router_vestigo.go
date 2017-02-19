package router

import (
	"context"
	"net/http"

	"github.com/husobee/vestigo"
)

// Make sure the Router conforms with the HTTPRouter interface
var _ HTTPRouter = newMuxVestigo(context.Background())

type muxVestigo struct {
	ctx context.Context
	mux *vestigo.Router
}

func newMuxVestigo(ctx context.Context) HTTPRouter {
	return &muxVestigo{
		ctx,
		vestigo.NewRouter(),
	}
}

func (m *muxVestigo) Add(method, path string, h http.Handler) {
	m.mux.Add(method, path, func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		p := vestigo.TrimmedParamNames(r)
		for i := range p {
			ctx = WithParamValue(ctx, p[i], vestigo.Param(r, p[i]))
		}

		for k := range r.URL.Query() {
			ctx = WithQueryValue(ctx, k, r.URL.Query().Get(k))
		}

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func (m *muxVestigo) Set404(h http.Handler) {
	vestigo.CustomNotFoundHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func (m *muxVestigo) Set405(h http.Handler) {
	vestigo.CustomMethodNotAllowedHandlerFunc(
		func(allowedMethods string) func(w http.ResponseWriter, r *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Allow", allowedMethods)
				h.ServeHTTP(w, r)
			}
		})
}

func (m *muxVestigo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(m.ctx)
	m.mux.ServeHTTP(w, r)
}
