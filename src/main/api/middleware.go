package api

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"internal/gzippool"

	"github.com/nats-io/nuid"
	//"github.com/rogpeppe/fastuuid"
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

// err4xx is wrapper for NotFound and MethodNotAllowed error handlers.
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

func mineHost(r *http.Request) string {
	var v string
	if v = r.Header.Get("X-Forwarded-For"); v == "" {
		if v = r.Header.Get("X-Real-IP"); v == "" {
			v = r.RemoteAddr
		}
	}

	v, _, _ = net.SplitHostPort(v)
	return v
}

func head(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errorFromContext(ctx)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx = contextWithTime(ctx, time.Now())
		ctx = contextWithHost(ctx, mineHost(r))
		ctx = contextWithUser(ctx, r.UserAgent())
		ctx = contextWithUUID(ctx, nuid.Next())

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errorFromContext(ctx)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// TODO:
		// authFunc()
		// here

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errorFromContext(ctx)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			z := gzippool.GetReader()
			defer gzippool.PutReader(z)

			_ = z.Reset(r.Body)
			r.Body = z
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			z := gzippool.GetWriter()
			defer gzippool.PutWriter(z)

			z.Reset(w)
			w = gzippool.NewResponseWriter(z, w)
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Add("Vary", "Accept-Encoding")
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func read(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errorFromContext(ctx)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ctx = contextWithError(ctx, err)
			println(err)
		}
		ctx = contextWithClen(ctx, int64(len(b)))
		ctx = contextWithData(ctx, b)
		_ = r.Body.Close()

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func body(h http.HandlerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			err := errorFromContext(ctx)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			h(w, r) // stdh

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func tail(log logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if log == nil {
				panic(fmt.Sprintf("%v log", log))
			}

			ctx := r.Context()
			err := errorFromContext(ctx)
			if err != nil {
				log.Printf("Error: %s\n", err.Error())
				return
			}

			log.Printf(
				"%s %s %s %s %d\n",
				timeFromContext(ctx),
				hostFromContext(ctx),
				userFromContext(ctx),
				uuidFromContext(ctx),
				clenFromContext(ctx),
			)

			// The End
			//r = r.WithContext(ctx)
			//next.ServeHTTP(w, r)
		})
	}
}
