package command

import (
	"context"
	"flag"
	"fmt"

	"main/version"

	"github.com/google/subcommands"
)

type versionCommand struct {
	baseCommand
	flag struct {
		withMeta bool
	}
}

func newVersionCommand() subcommands.Command {
	c := &versionCommand{
		baseCommand: baseCommand{
			name:  "version",
			brief: "print version",
			usage: "Print version to stdout",
		},
	}
	c.base = c
	return c
}

func (c *versionCommand) setFlags(f *flag.FlagSet) {
	f.BoolVar(&c.flag.withMeta,
		"with-meta",
		false,
		"print full version with build metadata",
	)
}

func (c *versionCommand) execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) error {
	v := version.String()
	if c.flag.withMeta {
		v = version.WithBuildInfo()
	}
	fmt.Println(v)
	return nil
}
