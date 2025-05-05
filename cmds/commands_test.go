package cmds

import (
	"testing"
)

func TestGetCommands(t *testing.T) {
	commands := GetCommands()

	// Check if we have the expected number of commands
	expectedCount := 3 // run, init, install
	if len(commands) != expectedCount {
		t.Errorf("GetCommands() returned %d commands, expected %d", len(commands), expectedCount)
	}

	// Check if all expected commands are present
	commandNames := make(map[string]bool)
	for _, cmd := range commands {
		commandNames[cmd.Name] = true
	}

	expectedCommands := []string{"run", "init", "install"}
	for _, name := range expectedCommands {
		if !commandNames[name] {
			t.Errorf("GetCommands() is missing the '%s' command", name)
		}
	}
}
