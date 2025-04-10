package devices

import (
	"github.com/spf13/cobra"
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

// GetDevicesCmd returns the devices command
func GetDevicesCmd() *cobra.Command {
	return devicesCmd
}