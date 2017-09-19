package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"internal/csvutil"
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

func (h *redisHelper) uploadSuggestion(_, data []byte) (interface{}, error) {
	c := h.getConn()
	defer h.delConn(c)

	c.Do("FT.DROP", "spc")
	c.Do("FT.DROP", "atc")
	c.Do("FT.DROP", "inn")
	c.Do("FT.DROP", "org")

	c.Do("FT.CREATE", "spc", "NOOFFSETS", "NOFIELDS", "NOSCOREIDX", "NOFREQS", "SCHEMA", "name", "TEXT", "SORTABLE")
	c.Do("FT.CREATE", "atc", "NOOFFSETS", "NOFIELDS", "NOSCOREIDX", "NOFREQS", "SCHEMA", "name", "TEXT", "SORTABLE")
	c.Do("FT.CREATE", "inn", "NOOFFSETS", "NOFIELDS", "NOSCOREIDX", "NOFREQS", "SCHEMA", "name", "TEXT", "SORTABLE")
	c.Do("FT.CREATE", "org", "NOOFFSETS", "NOFIELDS", "NOSCOREIDX", "NOFREQS", "SCHEMA", "name", "TEXT", "SORTABLE")

	csv := csvutil.NewRecordChan(bytes.NewReader(data), ',', false, 1)
	var err error
	var n int
	for v := range csv {
		if v.Error != nil {
			//continue
			return nil, v.Error
		}
		if len(v.Record) < 4 {
			return nil, fmt.Errorf("invalid csv: got %d, want %d", len(v.Record), 4)
		}

		switch v.Record[0] {
		case "atc":
			//c.Do("FT.ADD", "atc", v.Record[1], "1", "FIELDS", "name", strings.ToLower(v.Record[2]))
		case "info":
			c.Do("FT.ADD", "spc", v.Record[1], "1", "FIELDS", "name", strings.ToLower(v.Record[2]))
			fmt.Println(strings.ToLower(v.Record[2]))
		case "inn":
			c.Do("FT.ADD", "inn", v.Record[1], "1", "FIELDS", "name", strings.ToLower(v.Record[2]))
		case "org":
			c.Do("FT.ADD", "org", v.Record[1], "1", "FIELDS", "name", strings.ToLower(v.Record[2]))
		}
		if err != nil {
			return nil, err
		}
		n++
	}

	return []byte(fmt.Sprintf("OK: %d", n)), nil
}

func (h *redisHelper) selectSuggestion(_, data []byte) (interface{}, error) {
	v := struct {
		Name string `json:"name"`
	}{}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}

	n := len([]rune(v.Name))
	if n <= 2 {
		return nil, fmt.Errorf("too few characters: %d", n)
	}

	if n > 128 {
		return nil, fmt.Errorf("too many characters: %d", n)
	}

	c := h.getConn()
	defer h.delConn(c)

	res, err := redis.Values(c.Do("FT.SEARCH", "spc", v.Name, "NOCONTENT", "SORTBY", "name"))
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%v: %v", res[0], len(res)), nil
}
