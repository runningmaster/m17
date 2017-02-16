package router

import (
	"net/http"
)

type Router interface {
	Add(string, string, http.HandlerFunc)
}

type serveMux struct {
	//mux
}

// NewHTTPRouter returns router as http.Handler.
func NewHTTPRouter() (http.Handler, error) {
	return nil, nil
}

func (m *serveMux) Add(method, path string, h http.HandlerFunc) {
	//
}

func (m *serveMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//m.mux.ServeHTTP(w, r)
}
