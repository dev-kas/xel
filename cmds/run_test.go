package cmds

import (
	"flag"
	"os"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestRunCommand(t *testing.T) {
	// Create a temporary test file with integer result
	tempFile, err := os.CreateTemp("", "test_*.xel")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test content to the file
	testContent := "fn add(a, b) { return a + b } add(10, 20)"
	if _, err := tempFile.Write([]byte(testContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Create a temporary test file with string result
	tempFileString, err := os.CreateTemp("", "test_string_*.xel")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFileString.Name())

	// Write test content with string result
	stringContent := "\"Hello World\""
	if _, err := tempFileString.Write([]byte(stringContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFileString.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Create a temporary test file with null result
	tempFileNull, err := os.CreateTemp("", "test_null_*.xel")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFileNull.Name())

	// Write test content with null result (empty function that returns nothing)
	nullContent := "fn empty() {} empty()"
	if _, err := tempFileNull.Write([]byte(nullContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFileNull.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Create a temporary test file with syntax error
	tempFileError, err := os.CreateTemp("", "test_error_*.xel")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFileError.Name())

	// Write test content with syntax error
	errorContent := "10 + * 20"
	if _, err := tempFileError.Write([]byte(errorContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFileError.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Test cases
	tests := []struct {
		name      string
		args      []string
		wantError bool
	}{
		{
			name:      "no arguments",
			args:      []string{},
			wantError: true,
		},
		{
			name:      "file with wrong extension",
			args:      []string{"test.txt"},
			wantError: true,
		},
		{
			name:      "non-existent file",
			args:      []string{"nonexistent.xel"},
			wantError: true,
		},
		{
			name:      "valid file with integer result",
			args:      []string{tempFile.Name()},
			wantError: false,
		},
		{
			name:      "valid file with string result",
			args:      []string{tempFileString.Name()},
			wantError: false,
		},
		{
			name:      "valid file with null result",
			args:      []string{tempFileNull.Name()},
			wantError: false,
		},
		{
			name:      "file with syntax error",
			args:      []string{tempFileError.Name()},
			wantError: true,
		},
		{
			name:      "valid file with command line arguments",
			args:      []string{tempFile.Name(), "arg1", "arg2"},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new CLI context
			set := flag.NewFlagSet("test", 0)
			ctx := cli.NewContext(nil, set, nil)

			// Set the arguments
			if len(tt.args) > 0 {
				set.Parse(tt.args)
			}

			// Get the run command
			cmd := RunCommand()

			// Execute the command
			err := cmd.Action(ctx)

			// Check if error matches expectation
			if (err != nil) != tt.wantError {
				t.Errorf("RunCommand() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
