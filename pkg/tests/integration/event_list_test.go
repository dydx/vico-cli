package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

// TestEventListingWithPagination tests the event listing functionality with pagination.
func TestEventListingWithPagination(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=true to run.")
	}

	// Create a mock API server
	mockServer := NewMockAPIServer()
	defer mockServer.Close()

	// Override the API endpoint to use our mock server
	originalAPIEndpoint := os.Getenv("VICOHOME_API_ENDPOINT")
	defer func() {
		if originalAPIEndpoint != "" {
			os.Setenv("VICOHOME_API_ENDPOINT", originalAPIEndpoint)
		} else {
			os.Unsetenv("VICOHOME_API_ENDPOINT")
		}
	}()
	os.Setenv("VICOHOME_API_ENDPOINT", mockServer.Server.URL)

	// Set up mock credentials
	os.Setenv("VICOHOME_EMAIL", "test@example.com")
	os.Setenv("VICOHOME_PASSWORD", "password123")

	// Create mock event data
	mockEvents := generateMockEvents(100) // Generate 100 events

	// Set up different responses for different time ranges
	now := time.Now()
	oneDayAgo := now.Add(-24 * time.Hour)
	twoDaysAgo := now.Add(-48 * time.Hour)
	threeDaysAgo := now.Add(-72 * time.Hour)
	
	// Last 24 hours response (30 events)
	mockServer.AddEventResponse(
		fmt.Sprintf("%d", oneDayAgo.Unix()),
		fmt.Sprintf("%d", now.Unix()),
		map[string]interface{}{
			"result": 0,
			"msg":    "success",
			"data": map[string]interface{}{
				"list": mockEvents[0:30],
			},
		},
	)

	// 24-48 hours response (40 events)
	mockServer.AddEventResponse(
		fmt.Sprintf("%d", twoDaysAgo.Unix()),
		fmt.Sprintf("%d", oneDayAgo.Unix()),
		map[string]interface{}{
			"result": 0,
			"msg":    "success",
			"data": map[string]interface{}{
				"list": mockEvents[30:70],
			},
		},
	)

	// 48-72 hours response (30 events)
	mockServer.AddEventResponse(
		fmt.Sprintf("%d", threeDaysAgo.Unix()),
		fmt.Sprintf("%d", twoDaysAgo.Unix()),
		map[string]interface{}{
			"result": 0,
			"msg":    "success",
			"data": map[string]interface{}{
				"list": mockEvents[70:100],
			},
		},
	)

	// Test case 1: Fetch events for the last 24 hours
	t.Run("Fetch events for last 24 hours", func(t *testing.T) {
		mockServer.RequestLog = nil // Clear request log
		
		// Fetch events (24 hours)
		events, err := fetchEventsWithPagination(1, mockServer)
		
		// Validate
		if err != nil {
			t.Fatalf("Error fetching events: %v", err)
		}
		if len(events) != 30 {
			t.Errorf("Expected 30 events, got %d", len(events))
		}
		validateEventOrder(t, events)
	})

	// Test case 2: Fetch events for the last 48 hours
	t.Run("Fetch events for last 48 hours", func(t *testing.T) {
		mockServer.RequestLog = nil // Clear request log
		
		// Fetch events (48 hours)
		events, err := fetchEventsWithPagination(2, mockServer)
		
		// Validate
		if err != nil {
			t.Fatalf("Error fetching events: %v", err)
		}
		if len(events) != 70 {
			t.Errorf("Expected 70 events (30 + 40), got %d", len(events))
		}
		validateEventOrder(t, events)
	})

	// Test case 3: Fetch events for the last 72 hours
	t.Run("Fetch events for last 72 hours", func(t *testing.T) {
		mockServer.RequestLog = nil // Clear request log
		
		// Fetch events (72 hours)
		events, err := fetchEventsWithPagination(3, mockServer)
		
		// Validate
		if err != nil {
			t.Fatalf("Error fetching events: %v", err)
		}
		if len(events) != 100 {
			t.Errorf("Expected 100 events (30 + 40 + 30), got %d", len(events))
		}
		validateEventOrder(t, events)
	})
}

// fetchEventsWithPagination fetches events for a given number of days with pagination.
func fetchEventsWithPagination(days int, mockServer *MockAPIServer) ([]MockEvent, error) {
	// Login to get token
	token := ""
	for k, v := range mockServer.AuthTokens {
		if v {
			token = k
			break
		}
	}
	if token == "" {
		// Since we can't directly call the handleLogin method with our request,
		// we'll simulate a login by adding a token to the AuthTokens map
		token = fmt.Sprintf("mock-token-%d", time.Now().UnixNano())
		mockServer.AuthTokens[token] = true
	}
	
	// Calculate time ranges for pagination
	now := time.Now()
	end := now
	var allEvents []MockEvent
	
	// Fetch events in 24-hour chunks
	for i := 0; i < days; i++ {
		start := end.Add(-24 * time.Hour)
		events, err := fetchEventsForTimeRange(token, start, end, mockServer)
		if err != nil {
			return nil, err
		}
		// Append new events to accumulated events (proper order)
		allEvents = append(allEvents, events...)
		end = start
	}
	
	return allEvents, nil
}

