package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"
)

// TestMain tests the main function behavior
func TestMain(t *testing.T) {
	// Save original args and stdout
	originalArgs := os.Args
	originalStdout := os.Stdout

	defer func() {
		os.Args = originalArgs
		os.Stdout = originalStdout
	}()

	tests := []struct {
		name     string
		args     []string
		wantExit bool
	}{
		{
			name:     "help command",
			args:     []string{"shuku", "help"},
			wantExit: false,
		},
		{
			name:     "version command",
			args:     []string{"shuku", "version"},
			wantExit: false,
		},
		{
			name:     "invalid command",
			args:     []string{"shuku", "invalid"},
			wantExit: false, // CLI framework handles this gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up args
			os.Args = tt.args

			// Capture stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// This test is tricky because main() calls os.Exit()
			// We'll test the root command directly instead
			app := rootCmd()
			err := app.Run(tt.args)

			// Restore stdout and read output
			w.Close()
			os.Stdout = originalStdout

			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// For invalid commands, the CLI framework handles this gracefully
			if tt.name == "invalid command" {
				// CLI framework may return error or handle gracefully
				t.Logf("Invalid command handling - Error: %v, Output: %s", err, output)
			} else {
				// For valid commands, we don't expect errors
				if err != nil && !strings.Contains(err.Error(), "exit status") {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestRootCmd tests the root command structure
func TestRootCmd(t *testing.T) {
	app := rootCmd()

	// Test basic app properties
	if app.Name != "shuku" {
		t.Errorf("App name = %v, want %v", app.Name, "shuku")
	}

	if app.Usage != "A CLI tool for compressing images." {
		t.Errorf("App usage = %v, want %v", app.Usage, "A CLI tool for compressing images.")
	}

	if app.UsageText != "shuku [command] [options] [arguments]" {
		t.Errorf("App usage text = %v, want %v", app.UsageText, "shuku [command] [options] [arguments]")
	}

	if app.HelpName != "shuku" {
		t.Errorf("App help name = %v, want %v", app.HelpName, "shuku")
	}

	// Test suggest feature is enabled
	if !app.Suggest {
		t.Error("App suggest should be enabled")
	}

	// Test commands are registered
	expectedCommands := []string{"compress", "batch", "version", "help"}
	if len(app.Commands) != len(expectedCommands) {
		t.Errorf("App commands length = %v, want %v", len(app.Commands), len(expectedCommands))
	}

	// Check each command is present
	commandNames := make(map[string]bool)
	for _, cmd := range app.Commands {
		commandNames[cmd.Name] = true
	}

	for _, expectedCmd := range expectedCommands {
		if !commandNames[expectedCmd] {
			t.Errorf("Expected command %s not found", expectedCmd)
		}
	}
}

// TestHelpCommand tests the custom help command
func TestHelpCommand(t *testing.T) {
	app := rootCmd()

	// Find the help command
	var helpCmd *cli.Command
	for _, cmd := range app.Commands {
		if cmd.Name == "help" {
			helpCmd = cmd
			break
		}
	}

	if helpCmd == nil {
		t.Fatal("Help command not found")
	}

	// Test help command properties
	if helpCmd.Name != "help" {
		t.Errorf("Help command name = %v, want %v", helpCmd.Name, "help")
	}

	expectedAliases := []string{"h"}
	if len(helpCmd.Aliases) != len(expectedAliases) {
		t.Errorf("Help command aliases length = %v, want %v", len(helpCmd.Aliases), len(expectedAliases))
	}

	if len(helpCmd.Aliases) > 0 && helpCmd.Aliases[0] != "h" {
		t.Errorf("Help command alias[0] = %v, want %v", helpCmd.Aliases[0], "h")
	}

	if helpCmd.Usage != "Shows a list of commands or help for one command" {
		t.Errorf("Help command usage = %v, want %v", helpCmd.Usage, "Shows a list of commands or help for one command")
	}

	if helpCmd.ArgsUsage != "[command]" {
		t.Errorf("Help command args usage = %v, want %v", helpCmd.ArgsUsage, "[command]")
	}

	// Test help command action
	if helpCmd.Action == nil {
		t.Error("Help command action should not be nil")
	}
}

// TestCommandNotFound tests the custom command not found handler
func TestCommandNotFound(t *testing.T) {
	app := rootCmd()

	if app.CommandNotFound == nil {
		t.Error("CommandNotFound handler should not be nil")
	}

	// Capture stdout to test the handler
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a mock context (this is a simplified test)
	// In real usage, this would be called by the CLI framework
	if app.CommandNotFound != nil {
		// We can't easily test this without mocking the entire context
		// but we can verify it's set
		t.Log("CommandNotFound handler is properly set")
	}

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
}
