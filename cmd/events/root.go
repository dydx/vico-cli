package events

import (
	"github.com/spf13/cobra"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Manage Vicohome events",
	Long:  `List and get details for Vicohome events.`,
}

func init() {
	// Add subcommands
	eventsCmd.AddCommand(listCmd)
	eventsCmd.AddCommand(getCmd)
	eventsCmd.AddCommand(searchCmd)
}

// GetEventsCmd returns the events command
func GetEventsCmd() *cobra.Command {
	return eventsCmd
}