package main

import (
	"os"

	_ "expvar"
	_ "net/http/pprof"

	"main/cli"
)

func main() {
	os.Exit(cli.Run())
}
