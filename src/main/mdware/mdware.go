package mdware

import (
	"fmt"
	"net/http"
)

// Pipe joins several middleware in one pipeline.
// Pipe treats nil as http.DefaultServeMux.
func Pipe(pipes ...func(http.Handler) http.Handler) http.Handler {
	var h http.Handler
	for i := len(pipes) - 1; i >= 0; i-- {
		h = pipes[i](h)
	}

	if h == nil {
		h = http.DefaultServeMux
	}

	return h
}

// Head puts to context FIXME
func Head(genUUID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			var uuid string
			if genUUID != nil {
				uuid = genUUID()
			}
			fmt.Println("HeadMDWare", uuid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Tail ends habdlers pipeline.
func Tail(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fmt.Println("TailMDWare")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Wrap is wrapper for user http.HandlerFunc
func Wrap(h http.HandlerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			h(w, r) // stdh

			*r = *r.WithContext(ctx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
