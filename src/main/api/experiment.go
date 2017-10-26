package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"main/version"

	"internal/ctxutil"
	"internal/router"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, version.WithBuildInfo(), runtime.Version())
}

func help() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		res, err := ioutil.ReadFile(filepath.Join(filepath.Dir(os.Args[0]), version.AppName()+".txt"))
		if err != nil {
			ctx = ctxutil.WithError(ctx, err)
			*r = *r.WithContext(ctx)
			return
		}
		ctx = ctxutil.WithResult(ctx, res)
		*r = *r.WithContext(ctx)
	})
}

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
		ctx = dbx.ctx // get ctx from func
		if err != nil {
			ctx = ctxutil.WithError(ctx, err)
		}

		ctx = ctxutil.WithResult(ctx, res)
		*r = *r.WithContext(ctx)
	})
}
