package api

import (
	"log"
	"net/http"
	"os"
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

type handler struct {
	api    map[string]http.Handler
	err404 http.Handler
	err405 http.Handler
	rdb    rediser
	log    logger
}

func (h *handler) prepareAPI() *handler {
	p := &m.Pipe{}
	p.BeforeJoin(
		m.Head(uuid),
		m.Auth(auth),
		m.Gzip,
		m.Body,
	)
	p.AfterJoin(
		m.Resp,
		m.Fail,
		m.Tail(h.log),
	)

	h.api = map[string]http.Handler{
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

	h.err404 = m.Join(m.Head(uuid), m.Errc(http.StatusNotFound), m.Fail, m.Tail(h.log))
	h.err405 = m.Join(m.Head(uuid), m.Errc(http.StatusMethodNotAllowed), m.Fail, m.Tail(h.log))

	return h
}

func (h *handler) prepareRouter(r router.Router) (router.Router, error) {
	var s []string
	var err error
	for k, v := range h.api {
		s = strings.Split(k, " ")
		if len(s) != 2 {
			panic("api: invalid pair method-path")
		}
		err = r.Add(s[0], s[1], v)
		if err != nil {
			return nil, err
		}
	}

	err = r.Set404(h.err404)
	if err != nil {
		return nil, err
	}

	err = r.Set405(h.err405)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// MustWithRouter returns http.Handler based on given router.
func MustWithRouter(r router.Router, options ...func(*handler) error) (router.Router, error) {
	if r == nil {
		panic("nil router")
	}

	h := &handler{
		log: log.New(os.Stderr, "", log.LstdFlags),
		rdb: &redis.Pool{},
	}

	var err error
	for i := range options {
		err = options[i](h)
		if err != nil {
			return nil, err
		}
	}

	return h.prepareAPI().prepareRouter(r)
}

// Logger is option for passing logger interface.
func Logger(l logger) func(*handler) error {
	return func(h *handler) error {
		h.log = l
		return nil
	}
}

// Redis is interface for Redis Pool Connections.
func Redis(r rediser) func(*handler) error {
	return func(h *handler) error {
		h.rdb = r
		return nil
	}
}

func uuid() string {
	return nuid.Next()
}

func auth(r *http.Request) (string, int, error) {
	return "anonymous", http.StatusOK, nil
}
