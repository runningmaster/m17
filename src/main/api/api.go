package api

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"main/version"

	"github.com/tylerb/graceful"
)

// NewHTTPRouter returns router as http.Handler.
func NewHTTPRouter(secret string, l *log.Logger) (http.Handler, error) {
	return nil, nil
}

// NewHTTPServer returns *graceful.Server.
func NewHTTPServer(addr string, h http.Handler, t time.Duration, l *log.Logger) (*graceful.Server, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	if l != nil {
		l.Printf("%s is now ready to accept connections on %s", version.AppName(), u.Host)
	}

	return &graceful.Server{
			Server: &http.Server{
				Addr:    u.Host,
				Handler: h,
			},
			Timeout: t,
			Logger:  l,
		},
		nil
}
