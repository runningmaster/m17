package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"main/logger"
	"main/option"
)

// ListenAndServe returns *graceful.Server with http.Handler.
func ListenAndServe(ctx context.Context, options ...option.Fn) error {
	opt := &optionReceiver{}
	err := opt.Receive(options...)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:        opt.address,
		IdleTimeout: opt.timeout,
		Handler:     opt.handler,
	}

	log := opt.log
	log.Printf("now ready to accept connections on %s", srv.Addr)

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	go listenForShutdown(ctx, srv, log, ch)

	return srv.ListenAndServe()
}

func listenForShutdown(ctx context.Context, srv *http.Server, log logger.Logger, ch <-chan os.Signal) {
	<-ch
	log.Printf("\n")
	log.Printf("trying to shutdown...")
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Printf("%v", err)
	}
}
