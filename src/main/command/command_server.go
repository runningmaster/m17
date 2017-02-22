package command

import (
	"context"
	"flag"
	"time"

	"main/api"
	"main/redispool"
	"main/router"
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
	l := defaultOption.log
	p, err := redispool.New(
		ctx,
		redispool.Address(c.flag.redis),
		redispool.MaxIdle(c.flag.maxIdle),
		redispool.Timeout(c.flag.timeout),
		redispool.Logger(l),
	)
	if err != nil {
		return err
	}

	a, err := api.New(
		ctx,
		api.Redis(p),
		api.Logger(l),
	)
	if err != nil {
		return err
	}

	r, err := router.New(
		ctx,
		router.HTTPRouter,
		router.Logger(l),
	)
	if err != nil {
		return err
	}

	h, err := a.WithRouter(r)
	if err != nil {
		return err
	}

	return server.ListenAndServe(
		ctx,
		server.Address(c.flag.addr),
		server.Timeout(c.flag.timeout),
		server.Handler(h),
		server.Logger(l),
	)
}
