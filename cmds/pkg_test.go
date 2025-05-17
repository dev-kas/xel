package cmds

import (
	"flag"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestInstallCommand(t *testing.T) {
	// Create a new CLI context
	set := flag.NewFlagSet("test", 0)
	ctx := cli.NewContext(nil, set, nil)

	// Get the install command
	cmd := PackageCommands()

	// Check command properties
	if cmd.Name != "pkg" {
		t.Errorf("PackageCommands() name = %s, want %s", cmd.Name, "pkg")
	}

	if cmd.Usage != "Manage packages" {
		t.Errorf("PackageCommands() usage = %s, want %s", cmd.Usage, "Manage packages")
	}

	// Execute the command (should not return an error as it's unimplemented)
	err := cmd.Action(ctx)
	if err != nil {
		t.Errorf("PackageCommands() error = %v, want nil", err)
	}
}
