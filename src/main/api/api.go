package api

import (
	"context"
	"net/http"
	"strings"

	m "internal/middleware"
	"internal/router"

	"github.com/garyburd/redigo/redis"
	"github.com/nats-io/nuid"
	//"github.com/rogpeppe/fastuuid"
)

type logger interface {
	Printf(string, ...interface{})
}

type rediser interface {
	Get() redis.Conn
}

// Handler returns http.Handler based on given router.
func Handler(ctx context.Context, l logger, r router.Router, rdb rediser) (http.Handler, error) {
	api := prepareAPI(l, rdb)

	p := &m.Pipe{}
	p.AfterJoin(m.Fail, m.Tail(l))
	err404 := p.Join(m.ErrCode(http.StatusNotFound))
	err405 := p.Join(m.ErrCode(http.StatusMethodNotAllowed))

	return prepareRouter(r, api, err404, err405)
}

func prepareAPI(l logger, rdb rediser) map[string]http.Handler {
	p := &m.Pipe{}
	p.BeforeJoin(m.Head(uuid), m.Auth(auth), m.Gzip, m.Body)
	p.AfterJoin(m.JSON, m.Resp, m.Fail, m.Tail(l))

	return map[string]http.Handler{
		"GET /:foo/bar":   p.Join(m.Exec(test)),
		"GET /test/:foo":  p.Join(m.Exec(test)),
		"GET /redis/ping": p.Join(m.Exec(ping(rdb))),

		// => Debug mode only, when pref.Debug == true
		"GET /debug/vars":               p.Join(m.Exec(m.Stdh)), // expvar
		"GET /debug/pprof/":             p.Join(m.Exec(m.Stdh)), // net/http/pprof
		"GET /debug/pprof/cmdline":      p.Join(m.Exec(m.Stdh)), // net/http/pprof
		"GET /debug/pprof/profile":      p.Join(m.Exec(m.Stdh)), // net/http/pprof
		"GET /debug/pprof/symbol":       p.Join(m.Exec(m.Stdh)), // net/http/pprof
		"GET /debug/pprof/trace":        p.Join(m.Exec(m.Stdh)), // net/http/pprof
		"GET /debug/pprof/goroutine":    p.Join(m.Exec(m.Stdh)), // runtime/pprof
		"GET /debug/pprof/threadcreate": p.Join(m.Exec(m.Stdh)), // runtime/pprof
		"GET /debug/pprof/heap":         p.Join(m.Exec(m.Stdh)), // runtime/pprof
		"GET /debug/pprof/block":        p.Join(m.Exec(m.Stdh)), // runtime/pprof

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

func uuid() string {
	return nuid.Next()
}

func auth(r *http.Request) (string, int, error) {
	return "anonymous", http.StatusOK, nil
}
