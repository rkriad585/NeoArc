package main

import (
	"neoarc/internal/cli"
	"os"
)

func main() {
	os.Exit(cli.Run(os.Args))
}
