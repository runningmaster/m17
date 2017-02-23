package api

import (
	"context"
	"net/http"
	"strings"

	m "main/mdware"
	"main/router"

	"github.com/garyburd/redigo/redis"
)

var (
	err404 = m.Pipe(m.Err4xx(http.StatusNotFound), m.Tail)
	err405 = m.Pipe(m.Err4xx(http.StatusMethodNotAllowed), m.Tail)
)

var api = map[string]http.Handler{
	"GET /:foo/bar":   m.Pipe(m.Head(nil), m.Wrap(test), m.Tail),
	"GET /test/:foo":  m.Pipe(m.Head(nil), m.Wrap(test), m.Tail),
	"GET /redis/ping": m.Pipe(m.Head(nil), m.Wrap(ping), m.Tail),
}

type logger interface {
	Printf(string, ...interface{})
}

// Handler returns http.Handler based on given router.
func Handler(ctx context.Context, l logger, r router.Router, p *redis.Pool) (http.Handler, error) {

	return prepareRouter(r, api, err404, err405)
}

func prepareRouter(r router.Router, api map[string]http.Handler, err404, err405 http.Handler) (router.Router, error) {
	var s []string
	var err error
	for k, v := range api {
		s = strings.Split(k, " ")
		if len(s) != 2 {
			panic("invalid pair method-path")
		}
		err = r.Add(s[0], s[1], v)
		if err != nil {
			return nil, err
		}
	}

	err = r.Set404(err404)
	if err != nil {
		return nil, err
	}

	err = r.Set405(err405)
	if err != nil {
		return nil, err
	}
	return r, nil
}
