// Package events implements commands for retrieving and working with Vicohome events.
//
// This package provides commands for listing events, getting details for a specific event,
// and searching through events based on various criteria.
package events

import (
	"github.com/spf13/cobra"
	"os"
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
func GetBaseURL() string {
if v := os.Getenv("VICOHOME_BASE_URL"); v != "" {
return v
}
return "https://api-us.vicohome.io"
}
// GetEventsCmd returns the events command that provides access to event-related subcommands.
// This function is called by the root command to add event functionality to the CLI.
// It returns the events command with all subcommands (list, get, search) already attached.
func GetEventsCmd() *cobra.Command {
	return eventsCmd
}
