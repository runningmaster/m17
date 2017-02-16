package server

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/tylerb/graceful"
)

// NewHTTPGraceful returns *graceful.Server with http.Handler.
// TODO: change *log.Logger to interface when Go1.9 will released.
func NewHTTPGraceful(addr, name string, d time.Duration, h http.Handler, l *log.Logger) (*graceful.Server, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	if l != nil {
		l.Printf("%s is now ready to accept connections on %s", name, u.Host)
	}

	return &graceful.Server{
			Server: &http.Server{
				Addr:    u.Host,
				Handler: h,
			},
			Timeout: d,
			Logger:  l,
		},
		nil
}