// fetchEventsForTimeRange fetches events for a specific time range.
func fetchEventsForTimeRange(token string, start, end time.Time, mockServer *MockAPIServer) ([]MockEvent, error) {
	// Create request object
	req := map[string]interface{}{
		"startTimestamp": fmt.Sprintf("%d", start.Unix()),
		"endTimestamp":   fmt.Sprintf("%d", end.Unix()),
		"language":       "en",
		"countryNo":      "US",
	}
	
	// Simulate fetching events - create HTTP request to mock server
	reqJSON, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", mockServer.Server.URL+"/library/newselectlibrary", bytes.NewBuffer(reqJSON))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)
	
	// Send the request to the mock server
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()
	
	// Check response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from server: %d", resp.StatusCode)
	}
	
	// Decode the response
	var responseMap map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&responseMap); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}
	
	// Check for errors in the response
	if result, ok := responseMap["result"].(float64); ok && result != 0 {
		msg, _ := responseMap["msg"].(string)
		return nil, fmt.Errorf("API error: %s (code: %.0f)", msg, result)
	}
	
	// Now extract the events from the response
	return transformToListEvents(responseMap)
}

// transformToListEvents transforms the API response to a list of test event objects.
func transformToListEvents(resp map[string]interface{}) ([]MockEvent, error) {
	var events []MockEvent
	
	if data, ok := resp["data"].(map[string]interface{}); ok {
		if list, ok := data["list"].([]interface{}); ok {
			for _, item := range list {
				if eventMap, ok := item.(map[string]interface{}); ok {
					// Transform the event to our format (simplified for test)
					event := MockEvent{
						TraceID:    fmt.Sprintf("%v", eventMap["traceId"]),
						Timestamp:  fmt.Sprintf("%v", eventMap["timestamp"]),
						DeviceName: fmt.Sprintf("%v", eventMap["deviceName"]),
					}
					events = append(events, event)
				}
			}
		}
	}
	
	return events, nil
}

// generateMockEvents generates a list of mock events for testing.
func generateMockEvents(count int) []map[string]interface{} {
	mockEvents := make([]map[string]interface{}, 0, count)
	baseTime := time.Now()
	
	for i := 0; i < count; i++ {
		// Create event with decreasing timestamps (newer events first)
		eventTime := baseTime.Add(-time.Duration(i) * time.Minute)
		
		event := map[string]interface{}{
			"traceId":        fmt.Sprintf("trace-%d", i),
			"timestamp":      fmt.Sprintf("%d", eventTime.Unix()),
			"deviceName":     fmt.Sprintf("Device-%d", i%5), // 5 different devices
			"serialNumber":   fmt.Sprintf("SN%06d", i),
			"adminName":      "Admin User",
			"period":         "10.0",
			"birdName":       fmt.Sprintf("Bird-%d", i%10), // 10 different birds
			"birdLatin":      fmt.Sprintf("Latin-%d", i%10),
			"birdConfidence": 0.85 + (float64(i%10) / 100.0),
			"imageUrl":       fmt.Sprintf("https://example.com/image-%d.jpg", i),
			"videoUrl":       fmt.Sprintf("https://example.com/video-%d.mp4", i),
			"keyshots": []map[string]interface{}{
				{
					"imageUrl":       fmt.Sprintf("https://example.com/keyshot-%d.jpg", i),
					"message":        "Bird detected",
					"objectCategory": "bird",
					"subCategoryName": fmt.Sprintf("Bird-%d", i%10),
				},
			},
			"subcategoryInfoList": []map[string]interface{}{
				{
					"objectType":    "bird",
					"objectName":    fmt.Sprintf("Bird-%d", i%10),
					"birdStdName":   fmt.Sprintf("Latin-%d", i%10),
					"confidence":    0.85 + (float64(i%10) / 100.0),
				},
			},
		}
		
		mockEvents = append(mockEvents, event)
	}
	
	return mockEvents
}

// validateEventOrder checks that events are in the expected order (newest first).
func validateEventOrder(t *testing.T, events []MockEvent) {
	var lastTimestamp int64 = 0
	
	for i, event := range events {
		// Parse timestamp
		timestamp := parseTimestamp(event.Timestamp)
		if i > 0 && timestamp > lastTimestamp {
			t.Errorf("Events are not in correct order. Event %d has newer timestamp than event %d", i, i-1)
		}
		lastTimestamp = timestamp
	}
}

// parseTimestamp parses a timestamp string to a Unix timestamp.
func parseTimestamp(timestamp string) int64 {
	// Check if it's already a Unix timestamp
	unixTime, err := strconv.ParseInt(timestamp, 10, 64)
	if err == nil {
		return unixTime
	}
	
	// Try parsing as a date string
	formats := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	
	for _, format := range formats {
		t, err := time.Parse(format, timestamp)
		if err == nil {
			return t.Unix()
		}
	}
	
	// Default to current time if unable to parse
	return time.Now().Unix()
}