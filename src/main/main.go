package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"

	"main/command"
)

// logger is common interface for logging. Check when go1.9 will be released.
// see https://groups.google.com/forum/#!topic/golang-dev/F3l9Iz1JX4g .
type logger interface {
	Printf(string, ...interface{})
}

// common pattern:
// main -> command(flags, env) -> some deps -> router -> api -> server.ListenAndServe()
func main() {
	flag.Parse()

	ctx := context.Background()
	log := log.New(os.Stderr, "", log.LstdFlags)
	var _ logger = log // check interface

	if isSystemdBasedOS() {
		log.SetFlags(0)
	}

	code, err := command.Execute(ctx, log)
	if err != nil {
		log.Printf("%v", err)
	}
	os.Exit(code)
}

// systemdBasedOS returns true if systmd is running.
func isSystemdBasedOS() bool {
	return exec.Command("pidof", "systemd").Run() == nil
}
