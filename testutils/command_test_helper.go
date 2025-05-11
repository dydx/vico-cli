// Package testutils provides testing utilities for the vico-cli application.
package testutils

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ExecuteCommandCapturingOutput runs a command and captures both stdout and stderr.
// It returns the captured stdout, stderr, and any error from the command execution.
func ExecuteCommandCapturingOutput(t *testing.T, cmd *cobra.Command, args ...string) (string, string, error) {
	t.Helper()
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	
	// Save original streams
	oldStdout, oldStderr := os.Stdout, os.Stderr
	
	// Create pipes
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	
	// Replace stdout/stderr
	os.Stdout, os.Stderr = wOut, wErr
	
	// Set up copying to buffers
	outC := make(chan struct{})
	errC := make(chan struct{})
	go func() {
		io.Copy(stdout, rOut)
		close(outC)
	}()
	go func() {
		io.Copy(stderr, rErr)
		close(errC)
	}()
	
	// Set args and execute
	cmd.SetArgs(args)
	err := cmd.Execute()
	
	// Close write ends of pipes to finish reads
	wOut.Close()
	wErr.Close()
	
	// Wait for copying to complete
	<-outC
	<-errC
	
	// Restore original streams
	os.Stdout, os.Stderr = oldStdout, oldStderr
	
	return stdout.String(), stderr.String(), err
}

// ResetCommandFlags resets all flags on a command to their default values.
// This is useful for testing commands with multiple flag combinations.
func ResetCommandFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
		f.Value.Set(f.DefValue)
	})
}

// CreateTestCommand creates a simple test command with a given name and run function.
// This is useful for testing command routing and flag parsing.
func CreateTestCommand(name string, run func(*cobra.Command, []string)) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: "Test command",
		Long:  "Test command for testing",
		Run:   run,
	}
}