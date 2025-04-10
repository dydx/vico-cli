package cmd

import (
	"fmt"
	"os"

	"github.com/dydx/vico-cli/cmd/devices"
	"github.com/dydx/vico-cli/cmd/events"
	"github.com/spf13/cobra"
)

// Version is set during build via -ldflags
var Version = "dev"

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "vico-cli",
	Short: "Interact with Vicohome API",
	Long:  `A CLI tool for interacting with the Vicohome API to fetch and manage events.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// versionCmd represents the version command
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
