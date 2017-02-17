package router

import (
	"context"
	"net/http"

	"github.com/husobee/vestigo"
)

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
	//
}

func (m *muxVestigo) Set404(h http.Handler) {
	//
}

func (m *muxVestigo) Set405(h http.Handler) {
	//
}

func (m *muxVestigo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//m.mux.ServeHTTP(w, r)
}
