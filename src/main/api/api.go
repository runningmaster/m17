package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"main/router"
)

type API struct {
	mux    map[string]http.Handler
	err404 http.Handler
	err405 http.Handler
}

func (a *API) WithRouter(r router.Router) (router.Router, error) {
	var s []string
	var err error
	for k, v := range a.mux {
		s = strings.Split(k, " ")
		if len(s) != 2 {
			panic("invalid pair method-path")
		}
		err = r.Add(s[0], s[1], v)
		if err != nil {
			return nil, err
		}
	}

	err = r.Set404(a.err404)
	if err != nil {
		return nil, err
	}

	err = r.Set405(a.err405)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// New returns API.
func New(ctx context.Context, options ...func(*Option) error) (*API, error) {
	err := defaultOption.override(options...)
	if err != nil {
		return nil, err
	}

	if defaultOption.redisPool == nil {
		panic("must redis")
	}

	return &API{
			mux:    multiplexer,
			err404: err404,
			err405: err405,
		},
		nil
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
