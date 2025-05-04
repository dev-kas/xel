package cmds

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

// RunCommand returns the cli.Command for the run command
func RunCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Run a file with arguments",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("filename is required")
			}
			
			filename := c.Args().Get(0)
			
			// Check if file exists
			_, err := os.Stat(filename)
			if os.IsNotExist(err) {
				return fmt.Errorf("file %s does not exist", filename)
			}
			
			// Get remaining arguments
			args := []string{}
			if c.NArg() > 1 {
				for i := 1; i < c.NArg(); i++ {
					args = append(args, c.Args().Get(i))
				}
			}
			
			// For now, just print the file and args
			fmt.Printf("File: %s\n", filename)
			if len(args) > 0 {
				fmt.Printf("Args: %s\n", strings.Join(args, ":"))
			}
			
			return nil
		},
	}
}