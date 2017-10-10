package server

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"internal/logger"
)

const defaultURL = "http://localhost:8080"

// Server is wrapper for *http.Server with additional params.
type Server struct {
	ctx context.Context
	srv *http.Server
	log logger.Logger
}

// NewWithContext returns *Server with Context.
func NewWithContext(ctx context.Context, options ...func(*Server) error) (*Server, error) {
	if ctx == nil {
		panic("nil context")
	}

	s := &Server{
		ctx: ctx,
		srv: &http.Server{},
		log: logger.NewDefault(),
	}

	err := Address(defaultURL)(s)
	if err != nil {
		return nil, err
	}

	for i := range options {
		err = options[i](s)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Start starts HTTP server.
func (s *Server) Start() error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	go listenForShutdown(s, ch)

	err := s.srv.ListenAndServe()
	if err != nil && err == http.ErrServerClosed {
		return nil
	}
	return err
}

func listenForShutdown(s *Server, ch <-chan os.Signal) {
	s.log.Printf("now ready to accept connections on %s", s.srv.Addr)
	<-ch

	s.log.Printf("trying to shutdown...")
	err := s.srv.Shutdown(s.ctx)
	if err != nil {
		s.log.Printf("%v", err)
	}
}

// Logger is option for passing logger interface.
func Logger(l logger.Logger) func(*Server) error {
	return func(s *Server) error {
		s.log = l
		return nil
	}
}

// Address is TCP address to listen on, "http://localhost:8080" if empty.
func Address(a string) func(*Server) error {
	return func(s *Server) error {
		u, err := url.Parse(a)
		if err != nil {
			return err
		}
		s.srv.Addr = u.Host
		return nil
	}
}

// Handler to invoke, http.DefaultServeMux if nil.
func Handler(h http.Handler) func(*Server) error {
	return func(s *Server) error {
		s.srv.Handler = h
		return nil
	}
}

// IdleTimeout is the maximum amount of time to wait for the
// next request when keep-alives are enabled. If IdleTimeout
// is zero, the value of ReadTimeout is used. If both are
// zero, there is no timeout.
func IdleTimeout(d time.Duration) func(*Server) error {
	return func(s *Server) error {
		s.srv.IdleTimeout = d
		return nil
	}
}
