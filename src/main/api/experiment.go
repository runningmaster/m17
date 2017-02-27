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

func stdh(w http.ResponseWriter, r *http.Request) {
	if h, p := http.DefaultServeMux.Handler(r); p != "" {
		h.ServeHTTP(w, r)
	}
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

		v, err := redisPing(ctx, c.Get())
		if err != nil {
			ctx = contextWithError(ctx, err, http.StatusInternalServerError)
		}
		ctx = contextWithResult(ctx, v)
		*r = *r.WithContext(ctx)
	}
}

func redisPing(_ context.Context, c redis.Conn) (interface{}, error) {
	defer c.Close()
	return c.Do("PING")
}
