package router

import (
	"context"
	"fmt"
	"net/http"

	"strings"

	"github.com/go-zoo/bone"
)

// Make sure the Router conforms with the HTTPRouter interface
var _ HTTPRouter = newMuxBone(context.Background())

type muxBone struct {
	ctx    context.Context
	mux    *bone.Mux
	set405 http.Handler // Bone does not support MethodNotAllowed handler
}

func newMuxBone(ctx context.Context) HTTPRouter {
	if ctx == nil {
		panic("nil context")
	}
	return &muxBone{
		ctx,
		bone.New(),
		nil,
	}
}

func (m *muxBone) Add(method, path string, h http.Handler) {
	m.mux.Register(method, path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		p := bone.GetAllValues(r)
		fmt.Println(p)
		for k := range p {
			ctx = WithParamValue(ctx, k, bone.GetValue(r, k))
		}

		for k := range r.URL.Query() {
			ctx = WithQueryValue(ctx, k, r.URL.Query().Get(k))
		}

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	}))
}

func (m *muxBone) Set404(h http.Handler) {
	m.mux.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metod := r.Method
		var allowedMethods []string
		for _, v := range []string{"GET", "POST", "PUT", "DELETE", "HEAD", "PATCH", "OPTIONS"} {
			if v == metod {
				continue
			}
			r.Method = v
			if route := m.mux.GetRequestRoute(r); len(route) > 0 && route[0] == '/' {
				allowedMethods = append(allowedMethods, v)
			}
		}
		r.Method = metod

		if len(allowedMethods) > 0 {
			w.Header().Add("Allow", strings.Join(allowedMethods, ","))
			m.set405.ServeHTTP(w, r)
			return
		}

		h.ServeHTTP(w, r)
	}))

}

func (m *muxBone) Set405(h http.Handler) {
	m.set405 = h
}

func (m *muxBone) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(m.ctx)
	m.mux.ServeHTTP(w, r)
}
