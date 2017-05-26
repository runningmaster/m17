package api

import (
	"net/http"
	"strings"

	"internal/logger"
	"internal/mdware"
	"internal/router"

	"github.com/garyburd/redigo/redis"
	"github.com/nats-io/nuid"
	//"github.com/rogpeppe/fastuuid"
)

type rediser interface {
	Get() redis.Conn
}

type handler struct {
	api    map[string]http.Handler
	rdb    rediser
	log    logger.Logger
	err404 http.Handler
	err405 http.Handler
}

func (h *handler) prepareAPI() *handler {
	p := &mdware.Pipe{}
	p.BeforeJoin(
		mdware.Head(uuid),
		mdware.Auth(auth),
		mdware.Gzip,
		mdware.Body,
	)
	p.AfterJoin(
		mdware.Resp,
		mdware.Fail,
		mdware.Tail(h.log),
	)

	h.api = map[string]http.Handler{
		"GET /:foo/bar":   p.Join(mdware.Exec(test)),
		"GET /test/:foo":  p.Join(mdware.Exec(test)),
		"GET /redis/ping": p.Join(mdware.Exec(ping(h.rdb))),

		// => Debug mode only, when pref.Debug == true
		"GET /debug/vars":               p.Join(mdware.Exec(mdware.Stdh)), // expvar
		"GET /debug/pprof/":             p.Join(mdware.Exec(mdware.Stdh)), // net/http/pprof
		"GET /debug/pprof/cmdline":      p.Join(mdware.Exec(mdware.Stdh)), // net/http/pprof
		"GET /debug/pprof/profile":      p.Join(mdware.Exec(mdware.Stdh)), // net/http/pprof
		"GET /debug/pprof/symbol":       p.Join(mdware.Exec(mdware.Stdh)), // net/http/pprof
		"GET /debug/pprof/trace":        p.Join(mdware.Exec(mdware.Stdh)), // net/http/pprof
		"GET /debug/pprof/goroutine":    p.Join(mdware.Exec(mdware.Stdh)), // runtime/pprof
		"GET /debug/pprof/threadcreate": p.Join(mdware.Exec(mdware.Stdh)), // runtime/pprof
		"GET /debug/pprof/heap":         p.Join(mdware.Exec(mdware.Stdh)), // runtime/pprof
		"GET /debug/pprof/block":        p.Join(mdware.Exec(mdware.Stdh)), // runtime/pprof
	}

	h.err404 = mdware.Join(
		mdware.Head(uuid),
		mdware.Errc(http.StatusNotFound),
		mdware.Fail,
		mdware.Tail(h.log),
	)

	h.err405 = mdware.Join(
		mdware.Head(uuid),
		mdware.Errc(http.StatusMethodNotAllowed),
		mdware.Fail,
		mdware.Tail(h.log),
	)

	return h
}

func (h *handler) withRouter(r router.Router) (router.Router, error) {
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

// NewWithRouter returns http.Handler based on given router.
func NewWithRouter(r router.Router, options ...func(*handler) error) (router.Router, error) {
	if r == nil {
		panic("nil router")
	}

	h := &handler{
		log: logger.NewDefault(),
		rdb: &redis.Pool{},
	}

	var err error
	for i := range options {
		err = options[i](h)
		if err != nil {
			return nil, err
		}
	}

	return h.prepareAPI().withRouter(r)
}

// Logger is option for passing logger interface.
func Logger(l logger.Logger) func(*handler) error {
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
