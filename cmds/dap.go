package cmds

import (
	"os"
	"xel/dap"

	"github.com/urfave/cli/v2"
)

// DapCommand returns the cli.Command for the dap command
func DapCommand() *cli.Command {
	return &cli.Command{
		Name:  "dap",
		Usage: "Start the Debug Adapter Protocol server",
		Action: func(c *cli.Context) error {
			return dap.Serve(os.Stdin, os.Stdout)
		},
	}
}
