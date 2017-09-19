package api

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"internal/logger"

	"github.com/garyburd/redigo/redis"
)

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
	if h.fn == "" {
		return nil, fmt.Errorf("func is empty %q", h.fn)
	}

	return nil, fmt.Errorf("unknown func %q", h.fn)
}

/*
	"get-maker"
	"get-maker-sync"
	"set-maker"
	"del-maker"

	"get-class"
	"get-class-sync"
	"set-class"
	"del-class"

	"get-drug"
	"get-drug-sync"
	"set-drug"
	"set-drug-sale"
	"del-drug"

	"get-info"
	"get-info-sync"
	"set-info"
	"set-info-sale"
	"del-info"
*/
