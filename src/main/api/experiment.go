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
		dbx := &dbxHelper{
			ctx: ctx,
			rdb: rdb,
		}
		res, err := dbx.ping()
		if err != nil {
			ctx = ctxutil.WithError(ctx, err)
		}

		ctx = ctxutil.WithResult(ctx, res)
		*r = *r.WithContext(ctx)
	})
}

func exec(rdb rediser) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		dbx := &dbxHelper{
			ctx,
			rdb,
			nil,
			r,
			w,
			[]byte(r.Header.Get("Content-Meta")),
			ctxutil.BodyFrom(ctx),
		}

		res, err := dbx.exec(router.ParamValueFrom(ctx, "func"))
		if err != nil {
			ctx = ctxutil.WithError(ctx, err)
		}

		ctx = ctxutil.WithResult(ctx, res)
		*r = *r.WithContext(ctx)
	})
}
