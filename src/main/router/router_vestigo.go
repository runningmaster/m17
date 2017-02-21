package router

import (
	"context"
	"fmt"
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
		ctx: ctx,
		mux: vestigo.NewRouter(),
	}
}

func (m *muxVestigo) Add(method, path string, h http.Handler) error {
	err := validateAddParams(method, path, h)
	if err != nil {
		return err
	}

	m.mux.Add(method, path, func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		p := vestigo.TrimmedParamNames(r)
		for i := range p {
			ctx = ContextWithParamValue(ctx, p[i], vestigo.Param(r, p[i]))
		}

		for k := range r.URL.Query() {
			ctx = ContextWithQueryValue(ctx, k, r.URL.Query().Get(k))
		}

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})

	return nil
}

func (m *muxVestigo) Set404(h http.Handler) error {
	if h == nil {
		return fmt.Errorf("%v handler", h)
	}

	vestigo.CustomNotFoundHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
	return nil
}

func (m *muxVestigo) Set405(h http.Handler) error {
	if h == nil {
		return fmt.Errorf("%v handler", h)
	}

	vestigo.CustomMethodNotAllowedHandlerFunc(
		func(allowedMethods string) func(w http.ResponseWriter, r *http.Request) {
			return func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Allow", allowedMethods)
				h.ServeHTTP(w, r)
			}
		})
	return nil
}

func (m *muxVestigo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(m.ctx)
	m.mux.ServeHTTP(w, r)
}
