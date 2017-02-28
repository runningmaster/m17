package api

import (
	"context"
	"fmt"
	"net/http"

	m "internal/middleware"
	"internal/router"
)

func test(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Fprintf(w, "Hello, World! From test handler!\n")
	fmt.Fprintf(w, "Param foo: %s\n", router.ParamValueFromContext(ctx, "foo"))
	fmt.Fprintf(w, "Query foo: %s\n", router.QueryValueFromContext(ctx, "foo"))
	v, _ := ctx.Value("foo").(string)
	fmt.Fprintf(w, "Value foo: %s\n", v)

	*r = *r.WithContext(ctx)
}

func ping(rdb rediser) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		res, err := pingRedis(ctx, rdb)
		ctx = m.ContextWithError(ctx, err, http.StatusInternalServerError)
		ctx = m.ContextWithResult(ctx, res)
		*r = *r.WithContext(ctx)
	})
}

func pingRedis(_ context.Context, rdb rediser) (interface{}, error) {
	c := rdb.Get()
	defer c.Close()
	return c.Do("PING")
}
