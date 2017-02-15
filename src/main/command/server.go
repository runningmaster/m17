package command

import (
	"context"
	"flag"
	"log"
	"net/http"
	"net/url"
	"time"

	"main/version"

	"github.com/google/subcommands"
	"github.com/tylerb/graceful"
)

type serverCommand struct {
	baseCommand
	flagAddr   string
	flagRedis  string
	flagSecret string
}

func newServerCommand() subcommands.Command {
	c := &serverCommand{
		baseCommand: baseCommand{
			name:  "server",
			brief: "start server",
			usage: "Start HTTP server",
		},
	}
	c.cmd = c
	return c
}

func (c *serverCommand) setFlags(f *flag.FlagSet) {
	f.StringVar(&c.flagAddr,
		"addr",
		"http://127.0.0.1:8080",
		"Host server addres",
	)
	f.StringVar(&c.flagRedis,
		"redis",
		"redis://127.0.0.1:6379",
		"Redis server address",
	)
	f.StringVar(&c.flagSecret,
		"secret",
		"masterkey",
		"Default secret key",
	)
}

func (c *serverCommand) execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) error {
	s, err := newServer(c.flagAddr, nil, nil)
	if err != nil {
		return err
	}
	return s.ListenAndServe()
}

// newServer returns *graceful.Server.
func newServer(addr string, h http.Handler, l *log.Logger) (*graceful.Server, error) {
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
			Timeout: 5 * time.Second,
			Logger:  l,
		},
		nil
}
