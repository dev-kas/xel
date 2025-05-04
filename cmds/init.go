package cmds

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// InitCommand returns the cli.Command for the init command
func InitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize a new project",
		Action: func(c *cli.Context) error {
			fmt.Println("UNIMPLEMENTED FEATURE")
			return nil
		},
	}
}