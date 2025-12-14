// Package devices implements commands for managing and querying Vicohome devices.
//
// This package provides commands for listing all devices associated with a user account
// and retrieving detailed information about specific devices.
package devices

import (
	"github.com/spf13/cobra"
	"os"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Manage Vicohome devices",
	Long:  `List and get details for Vicohome devices.`,
}

func init() {
	// Add subcommands
	devicesCmd.AddCommand(listCmd)
	devicesCmd.AddCommand(getCmd)
}
func GetBaseURL() string {
if v := os.Getenv("VICOHOME_BASE_URL"); v != "" {
return v
}
return "https://api-us.vicohome.io"
}

// GetDevicesCmd returns the devices command that provides access to device-related subcommands.
// This function is called by the root command to add device functionality to the CLI.
// It returns the devices command with all subcommands (list, get) already attached.
func GetDevicesCmd() *cobra.Command {
	return devicesCmd
}
