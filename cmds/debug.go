package cmds

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"xel/helpers"
	xShared "xel/shared"

	"github.com/Masterminds/semver/v3"
	"github.com/chzyer/readline"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/evaluator"
	"github.com/dev-kas/virtlang-go/v4/parser"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

// DebugCommand returns the cli.Command for the debug command
// It executes a VirtLang script file in debug mode after performing version compatibility checks
// against the project's manifest (xel.json) if present.
func DebugCommand() *cli.Command {
	return &cli.Command{
		Name:  "debug",
		Usage: "Debug a VirtLang script file",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("filename is required")
			}

			filename := c.Args().Get(0)
			cwd, cwd_err := os.Getwd()
			if cwd_err != nil {
				return cwd_err
			}

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

			manifest := xShared.ProjectManifest{}

			// Check for xel.json in project directory
			manifestPath := filepath.Join(filepath.Dir(filename), "xel.json")
			if _, err := os.Stat(manifestPath); err == nil {
				// Parse manifest
				manifestContent, err := os.ReadFile(manifestPath)
				if err != nil {
					return fmt.Errorf("failed to read xel.json: %v", err)
				}

				if err := json.Unmarshal(manifestContent, &manifest); err != nil {
					return fmt.Errorf("failed to parse xel.json: %v", err)
				}

				// Check Xel version constraint if specified
				if manifest.Xel != "" {
					// Skip version check in development mode (when version is empty)
					if xShared.RuntimeVersion == "" {
						color.New(color.FgYellow).Printf("Skipping Xel version check in development mode\n")
					} else {
						constraint, err := semver.NewConstraint(manifest.Xel)
						if err != nil {
							return fmt.Errorf("invalid Xel version constraint in manifest: %v", err)
						}

						runtimeVersion, err := semver.NewVersion(xShared.RuntimeVersion)
						if err != nil {
							return fmt.Errorf("invalid runtime version format: %v", err)
						}

						if !constraint.Check(runtimeVersion) {
							return fmt.Errorf("xel version %s does not satisfy required version %s from xel.json, please upgrade your runtime",
								xShared.RuntimeVersion, manifest.Xel)
						}
					}
				}

				// Check Engine version constraint if specified
				if manifest.Engine != "" {
					// Skip version check in development mode (when version is empty)
					if xShared.EngineVersion == "" {
						color.New(color.FgYellow).Printf("Skipping Engine version check in development mode\n")
					} else {
						constraint, err := semver.NewConstraint(manifest.Engine)
						if err != nil {
							return fmt.Errorf("invalid Engine version constraint in manifest: %v", err)
						}

						engineVersion, err := semver.NewVersion(xShared.EngineVersion)
						if err != nil {
							return fmt.Errorf("invalid engine version format: %v", err)
						}

						if !constraint.Check(engineVersion) {
							return fmt.Errorf("engine version %s does not satisfy required version %s from xel.json, please upgrade your runtime",
								xShared.EngineVersion, manifest.Engine)
						}
					}
				}
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

			// Convert manifest to a VirtLang object
			// Note: This is not reliable because if `manifest` does not exist, it will only use default values
			manifestConverted := map[string]*shared.RuntimeValue{}
			nameVal := values.MK_STRING(fmt.Sprintf("\"%s\"", manifest.Name))
			manifestConverted["name"] = &nameVal
			descVal := values.MK_STRING(fmt.Sprintf("\"%s\"", manifest.Description))
			manifestConverted["description"] = &descVal
			versionVal := values.MK_STRING(fmt.Sprintf("\"%s\"", manifest.Version))
			manifestConverted["version"] = &versionVal
			xelVal := values.MK_STRING(fmt.Sprintf("\"%s\"", manifest.Xel))
			manifestConverted["xel"] = &xelVal
			engineVal := values.MK_STRING(fmt.Sprintf("\"%s\"", manifest.Engine))
			manifestConverted["engine"] = &engineVal
			mainVal := values.MK_STRING(fmt.Sprintf("\"%s\"", manifest.Main))
			manifestConverted["main"] = &mainVal
			authorVal := values.MK_STRING(fmt.Sprintf("\"%s\"", manifest.Author))
			manifestConverted["author"] = &authorVal
			licenseVal := values.MK_STRING(fmt.Sprintf("\"%s\"", manifest.License))
			manifestConverted["license"] = &licenseVal
			// Convert deps map separately
			depsObj := values.MK_OBJECT(map[string]*shared.RuntimeValue{})
			manifestConverted["deps"] = &depsObj
			depsMap := make(map[string]*shared.RuntimeValue)
			for k, v := range manifest.Deps {
				depVal := values.MK_STRING(fmt.Sprintf("\"%s\"", v))
				depsMap[k] = &depVal
			}
			(*manifestConverted["deps"]).Value = depsMap

			manifestObj := values.MK_OBJECT(manifestConverted)
			// Read file content
			content, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %v", filename, err)
			}

			program, parseErr := parser.New(filename).ProduceAST(string(content))
			if parseErr != nil {
				return parseErr
			}

			rootEnv := xShared.XelRootEnv
			RV_proc := map[string]*shared.RuntimeValue{}

			// Check if proc exists
			if _, err := rootEnv.LookupVar("proc"); err != nil {
				// Proc doesn't exist, create new one
				rootEnv.DeclareVar("proc", values.MK_OBJECT(RV_proc), false)
			} else {
				// Proc exists, get its value
				val, err := rootEnv.LookupVar("proc")
				if err != nil {
					return fmt.Errorf("failed to lookup proc variable: %v", err)
				}
				RV_proc = val.Value.(map[string]*shared.RuntimeValue)
			}

			RV_proc_args := values.MK_ARRAY(args)
			RV_proc["args"] = &RV_proc_args
			RV_proc["manifest"] = &manifestObj

			// Update the proc variable
			rootEnv.AssignVar("proc", values.MK_OBJECT(RV_proc))

			env := environment.NewEnvironment(rootEnv)
			env.DeclareVar("__filename__", values.MK_STRING(fmt.Sprintf("\"%s\"", filename)), true)
			env.DeclareVar("__dirname__", values.MK_STRING(fmt.Sprintf("\"%s\"", filepath.Dir(filename))), true)

			// We first have to pause the debugger or
			// else it will run the program without stopping
			xShared.XelRootDebugger.Pause()

			// We start this REPL as a goroutine so that
			// it doesnt block the main thread, and we can
			// also run it in parallel with the evaluation
			go debug_repl(filename, cwd)

			_, evalErr := evaluator.Evaluate(program, env, xShared.XelRootDebugger)
			if evalErr != nil {
				// show the stack trace
				stackTrace := xShared.XelRootDebugger.Snapshots[0]
				stackTraceStr := helpers.GenerateStackTrace(stackTrace.Stack, cwd)
				xShared.ColorPalette.Error.Println(stackTraceStr)
				return evalErr
			}

			return nil
		},
	}
}

