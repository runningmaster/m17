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

type Handler struct {
	api map[string]http.Handler
	rdb rediser
	rtr router.Router
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.rtr.ServeHTTP(w, r)
}

// Redis is interface for Redis Pool Connections
func Redis(r rediser) func(*Handler) error {
	return func(h *Handler) error {
		h.rdb = r
		return nil
	}
}

// Router is interface for HTTP Router inplementation
func Router(r router.Router) func(*Handler) error {
	return func(h *Handler) error {
		h.rtr = r
		return nil
	}
}

// NewHandler returns http.Handler based on given router.
func NewHandler(ctx context.Context, l logger, options ...func(*Handler) error) (http.Handler, error) {
	h := &Handler{}

	for i := range options {
		err = options[i](h)
		if err != nil {
			return err
		}
	}

	p := &m.Pipe{}
	p.BeforeJoin(
		m.Head(uuid),
		m.Auth(auth),
		m.Gzip,
		m.Body,
	)
	p.AfterJoin(
		m.JSON,
		m.Resp,
		m.Fail,
		m.Tail(l),
	)

	api := map[string]http.Handler{
		"GET /:foo/bar":   p.Join(m.Exec(test)),
		"GET /test/:foo":  p.Join(m.Exec(test)),
		"GET /redis/ping": p.Join(m.Exec(ping(h.rdb))),

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

	p := &m.Pipe{}
	p.AfterJoin(m.Fail, m.Tail(l))
	err404 := p.Join(m.ErrCode(http.StatusNotFound))
	err405 := p.Join(m.ErrCode(http.StatusMethodNotAllowed))

	return prepareRouter(r, api, err404, err405)
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
