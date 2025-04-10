package events

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dydx/vico-cli/pkg/auth"
	"github.com/spf13/cobra"
)

var (
	searchField string
	searchTerm  string
	searchHours int
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search events by field value",
	Long:  `Search for events that match a specific field value.`,
	Run: func(cmd *cobra.Command, args []string) {
		if searchField == "" {
			fmt.Println("Error: --field flag is required")
			cmd.Help()
			return
		}

		if searchTerm == "" && len(args) > 0 {
			searchTerm = args[0]
		}

		if searchTerm == "" {
			fmt.Println("Error: search term is required")
			cmd.Help()
			return
		}

		token, err := auth.Authenticate()
		if err != nil {
			fmt.Printf("Authentication failed: %v\n", err)
			return
		}

		// Use the same logic as list command to get events from the last searchHours
		end := time.Now()
		start := end.Add(-time.Duration(searchHours) * time.Hour)

		startTimestamp := fmt.Sprintf("%d", start.Unix())
		endTimestamp := fmt.Sprintf("%d", end.Unix())

		eventsReq := EventsRequest{
			StartTimestamp: startTimestamp,
			EndTimestamp:   endTimestamp,
			Language:       "en",
			CountryNo:      "US",
		}

		allEvents, err := fetchEvents(token, eventsReq)
		if err != nil {
			fmt.Printf("Error fetching events: %v\n", err)
			return
		}

		// Filter events based on search field and term
		var filteredEvents []Event
		for _, event := range allEvents {
			if matchesSearch(event, searchField, searchTerm) {
				filteredEvents = append(filteredEvents, event)
			}
		}

		// Display filtered events
		if len(filteredEvents) == 0 {
			fmt.Printf("No events found matching %s = '%s'\n", searchField, searchTerm)
			return
		}

		if outputFormat == "json" {
			// Output JSON format
			prettyJSON, err := json.MarshalIndent(filteredEvents, "", "  ")
			if err != nil {
				fmt.Printf("Error formatting JSON: %v\n", err)
				return
			}
			fmt.Println(string(prettyJSON))
		} else {
			// Output table format
			fmt.Printf("%-36s %-20s %-25s %-25s %-25s\n", 
				"Trace ID", "Timestamp", "Device Name", "Bird Name", "Bird Latin")
			fmt.Println("--------------------------------------------------------------------------------------------------")
			for _, event := range filteredEvents {
				fmt.Printf("%-36s %-20s %-25s %-25s %-25s\n", 
					event.TraceId, 
					event.Timestamp,
					event.DeviceName,
					event.BirdName,
					event.BirdLatin)
			}
		}
	},
}

func init() {
	searchCmd.Flags().StringVar(&searchField, "field", "", "Field to search (serialNumber, deviceName, birdName)")
	searchCmd.Flags().StringVar(&searchTerm, "value", "", "Value to search for")
	searchCmd.Flags().IntVar(&searchHours, "hours", 24, "Number of hours to search back for events")
	searchCmd.Flags().StringVar(&outputFormat, "format", "table", "Output format (table or json)")
	
	// Mark the field flag as required
	searchCmd.MarkFlagRequired("field")
}

// matchesSearch checks if an event matches the search criteria
func matchesSearch(event Event, field, term string) bool {
	term = strings.ToLower(term)
	
	switch strings.ToLower(field) {
	case "serialnumber":
		return strings.ToLower(event.SerialNumber) == term
	case "devicename":
		return strings.Contains(strings.ToLower(event.DeviceName), term)
	case "birdname":
		return strings.Contains(strings.ToLower(event.BirdName), term)
	default:
		// If the field isn't recognized, return false
		return false
	}
}