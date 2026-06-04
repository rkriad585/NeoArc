package main

import (
	"neoarc/internal/cli"
	"neoarc/internal/version"
	"os"
)

var (
	Version        string
	Commit         string
	PublisherName  string
	PublisherEmail string
)

func init() {
	if Version != "" {
		cli.Version = Version
		version.Version = Version
	}
	if Commit != "" {
		cli.Commit = Commit
		version.Commit = Commit
	}
}

func main() {
	os.Exit(cli.Run(os.Args))
}
