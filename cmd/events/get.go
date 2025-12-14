package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dydx/vico-cli/pkg/auth"
	"github.com/spf13/cobra"
)

// EventRequest represents the JSON request body sent to the Vicohome API
// when fetching a specific event by its trace ID.
type EventRequest struct {
	TraceID   string `json:"traceId"`   // Unique identifier for the event
	Language  string `json:"language"`  // Language code (e.g., "en" for English)
	CountryNo string `json:"countryNo"` // Country code (e.g., "US" for United States)
}

// getCmd represents the command to retrieve details for a specific event by its trace ID.
// It supports output in both table and JSON formats.
var getCmd = &cobra.Command{
	Use:   "get [traceID]",
	Short: "Get details for a specific event",
	Long:  `Fetch and display detailed information for a specific Vicohome event by its trace ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		traceID := args[0]

		token, err := auth.Authenticate()
		if err != nil {
			fmt.Printf("Authentication failed: %v\n", err)
			return
		}

		event, err := getEvent(token, traceID)
		if err != nil {
			fmt.Printf("Error fetching event: %v\n", err)
			return
		}

		// Display event details
		if outputFormat == "json" {
			// Output JSON format
			prettyJSON, err := json.MarshalIndent(event, "", "  ")
			if err != nil {
				fmt.Printf("Error formatting JSON: %v\n", err)
				return
			}
			fmt.Println(string(prettyJSON))
		} else {
			// Output formatted table
			fmt.Println("Event Details:")
			fmt.Println("------------------------------")
			fmt.Printf("Trace ID:       %s\n", event.TraceID)
			fmt.Printf("Timestamp:      %s\n", event.Timestamp)
			fmt.Printf("Device Name:    %s\n", event.DeviceName)
			fmt.Printf("Serial Number:  %s\n", event.SerialNumber)
			fmt.Printf("Admin Name:     %s\n", event.AdminName)
			fmt.Printf("Period:         %s\n", event.Period)
			fmt.Printf("Bird Name:      %s\n", event.BirdName)
			fmt.Printf("Bird Latin:     %s\n", event.BirdLatin)
			if event.BirdConfidence > 0 {
				fmt.Printf("Confidence:     %.2f%%\n", event.BirdConfidence*100)
			}
			fmt.Printf("KeyShot URL:    %s\n", event.KeyShotURL)
			fmt.Printf("Image URL:      %s\n", event.ImageURL)
			fmt.Printf("Video URL:      %s\n", event.VideoURL)
		}
	},
}

func init() {
	getCmd.Flags().StringVar(&outputFormat, "format", "table", "Output format (table or json)")
}

// getEvent fetches detailed information for a specific event from the Vicohome API.
// It takes an authentication token and the event's trace ID, and returns
// an Event object and any error encountered.
// This function handles the API request, response parsing, and error handling.
func getEvent(token string, traceID string) (Event, error) {
	req := EventRequest{
		TraceID:   traceID,
		Language:  "en",
		CountryNo: "US",
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return Event{}, fmt.Errorf("error marshaling request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", GetBaseURL()+"/library/newselectsinglelibrary", bytes.NewBuffer(reqBody))
	if err != nil {
		return Event{}, fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)

	// Use ExecuteWithRetry for automatic token refresh
	respBody, err := auth.ExecuteWithRetry(httpReq)
	if err != nil {
		return Event{}, fmt.Errorf("error making request: %w", err)
	}

	// Parse response
	var responseMap map[string]interface{}
	if err := json.Unmarshal(respBody, &responseMap); err != nil {
		return Event{}, fmt.Errorf("error unmarshaling response: %w\nResponse: %s", err, string(respBody))
	}

	// Check for API errors
	if code, ok := responseMap["code"].(float64); ok && code != 0 {
		msg, _ := responseMap["msg"].(string)
		return Event{}, fmt.Errorf("API returned error: %s (code: %.0f)", msg, code)
	}

	// Extract event data
	data, ok := responseMap["data"].(map[string]interface{})
	if !ok {
		return Event{}, fmt.Errorf("no event data found")
	}

	// First check if data has the traceId field, which indicates it's an event
	if _, hasTraceID := data["traceId"].(string); hasTraceID {
		return transformRawEvent(data), nil
	}

	// If we didn't find the event directly in data, try data.event as a fallback
	event, ok := data["event"].(map[string]interface{})
	if !ok {
		return Event{}, fmt.Errorf("no event data found")
	}

	return transformRawEvent(event), nil
}