func debug_repl(file, cwd string) {
	// We first save our debugger in a variable so that
	// we wont have to type the long thingy every time
	dbgr := xShared.XelRootDebugger

	// First thing we will do here
	// is stepping into the program node
	dbgr.StepInto()

	// This is the metadata for a command
	type Command struct {
		Name        string
		Description string
		Alias       string
		Usage       []string
		Execute     func(args ...string)
	}

	commands := []Command{
		{
			Name:        "quit",
			Description: "Exit the debugger",
			Alias:       "q",
			Usage:       []string{"quit"},
			Execute: func(args ...string) {
				os.Exit(0)
			},
		},
		{
			Name:        "continue",
			Description: "Resume execution until next breakpoint or end",
			Alias:       "c",
			Usage:       []string{"continue"},
			Execute: func(args ...string) {
				dbgr.Continue()
			},
		},
		{
			Name:        "step",
			Description: "Step into the next statement",
			Alias:       "s",
			Usage:       []string{"step"},
			Execute: func(args ...string) {
				dbgr.StepInto()
			},
		},
		{
			Name:        "next",
			Description: "Step over the current statement",
			Alias:       "n",
			Usage:       []string{"next"},
			Execute: func(args ...string) {
				dbgr.StepOver()
			},
		},
		{
			Name:        "out",
			Description: "Step out of the current function",
			Alias:       "o",
			Usage:       []string{"out"},
			Execute: func(args ...string) {
				dbgr.StepOut()
			},
		},
		{
			Name:        "breakpoint",
			Description: "Set or list breakpoints",
			Alias:       "bp",
			Usage: []string{
				"breakpoint <line_number>",
				"breakpoint <filename>:<line_number>",
				"breakpoint <filename> <line_number>",
				"breakpoint list",
				"breakpoint ls",
			},
			Execute: func(args ...string) {
				if len(args) == 0 {
					xShared.ColorPalette.Error.Println("See `help breakpoint` for usage.")
					return
				}

				// List all breakpoints
				if args[0] == "list" || args[0] == "ls" {
					bps := dbgr.BreakpointManager.Breakpoints
					if len(bps) == 0 {
						xShared.ColorPalette.Info.Println("No breakpoints set.")
						return
					}
					xShared.ColorPalette.Info.Println("Breakpoints:")
					i := 0
					for key := range bps {
						splitted := strings.Split(key, ":")
						filename := splitted[0]
						line := splitted[1]
						fmt.Printf("  [%d] %s:%s\n", i+1, filename, line)
						i++
					}
					return
				}

				// Set breakpoint: <line> or <filename>:<line>
				if len(args) == 1 {
					parts := strings.Split(args[0], ":")
					switch len(parts) {
					case 1:
						line, err := strconv.Atoi(parts[0])
						if err != nil {
							xShared.ColorPalette.Error.Println("Invalid line number:", parts[0])
							return
						}
						dbgr.BreakpointManager.Set(file, line)
					case 2:
						line, err := strconv.Atoi(parts[1])
						if err != nil {
							xShared.ColorPalette.Error.Println("Invalid line number:", parts[1])
							return
						}
						dbgr.BreakpointManager.Set(parts[0], line)
					default:
						xShared.ColorPalette.Error.Println("Invalid format. Use <line_number> or <filename>:<line_number>")
					}
					return
				}

				// Set breakpoint: <filename> <line>
				if len(args) == 2 {
					line, err := strconv.Atoi(args[1])
					if err != nil {
						xShared.ColorPalette.Error.Println("Invalid line number:", args[1])
						return
					}
					dbgr.BreakpointManager.Set(args[0], line)
					return
				}

				xShared.ColorPalette.Error.Println("Too many arguments. See `help breakpoint` for usage.")
			},
		},
		{
			Name:        "stacktrace",
			Description: "Display the current call stack",
			Alias:       "st",
			Usage:       []string{"stacktrace"},
			Execute: func(args ...string) {
				stack := dbgr.CallStack
				stacktraceStr := helpers.GenerateStackTrace(stack, cwd)
				xShared.ColorPalette.Info.Println(stacktraceStr)
			},
		},
		{
			Name:        "eval",
			Description: "Evaluate an expression and print the result",
			Alias:       "e",
			Usage:       []string{"eval <expression>"},
			Execute: func(args ...string) {
				if len(args) == 0 {
					xShared.ColorPalette.Error.Println("See `help print` for usage.")
					return
				}

				expr := strings.Join(args, " ")
				stmt, perr := parser.New(file).ProduceAST(expr)
				if perr != nil {
					xShared.ColorPalette.Error.Println("Parser error:", perr)
					return
				}
				res, eerr := evaluator.Evaluate(stmt, dbgr.Environment, nil)
				if eerr != nil {
					xShared.ColorPalette.Error.Println("Evaluation error:", eerr)
					return
				}
				xShared.ColorPalette.GrayMessage.Println(helpers.Stringify(*res, true))
			},
		},
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "(debug) > ",
		HistoryFile:     filepath.Join(os.TempDir(), "xel_debugger_history.tmp"),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		AutoComplete: readline.NewPrefixCompleter(
			readline.PcItem("continue"),
			readline.PcItem("quit"),
			readline.PcItem("step"),
			readline.PcItem("next"),
			readline.PcItem("out"),
			readline.PcItem("breakpoint"),
			readline.PcItem("backtrace"),
			readline.PcItem("help"),
			readline.PcItem("eval"),
		),
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break // ctrl+C twice
			}
			continue
		} else if err == io.EOF {
			break
		}

		input := strings.TrimSpace(line)
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]
		args := parts[1:]

		// help command
		if command == "help" || command == "h" {
			if len(args) > 0 {
				for _, arg := range args {
					for _, cmd := range commands {
						if cmd.Name == arg || cmd.Alias == arg {
							fmt.Printf("  %-12s (%s)\n", cmd.Name, cmd.Alias)
							fmt.Printf("    %s\n", cmd.Description)
							fmt.Println("    Usage:")
							for _, usage := range cmd.Usage {
								fmt.Printf("      %s\n", usage)
							}
							fmt.Println()
							break
						}
					}
				}
				continue
			}

			fmt.Println("Available commands:")
			for _, cmd := range commands {
				fmt.Printf("  %-12s (%s)\n", cmd.Name, cmd.Alias)
				fmt.Printf("    %s\n", cmd.Description)
				fmt.Println("    Usage:")
				for _, usage := range cmd.Usage {
					fmt.Printf("      %s\n", usage)
				}
				fmt.Println()
			}
			continue
		}

		found := false
		for _, cmd := range commands {
			if cmd.Name == command || cmd.Alias == command {
				cmd.Execute(args...)
				found = true
				break
			}
		}

		if !found {
			xShared.ColorPalette.Error.Printf("Unknown command: '%s'\n", command)
			fmt.Println("Type 'help' to see available commands.")
		}
	}
}
