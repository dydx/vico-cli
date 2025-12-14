package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dydx/vico-cli/pkg/auth"
	"github.com/spf13/cobra"
)

// Request represents the JSON request body sent to the Vicohome API
// when fetching events within a specific time range.
type Request struct {
	StartTimestamp string `json:"startTimestamp"` // Start time in Unix timestamp format
	EndTimestamp   string `json:"endTimestamp"`   // End time in Unix timestamp format
	Language       string `json:"language"`       // Language code (e.g., "en" for English)
	CountryNo      string `json:"countryNo"`      // Country code (e.g., "US" for United States)
}

// Event represents a Vicohome event with its properties as returned by the API.
// This structure contains information about bird sightings, including metadata
// about the device that captured the event, the bird identified, and media URLs.
type Event struct {
	TraceID        string  `json:"traceId"`
	Timestamp      string  `json:"timestamp"`
	DeviceName     string  `json:"deviceName"`
	SerialNumber   string  `json:"serialNumber"`
	AdminName      string  `json:"adminName"`
	Period         string  `json:"period"`
	BirdName       string  `json:"birdName"`
	BirdLatin      string  `json:"birdLatin"`
	BirdConfidence float64 `json:"birdConfidence"`
	KeyShotURL     string  `json:"keyShotUrl"`
	ImageURL       string  `json:"imageUrl"`
	VideoURL       string  `json:"videoUrl"`

	// Internal field - not exported to JSON
	keyshots []map[string]interface{} `json:"-"`
}

var (
	startTime    string
	endTime      string
	outputFormat string
)

