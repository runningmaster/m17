package api

import (
	"net/http"

	m "main/mdware"
)

var (
	err404 = m.Pipe(m.Err4xx(http.StatusNotFound), m.Tail)
	err405 = m.Pipe(m.Err4xx(http.StatusMethodNotAllowed), m.Tail)
)

var multiplexer = map[string]http.Handler{
	"GET /:foo/bar":   m.Pipe(m.Head(nil), m.Wrap(test), m.Tail),
	"GET /test/:foo":  m.Pipe(m.Head(nil), m.Wrap(test), m.Tail),
	"GET /redis/ping": m.Pipe(m.Head(nil), m.Wrap(ping), m.Tail),
}
