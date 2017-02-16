package router

import (
	"context"
	"net/http"
)

type Router interface {
	http.Handler
	Add(string, string, http.HandlerFunc)
	Set404(http.HandlerFunc)
	Set405(http.HandlerFunc)
}

type serveMux struct {
	//mux
}

// NewHTTPRouter returns router as http.Handler.
func NewHTTPRouter() (Router, error) {
	return nil, nil
}

func ContextParam(ctx context.Context, key string) string {
	return ""
}

func (m *serveMux) Add(method, path string, h http.HandlerFunc) {
	//
}

func (m *serveMux) Set404(h http.HandlerFunc) {
	//
}

func (m *serveMux) Set405(h http.HandlerFunc) {
	//
}

func (m *serveMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//m.mux.ServeHTTP(w, r)
}
