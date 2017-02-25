package api

import (
	"context"
	"fmt"
	"net/http"

	"internal/router"

	"github.com/garyburd/redigo/redis"
)

type redisConner interface {
	Get() redis.Conn
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

func ping(c redisConner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		res, err := redisPing(ctx, c.Get())
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
