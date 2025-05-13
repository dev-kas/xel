package cmds

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"xel/helpers"
	xShared "xel/shared"

	"github.com/dev-kas/virtlang-go/v2/environment"
	"github.com/dev-kas/virtlang-go/v2/errors"
	"github.com/dev-kas/virtlang-go/v2/evaluator"
	"github.com/dev-kas/virtlang-go/v2/parser"
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

			// Convert filename to absolute path
			filename, err = filepath.Abs(filename)
			if err != nil {
				return err
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

			program, parseErr := parser.New().ProduceAST(string(content))
			if parseErr != nil {
				if err, ok := parseErr.(*errors.SyntaxError); ok {
					start := err.Start + 1
					end := start + err.Difference

					startLine, startCol := helpers.PosToLineCol(content, start)
					endLine, endCol := helpers.PosToLineCol(content, end)

					color.Red("From %s:%d:%d to %s:%d:%d\n", filename, startLine, startCol, filename, endLine, endCol)
				}
				return parseErr
			}

			rootEnv := &xShared.XelRootEnv
			RV_proc := map[string]*shared.RuntimeValue{}

			if _, ok := rootEnv.Resolve("proc"); ok != nil {
				rootEnv.DeclareVar("proc", values.MK_OBJECT(RV_proc), false)
			} else {
				val, _ := rootEnv.LookupVar("proc")
				RV_proc = val.Value.(map[string]*shared.RuntimeValue)
			}

			RV_proc_args := values.MK_ARRAY(args)
			RV_proc["args"] = &RV_proc_args

			rootEnv.AssignVar("proc", values.MK_OBJECT(RV_proc))

			env := environment.NewEnvironment(rootEnv)
			env.DeclareVar("__filename__", values.MK_STRING(fmt.Sprintf("\"%s\"", filename)), true)
			env.DeclareVar("__dirname__", values.MK_STRING(fmt.Sprintf("\"%s\"", filepath.Dir(filename))), true)

			_, evalErr := evaluator.Evaluate(program, &env)

			if evalErr != nil {
				return evalErr
			}

			return nil
		},
	}
}
