package cli

import (
	"context"
	"flag"
	"time"

	"internal/redispool"
	"internal/router"
	"internal/server"

	"main/api"

	"github.com/google/subcommands"
)

func init() {
	subcommands.Register(newServerCommand(), "")
}

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
		"http://localhost:8080",
		"Host server addres",
	)
	f.StringVar(&c.flag.redis,
		"redis",
		"redis://localhost:6379",
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
	r, err := redispool.New(
		redispool.Address(c.flag.redis),
		redispool.MaxIdle(c.flag.maxIdle),
		redispool.IdleTimeout(c.flag.timeout),
	)
	if err != nil {
		return err
	}

	// ctx will be passed to http handlers via request
	h, err := api.NewWithRouter(
		router.NewMuxVestigo(ctx),
		api.Redis(r),
		api.Logger(c.log),
	)
	if err != nil {
		return err
	}

	// ctx will be passed to shutdown func
	s, err := server.NewWithContext(
		ctx,
		server.Handler(h),
		server.Address(c.flag.addr),
		server.IdleTimeout(c.flag.timeout),
		server.Logger(c.log),
	)
	if err != nil {
		return err
	}

	return s.Start()
}
