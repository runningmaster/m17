package command

import (
	"context"
	"flag"
	"time"

	"main/api"
	"main/client"
	"main/server"

	"github.com/google/subcommands"
)

type serverCommand struct {
	baseCommand
	flag struct {
		addr    string
		redis   string
		secret  string
		maxIdle int
		timeout time.Duration
	}
}

func newServerCommand() subcommands.Command {
	c := &serverCommand{
		baseCommand: baseCommand{
			name:  "server",
			brief: "start server",
			usage: "Start HTTP server",
		},
	}
	c.base = c
	return c
}

func (c *serverCommand) setFlags(f *flag.FlagSet) {
	f.StringVar(&c.flag.addr,
		"addr",
		"http://127.0.0.1:8080",
		"Host server addres",
	)
	f.StringVar(&c.flag.redis,
		"redis",
		"redis://127.0.0.1:6379",
		"Redis server address",
	)
	f.StringVar(&c.flag.secret,
		"secret",
		"masterkey",
		"Default secret key",
	)
	f.IntVar(&c.flag.maxIdle,
		"maxidle",
		128,
		"Server timeout",
	)
	f.DurationVar(&c.flag.timeout,
		"timeout",
		60*time.Second,
		"Server timeout",
	)
}

func (c *serverCommand) execute(ctx context.Context, _ *flag.FlagSet, _ ...interface{}) error {
	_, err := client.NewRedisPool(c.flag.redis, c.flag.maxIdle, c.flag.timeout)
	if err != nil {
		return err
	}

	h, err := api.New(ctx, nil)
	if err != nil {
		return err
	}

	g, err := server.New(c.flag.addr, c.appName(), c.flag.timeout, h, nil)
	if err != nil {
		return err
	}

	return g.ListenAndServe()
}
