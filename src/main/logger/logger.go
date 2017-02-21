package logger

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// Logger is common interface for logging.
// See https://groups.google.com/forum/#!topic/golang-dev/F3l9Iz1JX4g for details.
type Logger interface {
	Printf(string, ...interface{})
}

type loggerFunc func()

func (l loggerFunc) Printf(f string, a ...interface{}) {
	fmt.Fprintf(ioutil.Discard, f, a)
}

// New returns the implementation of Logger interface.
func New() Logger {
	l := log.New(os.Stderr, "", log.LstdFlags)

	if isSystemdBasedOS() {
		l.SetFlags(0)
	}

	return l
}

// systemdBasedOS returns true if systmd is running.
func isSystemdBasedOS() bool {
	return exec.Command("pidof", "systemd").Run() == nil
}
