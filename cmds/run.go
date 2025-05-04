package cmds

import (
	"fmt"
	"os"
	"strings"
	"xel/engine"

	"github.com/urfave/cli/v2"
)

// RunCommand returns the cli.Command for the run command
func RunCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Execute a VirtLang script file",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("filename is required")
			}

			filename := c.Args().Get(0)

			// Check if file has .xel extension
			if !strings.HasSuffix(filename, ".xel") {
				return fmt.Errorf("file must have .xel extension")
			}

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

			// Read file content
			content, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %v", filename, err)
			}

			// Pass content to the engine
			result, err := engine.Eval(string(content))
			if err != nil {
				return fmt.Errorf("execution failed: %v", err)
			}

			// Print the result
			if result != nil {
				// Check if the result has a Value field (assuming RuntimeValue has a Value field)
				if v, ok := result.Value.(int); ok {
					fmt.Printf("Result: %d\n", v)
				} else if v, ok := result.Value.(float64); ok {
					fmt.Printf("Result: %f\n", v)
				} else if v, ok := result.Value.(string); ok {
					fmt.Printf("Result: %s\n", v)
				} else if v, ok := result.Value.(bool); ok {
					fmt.Printf("Result: %t\n", v)
				} else {
					fmt.Printf("Result: %v\n", result)
				}
			} else {
				fmt.Println("Execution completed successfully with no return value")
			}

			return nil
		},
	}
}