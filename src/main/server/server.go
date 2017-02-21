package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"main/option"
)

type logger interface {
	Printf(string, ...interface{})
}

type loggerFunc func()

func (l loggerFunc) Printf(f string, a ...interface{}) {
	fmt.Fprintf(ioutil.Discard, f, a)
}

type server struct {
	log logger
	srv *http.Server
}

func (s *server) setAddress(a string) error {
	u, err := url.Parse(a)
	if err != nil {
		return err
	}
	s.srv.Addr = u.Host
	return nil
}

func (s *server) setTimeout(t time.Duration) error {
	s.srv.IdleTimeout = t
	return nil
}

func (s *server) setHandler(h http.Handler) error {
	s.srv.Handler = h
	return nil
}

func (s *server) setLogger(l logger) error {
	s.log = l
	return nil
}

// ListenAndServe returns *graceful.Server with http.Handler.
func ListenAndServe(ctx context.Context, options ...option.Fn) error {
	s := &server{
		log: loggerFunc(nil),
		srv: &http.Server{},
	}

	err := option.Receive(s, options...)
	if err != nil {
		return err
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	go listenForShutdown(ctx, s, ch)

	s.log.Printf(" now ready to accept connections on %s", s.srv.Addr)
	return s.srv.ListenAndServe()
}

func listenForShutdown(ctx context.Context, s *server, ch <-chan os.Signal) {
	<-ch
	s.log.Printf("\ntrying to shutdown...")
	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.log.Printf("%v", err)
	}
}
