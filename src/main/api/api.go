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
	err404 := p.join(errCode(http.StatusNotFound))
	err405 := p.join(errCode(http.StatusMethodNotAllowed))

	return prepareRouter(r, api, err404, err405)
}

func prepareAPI(l logger, c redisConner) map[string]http.Handler {
	p := &pipe{}
	p.head(uuid, auth, gunzip, gzip, body)
	// -> join result below will be here <-
	p.tail(mrshl, resp, errf, logg(l))

	return map[string]http.Handler{
		"GET /:foo/bar":   p.join(exec(http.HandlerFunc(test))),
		"GET /test/:foo":  p.join(exec(http.HandlerFunc(test))),
		"GET /redis/ping": p.join(exec(ping(c))),

		// => Debug mode only, when pref.Debug == true
		"GET /debug/vars":               p.join(exec(http.HandlerFunc(stdh))), // expvar
		"GET /debug/pprof/":             p.join(exec(http.HandlerFunc(stdh))), // net/http/pprof
		"GET /debug/pprof/cmdline":      p.join(exec(http.HandlerFunc(stdh))), // net/http/pprof
		"GET /debug/pprof/profile":      p.join(exec(http.HandlerFunc(stdh))), // net/http/pprof
		"GET /debug/pprof/symbol":       p.join(exec(http.HandlerFunc(stdh))), // net/http/pprof
		"GET /debug/pprof/trace":        p.join(exec(http.HandlerFunc(stdh))), // net/http/pprof
		"GET /debug/pprof/goroutine":    p.join(exec(http.HandlerFunc(stdh))), // runtime/pprof
		"GET /debug/pprof/threadcreate": p.join(exec(http.HandlerFunc(stdh))), // runtime/pprof
		"GET /debug/pprof/heap":         p.join(exec(http.HandlerFunc(stdh))), // runtime/pprof
		"GET /debug/pprof/block":        p.join(exec(http.HandlerFunc(stdh))), // runtime/pprof

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
