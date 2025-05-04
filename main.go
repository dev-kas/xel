package main

import (
	"fmt"
	"os"

	"xel/cmds"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "xel",
		Usage:    "A runtime for VirtLang",
		Commands: cmds.GetCommands(),
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
