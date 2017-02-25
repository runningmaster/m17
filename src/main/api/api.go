package api

import (
	"context"
	"net/http"
	"strings"

	"internal/router"
)

type logger interface {
	Printf(string, ...interface{})
}

// Handler returns http.Handler based on given router.
func Handler(ctx context.Context, l logger, r router.Router, c redisConner) (http.Handler, error) {
	api := prepareAPI(l, c)

	p := &pipe{}
	p.tail(logg(l))
	err404 := p.join(err4xx(http.StatusNotFound))
	err405 := p.join(err4xx(http.StatusMethodNotAllowed))

	return prepareRouter(r, api, err404, err405)
}

func prepareAPI(l logger, c redisConner) map[string]http.Handler {
	p := &pipe{}
	p.head(uuid, auth, gzip, read)
	p.tail(logg(l))

	return map[string]http.Handler{
		"GET /:foo/bar":   p.join(skipIfError(test)),
		"GET /test/:foo":  p.join(skipIfError(test)),
		"GET /redis/ping": p.join(skipIfError(ping(c))),
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
