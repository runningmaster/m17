package cli

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"

	"github.com/google/subcommands"
)

var _ logger = makeLogger(isSystemdBasedOS()) // check interface

func init() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
}

// logger is common interface for logging. Check when go1.9 will be released.
// see https://groups.google.com/forum/#!topic/golang-dev/F3l9Iz1JX4g .
type logger interface {
	Printf(string, ...interface{})
}

type flagSetter interface {
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

// systemdBasedOS returns true if systmd is running.
func isSystemdBasedOS() bool {
	return exec.Command("/usr/bin/pidof", "systemd").Run() == nil
}

// short means without timestamp
func makeLogger(short bool) logger {
	l := log.New(os.Stderr, "", log.LstdFlags)

	if short {
		l.SetFlags(0)
	}

	return l
}
