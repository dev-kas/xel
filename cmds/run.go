package cmds

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
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

			_, evalErr := evaluator.Evaluate(program, &env)

			if evalErr != nil {
				return evalErr
			}

			return nil
		},
	}
}
