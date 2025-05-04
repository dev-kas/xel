package cmds

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// InstallCommand returns the cli.Command for the install command
func InstallCommand() *cli.Command {
	return &cli.Command{
		Name:  "install",
		Usage: "Install a package",
		Action: func(c *cli.Context) error {
			fmt.Println("UNIMPLEMENTED FEATURE")
			return nil
		},
	}
}