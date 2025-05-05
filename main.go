package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"xel/cmds"
	"xel/globals"
	"xel/helpers"

	"github.com/chzyer/readline"
	"github.com/dev-kas/VirtLang-Go/environment"
	"github.com/dev-kas/VirtLang-Go/evaluator"
	"github.com/dev-kas/VirtLang-Go/parser"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

// Version is the current version of Xel
// This will be overridden during build by using ldflags
var Version = "dev"

func main() {
	app := &cli.App{
		Name:     "xel",
		Usage:    "A runtime for VirtLang",
		Version:  Version,
		Commands: cmds.GetCommands(),
		Action: func(c *cli.Context) error {
			color.Cyan("Welcome to Xel v%s REPL!", Version)
			color.RGB(105, 105, 105).Println("Type '!exit' to exit the REPL.")

			rl, err := readline.NewEx(&readline.Config{
				Prompt:            color.BlueString("> "),
				HistoryFile:       filepath.Join(os.TempDir(), "xel_history.tmp"),
				InterruptPrompt:   "^C",
				EOFPrompt:         "!exit",
				HistorySearchFold: true,
			})
			if err != nil {
				color.Red("Error initializing readline: %s", err.Error())
				return err
			}
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT)
			defer rl.Close()
			env := environment.NewEnvironment(nil)
			globals.Globalize(&env)
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
						color.RGB(105, 105, 105).Println("TODO: Stop current execution")
						continue
					}
				case err := <-errChan:
					if err.Error() == "EOF" {
						color.New(color.FgHiRed).Println("Exiting REPL.")
						return nil
					}
					if err.Error() == "Interrupt" {
						color.RGB(105, 105, 105).Println("TODO: Stop current execution")
						continue
					}
					color.Red("Error reading line: %s", err.Error())
				case line := <-inputChan:
					if line == "!exit" {
						color.New(color.FgHiRed).Println("Exiting REPL.")
						return nil
					}
					if line == "" {
						continue
					}

					p := parser.New()
					program, err := p.ProduceAST(line)
					if err != nil {
						color.Red("Error: %s", err.Error())
						continue
					}

					output, oerr := evaluator.Evaluate(program, &env)
					if oerr != nil {
						color.Red("Error: %s", oerr.Error())
						continue
					}
					if output != nil {
						color.RGB(105, 105, 105).Println(fmt.Sprintf("< %v", helpers.Stringify(*output, true)))
					}
				}
			}
		},
	}

	homedir, err := os.UserHomeDir()
	if err == nil {
		os.MkdirAll(filepath.Join(homedir, ".xel"), os.ModePerm)
		if Version != "dev" {
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
					if latestVersion != Version {
						_ = os.WriteFile(versionFile, []byte(latestVersion), os.ModePerm)
						fmt.Println(color.YellowString("--------------------------------------------------------------------------------------"))
						fmt.Print(color.YellowString("A new version of Xel is available! "))
						fmt.Println(color.RedString(Version), color.YellowString("->"), color.GreenString(latestVersion))
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
		color.Red("Error: %s", err.Error())
		os.Exit(1)
	}
}

func checkNewVersion() string {
	if Version == "dev" {
		return Version
	}

	type Release struct {
		TagName string `json:"tag_name"`
	}

	url := "https://api.github.com/repos/dev-kas/xel/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return Version
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Version
	}

	var release Release
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return Version
	}
	tagName := release.TagName
	if len(tagName) > 1 && tagName[0] == 'v' {
		tagName = tagName[1:]
	} else {
		return Version
	}
	if tagName != Version {
		return tagName
	}
	return Version
}
