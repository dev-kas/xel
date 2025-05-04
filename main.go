package main

import (
	"fmt"
	"os"

	"xel/cmds"

	"github.com/urfave/cli/v2"
)

// Version is the current version of Xel
// This will be overridden during build by using ldflags
var Version = "dev"

func main() {
	app := &cli.App{
		Name:     "xel",
		Usage:    "A runtime for VirtLang",
		Version:  Version,
		Commands: cmds.GetCommands(),
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
