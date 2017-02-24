package api

import (
	"context"
	"net/http"
	"strings"

	"internal/router"

	"github.com/garyburd/redigo/redis"
)

type logger interface {
	Printf(string, ...interface{})
}

// Handler returns http.Handler based on given router.
func Handler(ctx context.Context, l logger, r router.Router, p *redis.Pool) (http.Handler, error) {
	api := prepareAPI(l, p)
	err404 := pipe(err4xx(http.StatusNotFound), tail(l))
	err405 := pipe(err4xx(http.StatusMethodNotAllowed), tail(l))
	return prepareRouter(r, api, err404, err405)
}

func prepareAPI(l logger, p *redis.Pool) map[string]http.Handler {
	return map[string]http.Handler{
		"GET /:foo/bar":   pipe(head, auth, gzip, read, body(test), tail(l)),
		"GET /test/:foo":  pipe(head, auth, gzip, read, body(test), tail(l)),
		"GET /redis/ping": pipe(head, auth, gzip, read, body(ping(p)), tail(l)),
	}
}

func prepareRouter(r router.Router, api map[string]http.Handler, err404, err405 http.Handler) (router.Router, error) {
	var s []string
	var err error
	for k, v := range api {
		s = strings.Split(k, " ")
		if len(s) != 2 {
			panic("api: invalid pair method-path")
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
