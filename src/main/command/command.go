package command

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/google/subcommands"
)

// workaround
var errExec error

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

// baseCommand is base for another one
type baseCommand struct {
	cmd   interface{}
	name  string
	synop string
	usage string
	err   error
}

// Execute executes the command.
func Execute(ctx context.Context) (int, error) {
	status := subcommands.Execute(ctx)
	return int(status), errExec
}

// Name returns the name of the command.
func (c *baseCommand) Name() string {
	return c.name
}

// Synopsis returns a short string (less than one line) describing the command.
func (c *baseCommand) Synopsis() string {
	return c.synop
}

// Usage returns a long string explaining the command and giving usage
// information.
func (c *baseCommand) Usage() string {
	return fmt.Sprintf(`%s:
	%s
`, c.Name(), c.usage)
}

// SetFlags adds the flags for this command to the specified set.
func (c *baseCommand) SetFlags(f *flag.FlagSet) {
	if v, ok := c.cmd.(flagSetter); ok {
		v.setFlags(f)
	}
}

// Execute executes the command and returns an ExitStatus.
func (c *baseCommand) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	var err error
	if v, ok := c.cmd.(executer); ok {
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

// overrideFlagsEnv overrides flags from environment variables
func (c *baseCommand) overrideFlagsEnv(f *flag.FlagSet) error {
	var err error
	f.Visit(func(f *flag.Flag) {
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
