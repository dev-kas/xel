package cmds

import (
	"fmt"
	"os"
	"strings"
	"xel/engine"
	"xel/globals"

	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/shared"
	"github.com/dev-kas/VirtLang-Go/values"
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
			rawArgs := []string{}
			if c.NArg() > 1 {
				for i := 1; i < c.NArg(); i++ {
					rawArgs = append(rawArgs, c.Args().Get(i))
				}
			}

			// Convert arguments to []shared.RuntimeValue
			args := make([]shared.RuntimeValue, len(rawArgs))
			for i, arg := range rawArgs {
				args[i] = values.MK_STRING(arg)
			}

			// Read file content
			content, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %v", filename, err)
			}

			// Pass content to the engine
			_, err = engine.Eval(string(content), func(env *environment.Environment) {
				globals.Globalize(env)

				RV_proc_args := values.MK_ARRAY(args)
				RV_proc := map[string]*shared.RuntimeValue{
					"args": &RV_proc_args,
				}
				convertedProc := make(map[string]shared.RuntimeValue)
				for key, val := range RV_proc {
					convertedProc[key] = *val
				}
				env.DeclareVar("proc", values.MK_OBJECT(convertedProc), false)
			})

			if err != nil {
				return fmt.Errorf("execution failed: %v", err)
			}

			return nil
		},
	}
}
