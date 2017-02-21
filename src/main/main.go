package main

import (
	"context"
	"flag"
	"os"

	"main/command"
	"main/logger"
)

func main() {
	flag.Parse()

	ctx := context.Background()
	log := logger.New()

	code, err := command.Execute(ctx, command.Logger(log))
	if err != nil {
		log.Printf("%v", err)
	}
	os.Exit(code)
}
