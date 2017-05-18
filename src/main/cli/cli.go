package cli

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

func init() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
}

type flager interface {
	setFlags(*flag.FlagSet)
}

type executer interface {
	execute(context.Context, *flag.FlagSet, ...interface{}) error
}

// Run finds and executes the specific command.
func Run() int {
	flag.Parse()
	ctx := context.Background()
	res := subcommands.Execute(ctx)
	return int(res)
}
