package logger

import (
	"log"
	"os"
	"os/exec"
)

var _ Logger = NewDefault()

// Logger is common interface for logging. Check when go1.9 will be released.
// see https://groups.google.com/forum/#!topic/golang-dev/F3l9Iz1JX4g .
type Logger interface {
	Printf(string, ...interface{})
}

// NewDefault returns default logger
func NewDefault() Logger {
	l := log.New(os.Stderr, "", log.LstdFlags)
	if isSystemdBasedOS() {
		l.SetFlags(0)
	}
	return l
}

// systemdBasedOS returns true if systmd is running.
func isSystemdBasedOS() bool {
	return exec.Command("/usr/bin/pidof", "systemd").Run() == nil || exec.Command("/bin/pidof", "systemd").Run() == nil
}
