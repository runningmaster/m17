package command

import (
	"context"
	"flag"
	"time"

	"main/api"

	"github.com/google/subcommands"
)

type serverCommand struct {
	baseCommand
	flagAddr    string
	flagRedis   string
	flagSecret  string
	flagTimeout time.Duration
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
	f.DurationVar(&c.flagTimeout,
		"timeout",
		5*time.Second,
		"Server timeout",
	)
}

func (c *serverCommand) execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) error {
	r, err := api.NewHTTPRouter(c.flagSecret, nil)
	if err != nil {
		return err
	}

	s, err := api.NewHTTPServer(c.flagAddr, r, c.flagTimeout, nil)
	if err != nil {
		return err
	}
	return s.ListenAndServe()
}
