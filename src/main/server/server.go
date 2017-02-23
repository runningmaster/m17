package server

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

const defaultURL = "http://localhost:8080"

type logger interface {
	Printf(string, ...interface{})
}

// ListenAndServe starts HTTP server.
func ListenAndServe(ctx context.Context, l logger, options ...func(*http.Server) error) error {
	s := &http.Server{}
	err := Address(defaultURL)(s)
	if err != nil {
		return err
	}

	for i := range options {
		err = options[i](s)
		if err != nil {
			return err
		}
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	go listenForShutdown(ctx, l, s, ch)
	return s.ListenAndServe()
}

func listenForShutdown(ctx context.Context, l logger, s *http.Server, ch <-chan os.Signal) {
	log := l
	log.Printf("now ready to accept connections on %s", s.Addr)
	<-ch

	log.Printf("\n")
	log.Printf("trying to shutdown...")

	err := s.Shutdown(ctx)
	if err != nil {
		log.Printf("%v", err)
	}
}

// Address is TCP address to listen on, "http://localhost:8080" if empty
func Address(a string) func(*http.Server) error {
	return func(s *http.Server) error {
		u, err := url.Parse(a)
		if err != nil {
			return err
		}
		s.Addr = u.Host
		return nil
	}
}

// Handler to invoke, http.DefaultServeMux if nil
func Handler(h http.Handler) func(*http.Server) error {
	return func(s *http.Server) error {
		s.Handler = h
		return nil
	}
}

// IdleTimeout is the maximum amount of time to wait for the
// next request when keep-alives are enabled. If IdleTimeout
// is zero, the value of ReadTimeout is used. If both are
// zero, there is no timeout.
func IdleTimeout(d time.Duration) func(*http.Server) error {
	return func(s *http.Server) error {
		s.IdleTimeout = d
		return nil
	}
}
