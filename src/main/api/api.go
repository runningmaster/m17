package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"main/logger"
	m "main/mdware"
	"main/router"
)

type routeTable map[string]http.Handler

var table = routeTable{
	"GET /:foo/bar":  m.Pipe(m.Head(nil), m.Wrap(test), m.Tail),
	"GET /test/:foo": m.Pipe(m.Head(nil), m.Wrap(test), m.Tail),
}

// New returns API inits package.
func New(ctx context.Context, l logger.Logger) (http.Handler, error) {

	// make redis pool here
	return makeHTTPRouter(ctx, table)
}

func makeHTTPRouter(ctx context.Context, t routeTable) (router.HTTPRouter, error) {
	r, err := router.New(ctx, router.MuxBone)
	if err != nil {
		return nil, err
	}

	for k, v := range t {
		s := strings.Split(k, " ") // [m,p]
		r.Add(s[0], s[1], v)
	}

	r.Set404(m.Pipe(m.Err4xx(http.StatusNotFound), m.Tail))
	r.Set405(m.Pipe(m.Err4xx(http.StatusMethodNotAllowed), m.Tail))

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
	*r = *r.WithContext(ctx)
}
