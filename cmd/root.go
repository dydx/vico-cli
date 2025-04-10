package cmd

import (
	"os"

	"github.com/dydx/vico-cli/cmd/devices"
	"github.com/dydx/vico-cli/cmd/events"
	"github.com/spf13/cobra"
)

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

func init() {
	// No persistent flags needed

	// Add the commands
	rootCmd.AddCommand(devices.GetDevicesCmd())
	rootCmd.AddCommand(events.GetEventsCmd())
}
