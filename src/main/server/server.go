package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
)

// ListenAndServe starts *http.Server with http.Handler.
func ListenAndServe(ctx context.Context, options ...func(*Option) error) error {
	err := defaultOption.override(options...)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:        defaultOption.address,
		IdleTimeout: defaultOption.timeout,
		Handler:     defaultOption.handler,
	}

	log := defaultOption.log
	log.Printf("now ready to accept connections on %s", srv.Addr)

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	go listenForShutdown(ctx, srv, log, ch)

	return srv.ListenAndServe()
}

func listenForShutdown(ctx context.Context, srv *http.Server, log logger, ch <-chan os.Signal) {
	<-ch
	log.Printf("\n")
	log.Printf("trying to shutdown...")
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Printf("%v", err)
	}
}
