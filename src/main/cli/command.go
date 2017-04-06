package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
)

// baseCommand is base for another one.
type baseCommand struct {
	base  interface{}
	name  string
	brief string
	usage string
	logger
}

func (c *baseCommand) appName() string {
	return filepath.Base(os.Args[0])
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
	c.logger = makeLogger(isSystemdBasedOS())

	var err error
	if v, ok := c.base.(executer); ok {
		err = c.overrideFlagsEnv(f)
		if err == nil {
			err = v.execute(ctx, f, args...)
		}
	}

	if err != nil {
		c.logger.Printf("%v", err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
