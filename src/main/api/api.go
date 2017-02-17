package api

import (
	"context"
	"net/http"

	"fmt"
	"main/logger"
	"main/mdware"
	"main/router"
)

// Init inits package.
func New(ctx context.Context, l logger.Logger) (http.Handler, error) {
	r, err := router.New(ctx, "httprouter")
	if err != nil {
		return nil, err
	}
	// make redis pool here
	registerHTTPHandlers(r)
	return r, nil
}

func registerHTTPHandlers(r router.HTTPRouter) {
	head := mdware.Head(nil)
	tail := mdware.Tail
	r.Add("GET", "/test/:name/foo", mdware.Pipe(head, mdware.Wrap(test), tail))

	r.Set404(mdware.Pipe(mdware.Wrap(err4xx(http.StatusNotFound)), mdware.Tail))
	r.Set405(mdware.Pipe(mdware.Wrap(err4xx(http.StatusMethodNotAllowed)), mdware.Tail))
}

func test(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	key := router.ParamKey{Name: "name"}
	v, _ := ctx.Value(key).(string)
	fmt.Fprintf(w, "Hello, World! From test handler! Param: %s\n", v)
	*r = *r.WithContext(ctx)
}

func err4xx(code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		msg := fmt.Sprintf("%d %s", code, http.StatusText(code))
		http.Error(w, msg, code)
		*r = *r.WithContext(ctx)
	}
}
