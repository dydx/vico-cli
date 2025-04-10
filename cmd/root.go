// Package cmd implements the command structure for the Vicohome CLI application.
//
// This package uses the Cobra library to define commands, subcommands and flags
// that make up the CLI's interface. It acts as the command router, directing user
// input to the appropriate handlers.
package cmd

import (
	"fmt"
	"os"

	"github.com/dydx/vico-cli/cmd/devices"
	"github.com/dydx/vico-cli/cmd/events"
	"github.com/spf13/cobra"
)

// Version is set during build via -ldflags.
// It represents the current version of the CLI application.
var Version = "dev"

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "vico-cli",
	Short: "Interact with Vicohome API",
	Long:  `A CLI tool for interacting with the Vicohome API to fetch and manage events.`,
}

// Execute runs the root command and handles any resulting errors.
// If the command execution fails, the program will exit with a non-zero status code.
// This function is called by the main function and serves as the entry point for the CLI.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// versionCmd represents the version command, which displays the current version of the CLI.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version of vicohome",
	Long:  `Display the version of the vicohome CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("vicohome version %s\n", Version)
	},
}

func init() {
	// No persistent flags needed

	// Add the commands
	rootCmd.AddCommand(devices.GetDevicesCmd())
	rootCmd.AddCommand(events.GetEventsCmd())
	rootCmd.AddCommand(versionCmd)
}
