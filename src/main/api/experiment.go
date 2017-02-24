package api

import (
	"context"
	"fmt"
	"net/http"

	"main/router"

	"github.com/garyburd/redigo/redis"
)

type redisPool interface {
	Get() redis.Conn
}

func uuid() string {
	println("\tuuid")
	return "UUIDMustBeHereFIXME"
}

func xxxx(key string) bool {
	println("\tauth")
	if key != "" {
		return false
	}
	return true
}

func test(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Fprintf(w, "Hello, World! From test handler!\n")
	fmt.Fprintf(w, "Param foo: %s\n", router.ParamValueFromContext(ctx, "foo"))
	fmt.Fprintf(w, "Query foo: %s\n", router.QueryValueFromContext(ctx, "foo"))
	v, _ := ctx.Value("foo").(string)
	fmt.Fprintf(w, "Value foo: %s\n", v)
	*r = *r.WithContext(ctx)
}

func ping(p redisPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		res, err := redisPing(ctx, p.Get())
		if err != nil {
			// withFail
			fmt.Fprintf(w, "redis error: %v\n", err)

		}
		fmt.Fprintf(w, "redis result: %v\n", string(res))
		*r = *r.WithContext(ctx)
	}
}

func redisPing(_ context.Context, c redis.Conn) ([]byte, error) {
	defer c.Close()
	return redis.Bytes(c.Do("PING"))
}
