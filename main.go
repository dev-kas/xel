package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"xel/cmds"
	"xel/globals"
	"xel/helpers"
	"xel/shared"

	_ "xel/modules/array"
	_ "xel/modules/math"
	_ "xel/modules/strings"
	_ "xel/modules/threads"
	_ "xel/modules/time"
	_ "xel/modules/native"
	_ "xel/modules/os"
	_ "xel/modules/classes"
	_ "xel/modules/object"

	"github.com/chzyer/readline"
	"github.com/dev-kas/virtlang-go/v4/environment"
	"github.com/dev-kas/virtlang-go/v4/evaluator"
	"github.com/dev-kas/virtlang-go/v4/parser"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func main() {
	globals.Globalize(shared.XelRootEnv)
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("VirtLang Engine version: %s - Xel version: %s \n", shared.ColorPalette.Version(shared.EngineVersion), shared.ColorPalette.Version(c.App.Version))
	}

	app := &cli.App{
		Name:     "xel",
		Usage:    "A runtime for VirtLang",
		Version:  shared.RuntimeVersion,
		Commands: cmds.GetCommands(),
		Action: func(c *cli.Context) error {
			shared.ColorPalette.Welcome("Welcome to Xel v%s REPL (VirtLang v%s)!", shared.RuntimeVersion, shared.EngineVersion)
			shared.ColorPalette.GrayMessage.Println("Type '!exit' to exit the REPL.")

			rl, err := readline.NewEx(&readline.Config{
				Prompt:            shared.ColorPalette.PromptStr("> "),
				HistoryFile:       filepath.Join(os.TempDir(), "xel_history.tmp"),
				InterruptPrompt:   "^C",
				EOFPrompt:         "!exit",
				HistorySearchFold: true,
			})
			if err != nil {
				shared.ColorPalette.Error.Printf("Error initializing readline: %s", err.Error())
				return err
			}
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT)
			defer rl.Close()
			env := environment.NewEnvironment(shared.XelRootEnv)
			for {
				inputChan := make(chan string, 1)
				errChan := make(chan error, 1)

				go func() {
					input, err := rl.Readline()
					if err != nil {
						errChan <- err
					} else {
						inputChan <- input
					}
				}()

				select {
				case sig := <-sigChan:
					if sig == syscall.SIGINT {
						shared.ColorPalette.GrayMessage.Print("TODO: Stop current execution")
						continue
					}
				case err := <-errChan:
					if err.Error() == "EOF" {
						shared.ColorPalette.ExitMessage.Println("Exiting REPL.")
						return nil
					}
					if err.Error() == "Interrupt" {
						shared.ColorPalette.GrayMessage.Print("TODO: Stop current execution")
						continue
					}
					shared.ColorPalette.Error.Printf("Error reading line: %s", err.Error())
				case line := <-inputChan:
					if line == "!exit" {
						shared.ColorPalette.ExitMessage.Println("Exiting REPL.")
						return nil
					}
					if line == "" {
						continue
					}

					p := parser.New("<REPL>")
					program, err := p.ProduceAST(line)
					if err != nil {
						shared.ColorPalette.Error.Printf("Error: %s\n", err.Error())
						continue
					}

					output, oerr := evaluator.Evaluate(program, env, shared.XelRootDebugger)
					if oerr != nil && len(shared.XelRootDebugger.Snapshots) > 0 {
						stackTrace := shared.XelRootDebugger.Snapshots[0]
						stackTraceStr := helpers.GenerateStackTrace(stackTrace.Stack, "")
						shared.ColorPalette.Error.Println(stackTraceStr)
						shared.ColorPalette.Error.Printf("Error: %s\n", oerr.Error())
						shared.XelRootDebugger.Snapshots = nil

						continue
					}
					if output != nil {
						outputStr := helpers.Stringify(*output, true)
						lines := strings.Split(outputStr, "\n")
						for i, line := range lines {
							if i == 0 {
								shared.ColorPalette.GrayMessage.Printf("< %s", line)
							} else {
								shared.ColorPalette.GrayMessage.Printf("  %s", line)
							}
							if i < len(lines)-1 {
								fmt.Println()
							}
						}
						fmt.Println()
					}
				}
			}
		},
	}

	homedir, err := os.UserHomeDir()
	if err == nil {
		os.MkdirAll(filepath.Join(homedir, ".xel"), os.ModePerm)
		if shared.RuntimeVersion != "" {
			go func() {
				versionFile := filepath.Join(homedir, ".xel", "version-latest")
				needCheck := true

				if info, err := os.Stat(versionFile); err == nil {
					if time.Since(info.ModTime()) < 24*time.Hour {
						needCheck = false
					}
				}

				if needCheck {
					latestVersion := checkNewVersion()
					if latestVersion != shared.RuntimeVersion {
						_ = os.WriteFile(versionFile, []byte(latestVersion), os.ModePerm)
						fmt.Println(color.YellowString("--------------------------------------------------------------------------------------"))
						fmt.Print(color.YellowString("A new version of Xel is available! "))
						fmt.Println(color.RedString(shared.RuntimeVersion), color.YellowString("->"), color.GreenString(latestVersion))
						fmt.Println(color.YellowString("To update, run:"))
						fmt.Println(color.YellowString("curl -fsSL https://raw.githubusercontent.com/dev-kas/xel/master/scripts/update.sh | sh"))
						fmt.Println(color.YellowString("--------------------------------------------------------------------------------------"))
					}
				}
			}()
		}
	}

	err = app.Run(os.Args)
	if err != nil {
		shared.ColorPalette.Error.Println(err.Error())
		os.Exit(1)
	}
}

func checkNewVersion() string {
	if shared.RuntimeVersion == "" {
		return shared.RuntimeVersion
	}

	type Release struct {
		TagName string `json:"tag_name"`
	}

	url := "https://api.github.com/repos/dev-kas/xel/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return shared.RuntimeVersion
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return shared.RuntimeVersion
	}

	var release Release
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return shared.RuntimeVersion
	}
	tagName := release.TagName
	if len(tagName) > 1 && tagName[0] == 'v' {
		tagName = tagName[1:]
	} else {
		return shared.RuntimeVersion
	}
	if tagName != shared.RuntimeVersion {
		return tagName
	}
	return shared.RuntimeVersion
}
