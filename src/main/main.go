package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"main/command"
)

func main() {
	flag.Parse()

	ctx := context.Background()
	code, err := command.Execute(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
	}
	os.Exit(code)
}
