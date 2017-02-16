package command

import (
	"context"
	"flag"
	"time"

	"main/client"
	"main/router"
	"main/server"
	"main/version"

	"github.com/google/subcommands"
)

type serverCommand struct {
	baseCommand
	flagAddr    string
	flagRedis   string
	flagSecret  string
	flagMaxIdle int
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
	f.IntVar(&c.flagMaxIdle,
		"maxidle",
		128,
		"Server timeout",
	)
	f.DurationVar(&c.flagTimeout,
		"timeout",
		60*time.Second,
		"Server timeout",
	)
}

func (c *serverCommand) execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) error {
	_, err := client.NewRedisPool(c.flagRedis, c.flagMaxIdle, c.flagTimeout)
	if err != nil {
		return err
	}

	r, err := router.NewHTTPRouter()
	if err != nil {
		return err
	}

	g, err := server.NewHTTPGraceful(c.flagAddr, version.AppName(), c.flagTimeout, r, nil)
	if err != nil {
		return err
	}

	return g.ListenAndServe()
}
