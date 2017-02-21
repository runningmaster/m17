package command

import (
	"context"
	"flag"
	"fmt"
	"os"

	"main/logger"
	"main/option"
	"main/version"

	"github.com/google/subcommands"
)

func init() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(newServerCommand(), "")
	subcommands.Register(newVersionCommand(), "")
}

type flagSetter interface {
	setFlags(*flag.FlagSet)
}

type executer interface {
	execute(context.Context, *flag.FlagSet, ...interface{}) error
}

// baseCommand is base for another one.
type baseCommand struct {
	base  interface{}
	name  string
	brief string
	usage string
}

func (c *baseCommand) appName() string {
	return version.AppName()
}

// Name returns the name of the command.
func (c *baseCommand) Name() string {
	return c.name
}

// Synopsis returns a short string (less than one line) describing the command.
func (c *baseCommand) Synopsis() string {
	return c.brief
}

// Usage returns a long string explaining the command and giving usage
// information.
func (c *baseCommand) Usage() string {
	return fmt.Sprintf("%s [<flags>]:\n\t%s\n", c.Name(), c.usage)
}

// SetFlags adds the flags for this command to the specified set.
func (c *baseCommand) SetFlags(f *flag.FlagSet) {
	if v, ok := c.base.(flagSetter); ok {
		v.setFlags(f)
	}
}

// overrideFlagsEnv overrides flags from environment variables.
func (c *baseCommand) overrideFlagsEnv(f *flag.FlagSet) error {
	var err error
	f.VisitAll(func(f *flag.Flag) {
		env := os.Getenv(f.Name)
		if env != "" {
			err = f.Value.Set(env)
			if err != nil {
				return
			}
		}
	})

	return err
}

// Execute executes the command and returns an ExitStatus.
func (c *baseCommand) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	var err error
	if v, ok := c.base.(executer); ok {
		err = c.overrideFlagsEnv(f)
		if err == nil {
			err = v.execute(ctx, f, args...)
		}
	}

	if err != nil {
		errExec = err
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

// workaround for getting error from specific commands.
var errExec error
var logExec logger.Logger

// Execute finds and executes the specific command.
func Execute(ctx context.Context, options ...option.Fn) (int, error) {
	opt := &optionReceiver{}
	err := opt.Receive(options...)
	if err != nil {
		return int(subcommands.ExitFailure), err
	}
	logExec = opt.log

	status := subcommands.Execute(ctx)
	return int(status), errExec
}