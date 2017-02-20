package server

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"main/logger"
)

// ListenAndServe returns *graceful.Server with http.Handler.
func ListenAndServe(ctx context.Context, addr string, d time.Duration, h http.Handler) error {
	u, err := url.Parse(addr)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:        u.Host,
		Handler:     h,
		IdleTimeout: d,
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	go listenForShutdown(ctx, srv, ch)

	log := logger.FromContext(ctx)
	log.Printf(" now ready to accept connections on %s", u.Host)
	return srv.ListenAndServe()
}

func listenForShutdown(ctx context.Context, srv *http.Server, ch <-chan os.Signal) {
	<-ch
	log := logger.FromContext(ctx)
	log.Printf("\ntrying to shutdown...")
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Printf("%v", err)
	}
}