// listCmd represents the command to list events from the Vicohome API.
// It allows users to fetch events within a specified time range,
// and supports output in both table and JSON formats.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List events within a specified time range",
	Long: `Fetch and display events from Vicohome API for the specified time period. 
Times should be in format: 2025-05-18 14:59:25`,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse and validate time parameters
		start, end, err := parseTimeParameters(startTime, endTime)
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

		events, err := fetchEvents(token, eventsReq)
		if err != nil {
			fmt.Printf("Error fetching events: %v\n", err)
			return
		}

		// Display events
		if len(events) == 0 {
			fmt.Println("No events found in the specified time period.")
			return
		}

		// Write to stdout
		if outputFormat == "json" {
			// Output JSON format
			prettyJSON, err := json.MarshalIndent(events, "", "  ")
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
			for _, event := range events {
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

// supportedTimeFormats contains the timestamp formats that the handler can parse
var supportedTimeFormats = []string{
	"2006-01-02 15:04:05", // Standard format
	time.RFC3339,          // ISO 8601 format
}

// parseTimestamp attempts to parse a timestamp string using supported formats
func parseTimestamp(timestamp string) (time.Time, error) {
	var lastErr error

	// Try each supported format
	for _, format := range supportedTimeFormats {
		t, err := time.Parse(format, timestamp)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}

	// If we get here, none of the formats worked
	return time.Time{}, lastErr
}

func init() {
	currentTime := time.Now()
	defaultStart := currentTime.Add(-24 * time.Hour).Format("2006-01-02 15:04:05")
	defaultEnd := currentTime.Format("2006-01-02 15:04:05")

	listCmd.Flags().StringVar(&startTime, "startTime", defaultStart, "Start time (format: 2006-01-02 15:04:05)")
	listCmd.Flags().StringVar(&endTime, "endTime", defaultEnd, "End time (format: 2006-01-02 15:04:05)")
	listCmd.Flags().StringVar(&outputFormat, "format", "table", "Output format (table or json)")
}

// fetchEvents retrieves events from the Vicohome API within the specified time range.
// It takes an authentication token and a Request object containing the time range
// parameters, and returns a slice of Event objects and any error encountered.
// This function handles the API request, response parsing, and error handling.
func fetchEvents(token string, request Request) ([]Event, error) {
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", GetBaseURL()+"/library/newselectlibrary", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", token)

	// Use ExecuteWithRetry for automatic token refresh
	respBody, err := auth.ExecuteWithRetry(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	// Parse response
	var responseMap map[string]interface{}
	if err := json.Unmarshal(respBody, &responseMap); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w\nResponse: %s", err, string(respBody))
	}

	// Check for API errors
	if code, ok := responseMap["code"].(float64); ok && code != 0 {
		msg, _ := responseMap["msg"].(string)
		return nil, fmt.Errorf("API returned error: %s (code: %.0f)", msg, code)
	}

	// Extract the event list
	data, ok := responseMap["data"].(map[string]interface{})
	if !ok {
		return []Event{}, nil
	}

	eventList, ok := data["list"].([]interface{})
	if !ok {
		return []Event{}, nil
	}

	// Transform events to our simpler format
	events := make([]Event, 0, len(eventList))
	for _, item := range eventList {
		if eventMap, ok := item.(map[string]interface{}); ok {
			// Transform the event to our format
			transformedEvent := transformRawEvent(eventMap)
			events = append(events, transformedEvent)
		}
	}

	return events, nil
}

// transformRawEvent converts a map of event data from the API response into an Event struct.
// It safely extracts and type-converts the various event properties from the dynamic map
// into the strongly-typed Event structure. It handles special processing for bird information,
// timestamps, and keyshots. Default values are provided for missing or unidentified fields.
func transformRawEvent(eventMap map[string]interface{}) Event {
	event := Event{}

	// Extract string fields
	if val, ok := eventMap["traceId"].(string); ok {
		event.TraceID = val
	}

	// Fix: Handle timestamp as a number
	if val, ok := eventMap["timestamp"].(float64); ok {
		// Convert Unix timestamp to human-readable format
		t := time.Unix(int64(val), 0)
		event.Timestamp = t.Format("2006-01-02 15:04:05")
	} else if val, ok := eventMap["timestamp"].(string); ok {
		event.Timestamp = val
	}

	if val, ok := eventMap["deviceName"].(string); ok {
		event.DeviceName = val
	}
	if val, ok := eventMap["serialNumber"].(string); ok {
		event.SerialNumber = val
	}
	if val, ok := eventMap["adminName"].(string); ok {
		event.AdminName = val
	}

	// Fix: Handle period as a number
	if val, ok := eventMap["period"].(float64); ok {
		event.Period = fmt.Sprintf("%.2fs", val)
	} else if val, ok := eventMap["period"].(string); ok {
		event.Period = val
	}

	if val, ok := eventMap["imageUrl"].(string); ok {
		event.ImageURL = val
	}
	if val, ok := eventMap["videoUrl"].(string); ok {
		event.VideoURL = val
	}

	// Set default bird name
	event.BirdName = "Unidentified"

	// Process subcategoryInfoList for bird data
	if subcategoryInfoList, ok := eventMap["subcategoryInfoList"].([]interface{}); ok && len(subcategoryInfoList) > 0 {
		for _, info := range subcategoryInfoList {
			if infoMap, ok := info.(map[string]interface{}); ok {
				// Check if this is a bird entry
				if objectType, ok := infoMap["objectType"].(string); ok && objectType == "bird" {
					// Extract bird name
					if birdName, ok := infoMap["objectName"].(string); ok {
						event.BirdName = birdName
					}
					// Extract Latin name
					if birdLatin, ok := infoMap["birdStdName"].(string); ok {
						event.BirdLatin = birdLatin
					}
					// Extract confidence
					if confidence, ok := infoMap["confidence"].(float64); ok {
						event.BirdConfidence = confidence
					}
					break
				}
			}
		}
	}

	// Handle the keyshots field separately
	if keyshots, ok := eventMap["keyshots"].([]interface{}); ok {
		transformedKeyshots := make([]map[string]interface{}, 0, len(keyshots))
		for _, ks := range keyshots {
			if ksMap, ok := ks.(map[string]interface{}); ok {
				// Create new keyshot with just the desired fields
				newKeyshot := make(map[string]interface{})
				// Copy needed fields
				if url, ok := ksMap["imageUrl"].(string); ok {
					newKeyshot["imageUrl"] = url
					// Extract the first keyshot URL for the flat structure
					if event.KeyShotURL == "" {
						event.KeyShotURL = url
					}
				}
				if msg, ok := ksMap["message"].(string); ok {
					newKeyshot["message"] = msg
				}
				if cat, ok := ksMap["objectCategory"].(string); ok {
					newKeyshot["objectCategory"] = cat
				}
				if sub, ok := ksMap["subCategoryName"].(string); ok {
					newKeyshot["subCategoryName"] = sub
				}
				transformedKeyshots = append(transformedKeyshots, newKeyshot)
			}
		}
		event.keyshots = transformedKeyshots
	} else {
		event.keyshots = []map[string]interface{}{}
	}

	return event
}
