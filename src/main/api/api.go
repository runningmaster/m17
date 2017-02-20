package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"main/client"
	"main/logger"
	m "main/mdware"
	"main/router"

	"github.com/garyburd/redigo/redis"
)

type routeTable map[string]http.Handler

var table = routeTable{
	"GET /:foo/bar":   m.Pipe(m.Head(nil), m.Wrap(test), m.Tail),
	"GET /test/:foo":  m.Pipe(m.Head(nil), m.Wrap(test), m.Tail),
	"GET /redis/ping": m.Pipe(m.Head(nil), m.Wrap(ping), m.Tail),
}

// NewHandler returns HTTP API handler.
func NewHandler(ctx context.Context) (http.Handler, error) {

	// make redis pool here
	return makeHTTPRouter(ctx, table)
}

func makeHTTPRouter(ctx context.Context, t routeTable) (router.HTTPRouter, error) {
	r, err := router.New(ctx, router.MuxBone)
	if err != nil {
		return nil, err
	}

	var s []string
	for k, v := range t {
		s = strings.Split(k, " ")
		if len(s) != 2 {
			panic("invalid pair method-path")
		}
		err = r.Add(s[0], s[1], v)
		if err != nil {
			return nil, err
		}
	}

	err = r.Set404(m.Pipe(m.Err4xx(http.StatusNotFound), m.Tail))
	if err != nil {
		return nil, err
	}

	err = r.Set405(m.Pipe(m.Err4xx(http.StatusMethodNotAllowed), m.Tail))
	if err != nil {
		return nil, err
	}
	return r, nil
}

func test(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Fprintf(w, "Hello, World! From test handler!\n")
	fmt.Fprintf(w, "Param foo: %s\n", router.ContextParamValue(ctx, "foo"))
	fmt.Fprintf(w, "Query foo: %s\n", router.ContextQueryValue(ctx, "foo"))
	v, _ := ctx.Value("foo").(string)
	fmt.Fprintf(w, "Value foo: %s\n", v)
	*r = *r.WithContext(ctx)
}

func ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.ContextLogger(ctx)

	res, err := redisPing(ctx)
	if err != nil {
		log.Printf("redis error: %v\n", err)
		fmt.Fprintf(w, "redis error: %v\n", err)

	}
	log.Printf("redis result: %v\n", string(res))
	fmt.Fprintf(w, "redis result: %v\n", string(res))
	*r = *r.WithContext(ctx)
}

func redisPing(ctx context.Context) ([]byte, error) {
	cli := client.ContextRedisConn(ctx)
	defer cli.Close()

	return redis.Bytes(cli.Do("PING"))
}
