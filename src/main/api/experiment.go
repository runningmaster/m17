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

	"github.com/garyburd/redigo/redis"
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

func ping(rdb rediser) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		c := rdb.Get()
		defer c.Close()

		res, err := redis.String(c.Do("PING"))
		if err != nil {
			ctx = ctxutil.WithError(ctx, err)
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
