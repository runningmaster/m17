package api

import (
	"fmt"
	"net/http"
)

// pipe joins several middleware in one pipeline.
func pipe(pipes ...func(http.Handler) http.Handler) http.Handler {
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// tail code may be here
	})
	for i := len(pipes) - 1; i >= 0; i-- {
		h = pipes[i](h)
	}
	return h
}

// err4xx is wrapper for NotFound and MethodNotAllowed error handlers
func err4xx(code int) func(http.Handler) http.Handler {
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

// head puts to context FIXME
func head(uuid func() string, auth func(string) bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			println("head")
			ctx := r.Context()
			if uuid != nil {
				uuid()
			}
			if auth != nil {
				auth("")
			}

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println("gzip")
		ctx := r.Context()
		/*
			err := failFrom(ctx)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if gziputil.InString(r.Header.Get("Content-Encoding")) {
				z, err := gziputil.GetReader()
				if err != nil {
					ctx = withFail(ctx, err)
				}
				defer func() { _ = gziputil.PutReader(z) }()
				err = z.Reset(r.Body)
				if err != nil {
					ctx = withFail(ctx, err)
				}
				r.Body = z
			}

			if gziputil.InString(r.Header.Get("Accept-Encoding")) {
				z, err := gziputil.GetWriter()
				if err != nil {
					ctx = withFail(ctx, err)
				}
				defer func() { _ = gziputil.PutWriter(z) }()
				z.Reset(w)
				w = gziputil.ResponseWriter{Writer: z, ResponseWriter: w}
				w.Header().Add("Vary", "Accept-Encoding")
				w.Header().Set("Content-Encoding", "gzip")
			}
		*/
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// tail puts to context FIXME
func tail(log logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			log.Printf("tail\n")

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// tail ends habdlers pipeline.
//func tail(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		ctx := r.Context()
//		fmt.Println("TailMDWare")
//		r = r.WithContext(ctx)
//		next.ServeHTTP(w, r)
//	})
//}

// wrap is wrapper for user http.HandlerFunc
func wrap(h http.HandlerFunc) func(http.Handler) http.Handler {
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
