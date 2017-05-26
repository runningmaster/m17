package api

import (
	"context"
	"io"

	"internal/logger"

	"github.com/garyburd/redigo/redis"
)

//
type rediser interface {
	Get() redis.Conn
}

type redisHelper struct {
	ctx context.Context
	rdb rediser
	log logger.Logger
}

func newRedisHelper(ctx context.Context, rdb rediser, log logger.Logger) *redisHelper {
	return &redisHelper{
		ctx: ctx,
		rdb: rdb,
		log: log,
	}
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

func (h *redisHelper) uploadData(meta, data []byte) (interface{}, error) {
	return nil, nil
}
