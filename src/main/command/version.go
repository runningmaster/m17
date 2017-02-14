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
	flagFull bool
}

func newVersionCommand() subcommands.Command {
	c := &versionCommand{
		baseCommand: baseCommand{
			name:  "version",
			synop: "print versionr",
			usage: "Print version to stdout",
		},
	}
	c.cmd = c
	return c
}

func (c *versionCommand) setFlags(f *flag.FlagSet) {
	f.BoolVar(&c.flagFull, "full", false, "print full version with build info")
}

func (c *versionCommand) execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) error {
	v := version.String()
	if c.flagFull {
		v = version.WithBuildInfo()
	}
	fmt.Println(v)
	return nil
}
