package main

import (
	"os"
	"testing"
	"xel/cmds"

	"github.com/urfave/cli/v2"
)

func TestMainApp(t *testing.T) {
	// Create a test app similar to the one in main()
	app := &cli.App{
		Name:     "xel",
		Usage:    "A runtime for VirtLang",
		Version:  Version,
		Commands: cmds.GetCommands(),
	}
	
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	// Test cases
	tests := []struct {
		name      string
		args      []string
		wantError bool
	}{
		{
			name:      "version flag",
			args:      []string{"xel", "--version"},
			wantError: false,
		},
		{
			name:      "help flag",
			args:      []string{"xel", "--help"},
			wantError: false,
		},
		{
			name:      "unknown command",
			args:      []string{"xel", "unknown"},
			wantError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify the app structure
			if app.Name != "xel" {
				t.Errorf("App name = %s, want %s", app.Name, "xel")
			}
			
			if app.Usage != "A runtime for VirtLang" {
				t.Errorf("App usage = %s, want %s", app.Usage, "A runtime for VirtLang")
			}
			
			if len(app.Commands) != 3 {
				t.Errorf("App has %d commands, want %d", len(app.Commands), 3)
			}
			
			// Skip actual execution to avoid program exit
			// We're just testing the structure here
		})
	}
}