package cmds

import (
	"fmt"
	"os"
	"strings"
	"xel/engine"
	"xel/globals"
	"xel/helpers"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/shared"
	"github.com/dev-kas/virtlang-go/v2/values"
	"github.com/fatih/color"
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
				args[i] = values.MK_STRING(fmt.Sprintf("\"%s\"", arg))
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
				if ok := err.(*errors.SyntaxError); ok != nil {
					start := ok.Start + 1
					end := start + ok.Difference

					startLine, startCol := helpers.PosToLineCol(content, start)
					endLine, endCol := helpers.PosToLineCol(content, end)

					color.Red("From %s:%d:%d to %s:%d:%d\n", filename, startLine, startCol, filename, endLine, endCol)
				}
				return err
			}

			return nil
		},
	}
}
