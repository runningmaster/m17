package api

import (
	"context"
	"fmt"
	"net/http"

	"main/router"
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

func ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := redisPing(ctx)
	if err != nil {
		fmt.Fprintf(w, "redis error: %v\n", err)

	}
	fmt.Fprintf(w, "redis result: %v\n", string(res))
	*r = *r.WithContext(ctx)
}

func redisPing(ctx context.Context) ([]byte, error) {
	return nil, nil
}
