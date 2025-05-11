package cmd

import (
	"strings"
	"testing"

	"github.com/dydx/vico-cli/testutils"
	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	// Test default behavior (no args)
	stdout, stderr, err := testutils.ExecuteCommandCapturingOutput(t, rootCmd)
	assert.NoError(t, err)
	assert.Empty(t, stderr)
	// Root command with no args should display help
	assert.Contains(t, stdout, "Usage:")
	assert.Contains(t, stdout, "vico-cli")
	assert.Contains(t, stdout, "Available Commands:")
}

func TestVersionCommand(t *testing.T) {
	// Save original version and restore after test
	originalVersion := Version
	defer func() { Version = originalVersion }()
	
	// Set a test version
	Version = "0.1.0-test"
	
	stdout, stderr, err := testutils.ExecuteCommandCapturingOutput(t, rootCmd, "version")
	assert.NoError(t, err)
	assert.Empty(t, stderr)
	assert.Equal(t, "vicohome version 0.1.0-test\n", stdout)
}

func TestUnknownCommand(t *testing.T) {
	_, stderr, err := testutils.ExecuteCommandCapturingOutput(t, rootCmd, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, stderr, "unknown command")
}

func TestHelpCommand(t *testing.T) {
	// Test explicit help command
	stdout, stderr, err := testutils.ExecuteCommandCapturingOutput(t, rootCmd, "help")
	assert.NoError(t, err)
	assert.Empty(t, stderr)
	
	// Verify help output contains expected sections
	assert.Contains(t, stdout, "Usage:")
	assert.Contains(t, stdout, "vico-cli")
	assert.Contains(t, stdout, "Available Commands:")
	assert.Contains(t, stdout, "devices") // Check for subcommands
	assert.Contains(t, stdout, "events")
	assert.Contains(t, stdout, "version")
	
	// Test --help flag
	helpFlagOut, _, err := testutils.ExecuteCommandCapturingOutput(t, rootCmd, "--help")
	assert.NoError(t, err)
	
	// Help flag should produce the same output as help command
	assert.Equal(t, stdout, helpFlagOut)
}

func TestCommandCompleteness(t *testing.T) {
	// Test that all expected commands are available
	// This helps ensure that if new commands are added, they're properly registered
	stdout, _, _ := testutils.ExecuteCommandCapturingOutput(t, rootCmd, "help")
	
	expectedCommands := []string{
		"devices",
		"events",
		"version",
		"help",
	}
	
	for _, cmd := range expectedCommands {
		assert.True(t, strings.Contains(stdout, cmd), "Expected command '%s' not found in help output", cmd)
	}
}