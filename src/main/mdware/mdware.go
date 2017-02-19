package mdware

import (
	"fmt"
	"net/http"
)

// Pipe joins several middleware in one pipeline.
func Pipe(pipes ...func(http.Handler) http.Handler) http.Handler {
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	for i := len(pipes) - 1; i >= 0; i-- {
		h = pipes[i](h)
	}

	return h
}

// Err4xx is wrapper for NotFound and MethodNotAllowed error handlers
func Err4xx(code int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			msg := fmt.Sprintf("%d %s", code, http.StatusText(code))
			http.Error(w, msg, code)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
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

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// Tail ends habdlers pipeline.
func Tail(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fmt.Println("TailMDWare")
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// Wrap is wrapper for user http.HandlerFunc
func Wrap(h http.HandlerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			fmt.Println("WrapMDWare")
			h(w, r) // stdh

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
