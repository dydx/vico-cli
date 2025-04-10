package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dydx/vico-cli/pkg/auth"
	"github.com/spf13/cobra"
)

type EventsRequest struct {
	StartTimestamp string `json:"startTimestamp"`
	EndTimestamp   string `json:"endTimestamp"`
	Language       string `json:"language"`
	CountryNo      string `json:"countryNo"`
}

type Event struct {
	TraceId      string                   `json:"traceId"`
	Timestamp    string                   `json:"timestamp"`
	DeviceName   string                   `json:"deviceName"`
	SerialNumber string                   `json:"serialNumber"`
	AdminName    string                   `json:"adminName"`
	Period       string                   `json:"period"`
	Keyshots     []map[string]interface{} `json:"keyshots"`
	ImageUrl     string                   `json:"imageUrl"`
	VideoUrl     string                   `json:"videoUrl"`
}

var hours int
var outputFormat string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List events from the last N hours",
	Long:  `Fetch and display events from Vicohome API for the specified time period.`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := auth.Authenticate()
		if err != nil {
			fmt.Printf("Authentication failed: %v\n", err)
			return
		}

		end := time.Now()
		start := end.Add(-time.Duration(hours) * time.Hour)

		startTimestamp := fmt.Sprintf("%d", start.Unix())
		endTimestamp := fmt.Sprintf("%d", end.Unix())

		eventsReq := EventsRequest{
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
			fmt.Printf("%-36s %-20s %-30s %-36s\n", 
				"Trace ID", "Timestamp", "Device Name", "Serial Number")
			fmt.Println("------------------------------------------------------------------------------------------------")
			for _, event := range events {
				fmt.Printf("%-36s %-20s %-30s %-36s\n", 
					event.TraceId, 
					event.Timestamp,
					event.DeviceName,
					event.SerialNumber)
			}
		}
	},
}

func init() {
	listCmd.Flags().IntVar(&hours, "hours", 24, "Number of hours to fetch events for")
	listCmd.Flags().StringVar(&outputFormat, "format", "table", "Output format (table or json)")
}

func fetchEvents(token string, request EventsRequest) ([]Event, error) {
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api-us.vicohome.io/library/newselectlibrary", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
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

func transformRawEvent(eventMap map[string]interface{}) Event {
	event := Event{}

	// Extract string fields
	if val, ok := eventMap["traceId"].(string); ok {
		event.TraceId = val
	}
	if val, ok := eventMap["timestamp"].(string); ok {
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
	if val, ok := eventMap["period"].(string); ok {
		event.Period = val
	}
	if val, ok := eventMap["imageUrl"].(string); ok {
		event.ImageUrl = val
	}
	if val, ok := eventMap["videoUrl"].(string); ok {
		event.VideoUrl = val
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
		event.Keyshots = transformedKeyshots
	} else {
		event.Keyshots = []map[string]interface{}{}
	}

	return event
}