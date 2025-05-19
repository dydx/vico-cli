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
	searchField     string
	searchTerm      string
	searchStartTime string
	searchEndTime   string
)

// searchCmd represents the command to search for events that match specific criteria.
// It allows filtering events by field values (such as device name or bird name)
// within a specified time range, and supports output in both table and JSON formats.
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search events by field value",
	Long: `Search for events that match a specific field value within a specified time range.
Times should be in format: 2025-05-18 14:59:25`,
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

		// Parse and validate time parameters
		start, end, err := parseTimeParameters(searchStartTime, searchEndTime)
		if err != nil {
			fmt.Printf("Error parsing time parameters: %v\n", err)
			return
		}

		token, err := auth.Authenticate()
		if err != nil {
			fmt.Printf("Authentication failed: %v\n", err)
			return
		}

		startTimestamp := fmt.Sprintf("%d", start.Unix())
		endTimestamp := fmt.Sprintf("%d", end.Unix())

		eventsReq := Request{
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
					event.TraceID,
					event.Timestamp,
					event.DeviceName,
					event.BirdName,
					event.BirdLatin)
			}
		}
	},
}

func init() {
	currentTime := time.Now()
	defaultStart := currentTime.Add(-24 * time.Hour).Format("2006-01-02 15:04:05")
	defaultEnd := currentTime.Format("2006-01-02 15:04:05")

	searchCmd.Flags().StringVar(&searchField, "field", "", "Field to search (serialNumber, deviceName, birdName)")
	searchCmd.Flags().StringVar(&searchTerm, "value", "", "Value to search for")
	searchCmd.Flags().StringVar(&searchStartTime, "startTime", defaultStart, "Start time (format: 2006-01-02 15:04:05)")
	searchCmd.Flags().StringVar(&searchEndTime, "endTime", defaultEnd, "End time (format: 2006-01-02 15:04:05)")
	searchCmd.Flags().StringVar(&outputFormat, "format", "table", "Output format (table or json)")

	// Mark the field flag as required
	searchCmd.MarkFlagRequired("field")
}

// matchesSearch checks if an event matches the search criteria provided by the user.
// It compares the specified field in the event with the search term, using case-insensitive
// matching. For some fields like deviceName and birdName, it uses substring matching,
// while for serialNumber it requires an exact match.
//
// Parameters:
//   - event: The Event to check
//   - field: The field name to check against (serialNumber, deviceName, birdName)
//   - term: The value to search for
//
// Returns:
//   - true if the event matches the search criteria, false otherwise
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
