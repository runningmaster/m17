package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"

	"main/command"
)

func main() {
	flag.Parse()

	ctx := context.Background()
	log := log.New(os.Stderr, "", log.LstdFlags)
	if isSystemdBasedOS() {
		log.SetFlags(0)
	}

	code, err := command.Execute(
		ctx,
		command.Logger(log),
	)
	if err != nil {
		log.Printf("%v", err)
	}
	os.Exit(code)
}

// isSystemdBasedOS returns true if systmd is running.
func isSystemdBasedOS() bool {
	return exec.Command("pidof", "systemd").Run() == nil
}
