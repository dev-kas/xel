package cmds

import "github.com/urfave/cli/v2"

// GetCommands returns all available commands
func GetCommands() []*cli.Command {
	return []*cli.Command{
		RunCommand(),
		InitCommand(),
		InstallCommand(),
	}
}
