package middleware

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
)

type logger interface {
	Printf(string, ...interface{})
}

// Pipe does
type Pipe struct {
	before []func(http.Handler) http.Handler
	// Join(...) here
	after []func(http.Handler) http.Handler
}

// BeforeJoin joins middleware handlers in order executing before Join method.
func (p *Pipe) BeforeJoin(pipes ...func(http.Handler) http.Handler) {
	for i := range pipes {
		p.before = append(p.before, pipes[i])
	}
}

// AfterJoin joins middleware handlers in order executing before Join method.
func (p *Pipe) AfterJoin(pipes ...func(http.Handler) http.Handler) {
	for i := range pipes {
		p.after = append(p.after, pipes[i])
	}
}

// Join joins several middlewares in one pipeline.
func (p *Pipe) Join(pipes ...func(http.Handler) http.Handler) http.Handler {
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	for i := len(p.after) - 1; i >= 0; i-- {
		h = p.after[i](h)
	}
	for i := len(pipes) - 1; i >= 0; i-- {
		h = pipes[i](h)
	}
	for i := len(p.before) - 1; i >= 0; i-- {
		h = p.before[i](h)
	}

	return h
}

// Head does some actions the first in handlers pipeline.  Must be first in pipeline.
func Head(uuidFn func() string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = contextWithTime(ctx, time.Now())
			ctx = contextWithHost(ctx, mineHost(r))
			ctx = contextWithUser(ctx, r.UserAgent())

			if uuidFn != nil {
				uuid := uuidFn()
				ctx = contextWithUUID(ctx, uuid)
				w.Header().Set("X-Request-ID", uuid)
			}

			w.Header().Set("X-Powered-By", fmt.Sprintf("go version %s", runtime.Version()))
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

// Auth checks user's access to service.
func Auth(authFn func(*http.Request) (string, int, error)) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if authFn != nil {
				ctx := r.Context()
				name, code, err := authFn(r)
				if err != nil {
					ctx = contextWithError(ctx, err, code)
				}
				ctx = contextWithAuth(ctx, name)
				r = r.WithContext(ctx)
			}

			h.ServeHTTP(w, r)
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

// Gzip wraps reader and writer for decompress and ompress data.
func Gzip(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			ctx := r.Context()
			err := errorFromContext(ctx)
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}

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

		h.ServeHTTP(w, r)
	})
}

// Body reads data from Request.Body into context.Conext.
func Body(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errorFromContext(ctx)
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}

		if !strings.HasPrefix(r.Method, "P") {
			h.ServeHTTP(w, r)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ctx = contextWithError(ctx, err, http.StatusBadRequest)
		}
		ctx = contextWithClen(ctx, int64(len(b)))
		ctx = contextWithBody(ctx, b)
		_ = r.Body.Close()

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

// Exec execites main user handler for registared URL.
func Exec(v interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			err := errorFromContext(ctx)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			switch h := v.(type) {
			case http.HandlerFunc:
				h(w, r)
			case http.Handler:
				h.ServeHTTP(w, r)
			default:
				panic("unknown handler")
			}

			next.ServeHTTP(w, r)
		})
	}
}

// JSON makes JSON from data. Must be before response.
func JSON(h http.Handler) http.Handler {
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
				ctx = contextWithError(ctx, err)
			} else {
				ctx = contextWithData(ctx, b)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
			}
		} else {
			if v, ok := res.([]byte); !ok {
				ctx = contextWithError(ctx, fmt.Errorf("result must be []byte"))
			} else {
				ctx = contextWithData(ctx, v)
			}
		}

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

// Resp writes result data to response.
func Resp(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := errorFromContext(ctx)
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}

		// skip if stdh executed
		if cl := w.Header().Get("Content-Length"); cl != "" {
			var val int64
			val, err = strconv.ParseInt(cl, 10, 64)
			if err != nil {
				ctx = contextWithError(ctx, err)
			}
			ctx = contextWithSize(ctx, val)

			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
			return
		}

		data := dataFromContext(ctx)
		if data == nil {
			ctx = contextWithError(ctx, fmt.Errorf("%v data", data))

			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
			return
		}

		code := codeFromContext(ctx)
		w.WriteHeader(code)

		n, err := w.Write(data)
		if err != nil {
			ctx = contextWithError(ctx, err)
		}
		_, err = w.Write([]byte("\n"))
		if err != nil {
			ctx = contextWithError(ctx, err)
		} else {
			n++
		}
		ctx = contextWithSize(ctx, int64(n))

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

// Fail writes error message to response. Must be after resp.
func Fail(h http.Handler) http.Handler {
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
				ctx = contextWithError(ctx, err)
			}
			ctx = contextWithSize(ctx, int64(n))
			r = r.WithContext(ctx)
		}

		h.ServeHTTP(w, r)
	})
}

// ErrCode is wrapper for NotFound and MethodNotAllowed error handlers.
func ErrCode(code int) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = contextWithError(ctx, fmt.Errorf("router says"), code)

			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

// Tail does some last actions (logging, send metrics). Must be in the end of pipe.
func Tail(log logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if log == nil {
				panic(fmt.Sprintf("%v log", log))
			}

			ctx := r.Context()
			err := errorFromContext(ctx)
			var errText string
			if err != nil {
				errText = fmt.Sprintf(" %s", err.Error())
			}
			log.Printf(
				"%s %s %s %s %s %d %d %d%s\n",
				timeFromContext(ctx),
				hostFromContext(ctx),
				userFromContext(ctx),
				uuidFromContext(ctx),
				authFromContext(ctx),
				clenFromContext(ctx),
				sizeFromContext(ctx),
				codeFromContext(ctx),
				errText,
			)

			h.ServeHTTP(w, r)
		})
	}
}

// Stdh executes standard handlers regestered in http.DefaultServeMux.
func Stdh(w http.ResponseWriter, r *http.Request) {
	if h, p := http.DefaultServeMux.Handler(r); p != "" {
		h.ServeHTTP(w, r)
	}
}
