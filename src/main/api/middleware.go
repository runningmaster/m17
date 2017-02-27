package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"internal/gzippool"

	"github.com/nats-io/nuid"
	//"github.com/rogpeppe/fastuuid"
)

type pipe struct {
	headPipe []func(http.Handler) http.Handler
	tailPipe []func(http.Handler) http.Handler
}

func (p *pipe) head(pipes ...func(http.Handler) http.Handler) {
	for i := range pipes {
		p.headPipe = append(p.headPipe, pipes[i])
	}
}

func (p *pipe) tail(pipes ...func(http.Handler) http.Handler) {
	for i := range pipes {
		p.tailPipe = append(p.tailPipe, pipes[i])
	}
}

// join joins several middleware in one pipeline.
func (p *pipe) join(pipes ...func(http.Handler) http.Handler) http.Handler {
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	for i := len(p.tailPipe) - 1; i >= 0; i-- {
		h = p.tailPipe[i](h)
	}
	for i := len(pipes) - 1; i >= 0; i-- {
		h = pipes[i](h)
	}
	for i := len(p.headPipe) - 1; i >= 0; i-- {
		h = p.headPipe[i](h)
	}

	return h
}

// errCode is wrapper for NotFound and MethodNotAllowed error handlers.
func errCode(code int) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = contextWithError(ctx, fmt.Errorf("router says"), code)

			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

func uuid(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = contextWithTime(ctx, time.Now())
		ctx = contextWithHost(ctx, mineHost(r))
		ctx = contextWithUser(ctx, r.UserAgent())
		uuid := nuid.Next()
		ctx = contextWithUUID(ctx, uuid)

		w.Header().Set("X-Powered-By", fmt.Sprintf("go version %s", runtime.Version()))
		w.Header().Set("X-Request-ID", uuid)

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
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

func auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// TODO: authFunc() here

		// code := http.StatusUnauthorized
		// err := errByCode(code)
		// ctx = contextWithError(ctx, err, code)

		// code = http.StatusForbidden
		// err := errByCode(code)
		// ctx = contextWithError(ctx, err, code)

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func gunzip(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errorFromContext(ctx)
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			z := gzippool.GetReader()
			defer gzippool.PutReader(z)
			_ = z.Reset(r.Body)
			r.Body = z
		}

		h.ServeHTTP(w, r)
	})
}

func gzip(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			z := gzippool.GetWriter()
			defer gzippool.PutWriter(z)
			z.Reset(w)
			w = gzippool.NewResponseWriter(z, w)
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Add("Vary", "Accept-Encoding")
		}

		h.ServeHTTP(w, r)
	})
}

func body(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errorFromContext(ctx)
		if err != nil || r.Method != "POST" || r.Method != "PUT" {
			h.ServeHTTP(w, r)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ctx = contextWithError(ctx, err, http.StatusBadRequest)
		}
		ctx = contextWithClen(ctx, int64(len(b)))
		ctx = contextWithData(ctx, b)
		_ = r.Body.Close()

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func exec(h http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			err := errorFromContext(ctx)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			h.ServeHTTP(w, r)
			next.ServeHTTP(w, r)
		})
	}
}

func mrshl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errorFromContext(ctx)
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}

		res := resultFromContext(ctx)
		if w.Header().Get("Content-Type") == "" {
			b, err := json.Marshal(res)
			if err != nil {
				ctx = contextWithError(ctx, err, http.StatusInternalServerError)
			} else {
				ctx = contextWithResponse(ctx, b)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
			}
		} else {
			if v, ok := res.([]byte); !ok {
				ctx = contextWithError(ctx, fmt.Errorf("%v result", v), http.StatusInternalServerError)
			} else {
				ctx = contextWithResponse(ctx, v)
			}
		}

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func resp(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errorFromContext(ctx)
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}

		// skip if stdh executed
		if cl := w.Header().Get("Content-Length"); cl != "" {
			v, err := strconv.ParseInt(cl, 10, 64)
			if err != nil {
				ctx = contextWithError(ctx, err, http.StatusInternalServerError)
			}
			ctx = contextWithSize(ctx, v)

			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
			return
		}

		resp := responseFromContext(ctx)
		if resp == nil {
			ctx = contextWithError(ctx, fmt.Errorf("%v response", resp), http.StatusInternalServerError)

			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
			return
		}

		code := codeFromContext(ctx)
		w.WriteHeader(code)

		n, err := w.Write(resp)
		if err != nil {
			ctx = contextWithError(ctx, err, http.StatusInternalServerError)
		}
		_, err = w.Write([]byte("\n"))
		if err != nil {
			ctx = contextWithError(ctx, err, http.StatusInternalServerError)
		} else {
			n++
		}
		ctx = contextWithSize(ctx, int64(n))

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func errf(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errorFromContext(ctx)
		if err != nil {
			msg := fmt.Sprintf("%s\n", err.Error())
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")

			code := codeFromContext(ctx)
			w.WriteHeader(code)

			n, err := w.Write([]byte(msg))
			if err != nil {
				ctx = contextWithError(ctx, err, http.StatusInternalServerError)
			}
			ctx = contextWithSize(ctx, int64(n))
			r = r.WithContext(ctx)
		}

		h.ServeHTTP(w, r)
	})
}

func logg(log logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if log == nil {
				panic(fmt.Sprintf("%v log", log))
			}

			ctx := r.Context()
			log.Printf(
				"%s %s %s %s %d\n",
				timeFromContext(ctx),
				hostFromContext(ctx),
				userFromContext(ctx),
				uuidFromContext(ctx),
				clenFromContext(ctx),
			)

			h.ServeHTTP(w, r)
		})
	}
}
