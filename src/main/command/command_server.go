package command

import (
	"context"
	"flag"
	"time"

	"main/api"
	"main/redis"
	"main/server"

	"github.com/google/subcommands"
)

const (
	defaultRedisAddress = "redis://127.0.0.1:6379"
	defaultRedisMaxIdle = 128
	defaultRedisTimeout = 60 * time.Second

	defaultServerAddress = "http://127.0.0.1:8080"
	defaultServerSecret  = "masterkey"
	defaultServerTimeout = 60 * time.Second
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
		defaultServerAddress,
		"Host server addres",
	)
	f.StringVar(&c.flag.redis,
		"redis",
		defaultRedisAddress,
		"Redis server address",
	)
	f.StringVar(&c.flag.secret,
		"secret",
		defaultServerSecret,
		"Default secret key",
	)
	f.IntVar(&c.flag.maxIdle,
		"maxidle",
		defaultRedisMaxIdle,
		"Server timeout",
	)
	f.DurationVar(&c.flag.timeout,
		"timeout",
		defaultRedisTimeout,
		"Server timeout",
	)
}

func (c *serverCommand) execute(ctx context.Context, _ *flag.FlagSet, _ ...interface{}) error {
	r, err := redis.NewPool(
		ctx,
		redis.Address(c.flag.redis),
		redis.MaxIdle(c.flag.maxIdle),
		redis.Timeout(c.flag.timeout),
	)
	if err != nil {
		return err
	}

	h, err := api.NewHandler(
		ctx,
		api.Router()
		api.Redis(r),
		api.Logger(l),
	)
	if err != nil {
		return err
	}

	return server.ListenAndServe(
		ctx,
		server.Address(c.flag.addr),
		server.Timeout(c.flag.timeout),
		server.Handler(h),
		server.Logger(log),
	)
}
