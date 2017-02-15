package command

import (
	"context"
	"flag"

	"main/server"

	"github.com/google/subcommands"
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
	return server.ListenAndServe("")
}
