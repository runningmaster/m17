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

	log := logger.New()
	ctx := logger.ContextWithLogger(context.Background(), log)

	code, err := command.Execute(ctx)
	if err != nil {
		log.Printf("%v", err)
	}
	os.Exit(code)
}
