package api

import (
	"fmt"
	"net/http"

	"internal/ctxutil"
	"internal/router"
)

func test(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Fprintf(w, "Hello, World! From test handler!\n")
	fmt.Fprintf(w, "Param foo: %s\n", router.ParamValueFrom(ctx, "foo"))
	fmt.Fprintf(w, "Query foo: %s\n", router.QueryValueFrom(ctx, "foo"))
	v, _ := ctx.Value("foo").(string)
	fmt.Fprintf(w, "Value foo: %s\n", v)

	*r = *r.WithContext(ctx)
}

func ping(rdb rediser) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		res, err := newRedisHelper(ctx, rdb, nil).ping()
		if err != nil {
			ctx = ctxutil.WithError(ctx, err)
		}

		ctx = ctxutil.WithResult(ctx, res)
		*r = *r.WithContext(ctx)
	})
}

func upload(rdb rediser) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		res, err := newRedisHelper(ctx, rdb, nil).uploadData([]byte(r.Header.Get("Content-Meta")), ctxutil.BodyFrom(ctx))
		if err != nil {
			ctx = ctxutil.WithError(ctx, err)
		}

		ctx = ctxutil.WithResult(ctx, res)
		*r = *r.WithContext(ctx)
	})
}

func uploadSuggestion(rdb rediser) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		res, err := newRedisHelper(ctx, rdb, nil).uploadSuggestion([]byte(r.Header.Get("Content-Meta")), ctxutil.BodyFrom(ctx))
		if err != nil {
			ctx = ctxutil.WithError(ctx, err)
		}

		ctx = ctxutil.WithResult(ctx, res)
		*r = *r.WithContext(ctx)
	})
}

func selectSuggestion(rdb rediser) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		res, err := newRedisHelper(ctx, rdb, nil).selectSuggestion([]byte(r.Header.Get("Content-Meta")), ctxutil.BodyFrom(ctx))
		if err != nil {
			ctx = ctxutil.WithError(ctx, err)
		}

		ctx = ctxutil.WithResult(ctx, res)
		*r = *r.WithContext(ctx)
	})
}
