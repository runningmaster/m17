package router

import (
	"context"
	"net/http"

	"github.com/go-zoo/bone"
)

type muxBone struct {
	ctx context.Context
	mux *bone.Mux
}

func newMuxBone(ctx context.Context) HTTPRouter {
	return &muxBone{
		ctx,
		bone.New(),
	}
}

func (m *muxBone) Add(method, path string, h http.Handler) {
	//
}

func (m *muxBone) Set404(h http.Handler) {
	//
}

func (m *muxBone) Set405(h http.Handler) {
	//
}

func (m *muxBone) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//m.mux.ServeHTTP(w, r)
}
