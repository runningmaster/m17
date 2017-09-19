package api

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"internal/logger"

	"github.com/garyburd/redigo/redis"
)

var apiFunc = map[string]func(h *redisHelper) (interface{}, error){
	"get-maker":      fooBar,
	"get-maker-sync": fooBar,
	"set-maker":      fooBar,
	"del-maker":      fooBar,

	"get-class":      fooBar,
	"get-class-sync": fooBar,
	"set-class":      fooBar,
	"del-class":      fooBar,

	"get-drug":      fooBar,
	"get-drug-sync": fooBar,
	"set-drug":      fooBar,
	"set-drug-sale": fooBar,
	"del-drug":      fooBar,

	"get-info":      fooBar,
	"get-info-sync": fooBar,
	"set-info":      fooBar,
	"set-info-sale": fooBar,
	"del-info":      fooBar,
}

//
type rediser interface {
	Get() redis.Conn
}

type redisHelper struct {
	ctx  context.Context
	rdb  rediser
	log  logger.Logger
	r    *http.Request
	w    http.ResponseWriter
	meta []byte
	data []byte
	fn   string
}

func (h *redisHelper) getConn() redis.Conn {
	return h.rdb.Get()
}

func (h *redisHelper) delConn(c io.Closer) {
	_ = c.Close
}

func (h *redisHelper) ping() (interface{}, error) {
	c := h.getConn()
	defer h.delConn(c)

	return redis.Bytes(c.Do("PING"))
}

func (h *redisHelper) exec() (interface{}, error) {
	if f, ok := apiFunc[h.fn]; ok {
		return f(h)
	}

	return nil, fmt.Errorf("unknown func %q", h.fn)
}

func fooBar(h *redisHelper) (interface{}, error) {
	return fmt.Sprintf("foo bar: %s", h.fn), nil
}
