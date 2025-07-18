package cmds

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dev-kas/xel/helpers"

	xShared "github.com/dev-kas/xel/shared"

	"github.com/Masterminds/semver/v3"

	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/evaluator"
	"github.com/dev-kas/virtlang-go/v4/parser"
	"github.com/dev-kas/virtlang-go/v4/shared"
	"github.com/dev-kas/virtlang-go/v4/values"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

// RunCommand returns the cli.Command for the run command
// It executes a VirtLang script file after performing version compatibility checks
// against the project's manifest (xel.json) if present.
func RunCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Execute a VirtLang script file",
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

			// Check for xel.json in project directory
			manifest, _, err := helpers.FetchManifest(filepath.Dir(filename), filepath.Dir(filename))
			if err != nil {
				// Ignore error if manifest is not found
				if err.Error() != "cannot find manifest" {
					return err
				}
			}
			if manifest != nil {
				// Check Xel version constraint if specified
				if *manifest.Xel != "" {
					// Skip version check in development mode (when version is empty)
					if xShared.RuntimeVersion == "" {
						color.New(color.FgYellow).Printf("Skipping Xel version check in development mode\n")
					} else {
						constraint, err := semver.NewConstraint(*manifest.Xel)
						if err != nil {
							return fmt.Errorf("invalid Xel version constraint in manifest: %v", err)
						}

						runtimeVersion, err := semver.NewVersion(xShared.RuntimeVersion)
						if err != nil {
							return fmt.Errorf("invalid runtime version format: %v", err)
						}

						if !constraint.Check(runtimeVersion) {
							return fmt.Errorf("xel version %s does not satisfy required version %s from xel.json, please upgrade your runtime",
								xShared.RuntimeVersion, *manifest.Xel)
						}
					}
				}

				// Check Engine version constraint if specified
				if *manifest.Engine != "" {
					// Skip version check in development mode (when version is empty)
					if xShared.EngineVersion == "" {
						color.New(color.FgYellow).Printf("Skipping Engine version check in development mode\n")
					} else {
						constraint, err := semver.NewConstraint(*manifest.Engine)
						if err != nil {
							return fmt.Errorf("invalid Engine version constraint in manifest: %v", err)
						}

						engineVersion, err := semver.NewVersion(xShared.EngineVersion)
						if err != nil {
							return fmt.Errorf("invalid engine version format: %v", err)
						}

						if !constraint.Check(engineVersion) {
							return fmt.Errorf("engine version %s does not satisfy required version %s from xel.json, please upgrade your runtime",
								xShared.EngineVersion, *manifest.Engine)
						}
					}
				}
			} else {
				// No manifest found, use default values
				manifest = &xShared.ProjectManifest{
					Name:    "Unknown",
					Version: "0.0.0",
					Deps:    &map[string]string{},
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
				args[i] = values.MK_STRING(arg)
			}

			// Convert manifest to a VirtLang object
			// Note: This is not reliable because if `manifest` does not exist, it will only use default values
			manifestConverted := map[string]*shared.RuntimeValue{}
			nameVal := values.MK_STRING(manifest.Name)
			manifestConverted["name"] = &nameVal
			descVal := values.MK_STRING(manifest.Description)
			manifestConverted["description"] = &descVal
			versionVal := values.MK_STRING(manifest.Version)
			manifestConverted["version"] = &versionVal
			xelVal := values.MK_STRING(*manifest.Xel)
			manifestConverted["xel"] = &xelVal
			engineVal := values.MK_STRING(*manifest.Engine)
			manifestConverted["engine"] = &engineVal
			mainVal := values.MK_STRING(manifest.Main)
			manifestConverted["main"] = &mainVal
			authorVal := values.MK_STRING(manifest.Author)
			manifestConverted["author"] = &authorVal
			licenseVal := values.MK_STRING(manifest.License)
			manifestConverted["license"] = &licenseVal
			// Convert deps map separately
			depsObj := values.MK_OBJECT(map[string]*shared.RuntimeValue{})
			manifestConverted["deps"] = &depsObj
			depsMap := make(map[string]*shared.RuntimeValue)
			for k, v := range *manifest.Deps {
				depVal := values.MK_STRING(v)
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
			env.DeclareVar("__filename__", values.MK_STRING(filename), true)
			env.DeclareVar("__dirname__", values.MK_STRING(filepath.Dir(filename)), true)

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
