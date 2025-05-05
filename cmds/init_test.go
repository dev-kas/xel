package cmds

import (
	"flag"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestInitCommand(t *testing.T) {
	// Create a new CLI context
	set := flag.NewFlagSet("test", 0)
	ctx := cli.NewContext(nil, set, nil)

	// Get the init command
	cmd := InitCommand()

	// Check command properties
	if cmd.Name != "init" {
		t.Errorf("InitCommand() name = %s, want %s", cmd.Name, "init")
	}

	if cmd.Usage != "Initialize a new project" {
		t.Errorf("InitCommand() usage = %s, want %s", cmd.Usage, "Initialize a new project")
	}

	// Execute the command (should not return an error as it's unimplemented)
	err := cmd.Action(ctx)
	if err != nil {
		t.Errorf("InitCommand() error = %v, want nil", err)
	}
}
