package models

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// TestEventJSONMarshaling tests the marshaling and unmarshaling of Event objects to/from JSON.
func TestEventJSONMarshaling(t *testing.T) {
	// Create a sample event
	event := Event{
		TraceID:        "trace-123",
		Timestamp:      "2023-05-09T10:15:30Z",
		DeviceName:     "TestDevice",
		SerialNumber:   "SN123456",
		AdminName:      "Admin User",
		Period:         "10.0s",
		BirdName:       "American Robin",
		BirdLatin:      "Turdus migratorius",
		BirdConfidence: 0.95,
		KeyShotURL:     "https://example.com/keyshot.jpg",
		ImageURL:       "https://example.com/image.jpg",
		VideoURL:       "https://example.com/video.mp4",
		keyshots:       []map[string]interface{}{{"url": "https://example.com/keyshot.jpg"}},
	}

	// Test marshaling
	t.Run("Marshal Event to JSON", func(t *testing.T) {
		jsonData, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("Failed to marshal Event to JSON: %v", err)
		}

		// Verify internal fields are not included
		if string(jsonData) == "" {
			t.Error("JSON output is empty")
		}
		if string(jsonData) == "{}" {
			t.Error("JSON output is empty object")
		}

		// Verify internal field (keyshots) is not marshaled
		if string(jsonData) != "" && string(jsonData) != "{}" {
			var jsonMap map[string]interface{}
			err = json.Unmarshal(jsonData, &jsonMap)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON to map: %v", err)
			}

			// Check that keyshots field is not present
			if _, ok := jsonMap["keyshots"]; ok {
				t.Error("Internal field 'keyshots' was marshaled to JSON")
			}
		}
	})

	// Test unmarshaling
	t.Run("Unmarshal JSON to Event", func(t *testing.T) {
		jsonData := `{
			"traceId": "trace-123",
			"timestamp": "2023-05-09T10:15:30Z",
			"deviceName": "TestDevice",
			"serialNumber": "SN123456",
			"adminName": "Admin User",
			"period": "10.0s",
			"birdName": "American Robin",
			"birdLatin": "Turdus migratorius",
			"birdConfidence": 0.95,
			"keyShotUrl": "https://example.com/keyshot.jpg",
			"imageUrl": "https://example.com/image.jpg",
			"videoUrl": "https://example.com/video.mp4"
		}`

		var unmarshaledEvent Event
		err := json.Unmarshal([]byte(jsonData), &unmarshaledEvent)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON to Event: %v", err)
		}

		// Create expected event (without keyshots as it's not in JSON)
		expectedEvent := Event{
			TraceID:        "trace-123",
			Timestamp:      "2023-05-09T10:15:30Z",
			DeviceName:     "TestDevice",
			SerialNumber:   "SN123456",
			AdminName:      "Admin User",
			Period:         "10.0s",
			BirdName:       "American Robin",
			BirdLatin:      "Turdus migratorius",
			BirdConfidence: 0.95,
			KeyShotURL:     "https://example.com/keyshot.jpg",
			ImageURL:       "https://example.com/image.jpg",
			VideoURL:       "https://example.com/video.mp4",
		}

		// Check that deserialized event matches expected
		if !reflect.DeepEqual(unmarshaledEvent, expectedEvent) {
			t.Errorf("Unmarshaled event doesn't match expected event.\nGot: %+v\nExpected: %+v", 
				unmarshaledEvent, expectedEvent)
		}
	})

	// Test round-trip marshaling and unmarshaling
	t.Run("Marshal and then Unmarshal (round-trip)", func(t *testing.T) {
		// Create a copy of the event without the keyshots (internal field)
		eventWithoutKeyshots := Event{
			TraceID:        event.TraceID,
			Timestamp:      event.Timestamp,
			DeviceName:     event.DeviceName,
			SerialNumber:   event.SerialNumber,
			AdminName:      event.AdminName,
			Period:         event.Period,
			BirdName:       event.BirdName,
			BirdLatin:      event.BirdLatin,
			BirdConfidence: event.BirdConfidence,
			KeyShotURL:     event.KeyShotURL,
			ImageURL:       event.ImageURL,
			VideoURL:       event.VideoURL,
		}

		// Marshal to JSON
		jsonData, err := json.Marshal(event)
		if err != nil {
			t.Fatalf("Failed to marshal Event to JSON: %v", err)
		}

		// Unmarshal back to a new Event
		var roundTripEvent Event
		err = json.Unmarshal(jsonData, &roundTripEvent)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON back to Event: %v", err)
		}

		// Compare with expected event (without keyshots)
		if !reflect.DeepEqual(roundTripEvent, eventWithoutKeyshots) {
			t.Errorf("Round-trip marshaling/unmarshaling produced different event.\nGot: %+v\nExpected: %+v", 
				roundTripEvent, eventWithoutKeyshots)
		}
	})
}

// TestEventAPIMapping tests the mapping between API responses and Event objects.
func TestEventAPIMapping(t *testing.T) {
	// Test complete API response mapping
	t.Run("Map complete API response to Event", func(t *testing.T) {
		// Sample API response JSON (complete with all fields)
		apiJSON := `{
			"traceId": "api-trace-123",
			"timestamp": "2023-05-09T15:30:45Z",
			"deviceName": "API Device",
			"serialNumber": "APISN789",
			"adminName": "API Admin",
			"period": "15.5s",
			"birdName": "Blue Jay",
			"birdLatin": "Cyanocitta cristata",
			"birdConfidence": 0.87,
			"keyShotUrl": "https://example.com/api-keyshot.jpg",
			"imageUrl": "https://example.com/api-image.jpg",
			"videoUrl": "https://example.com/api-video.mp4",
			"extraField1": "This field should be ignored",
			"extraField2": 123
		}`

		// Expected Event after mapping
		expectedEvent := Event{
			TraceID:        "api-trace-123",
			Timestamp:      "2023-05-09T15:30:45Z",
			DeviceName:     "API Device",
			SerialNumber:   "APISN789",
			AdminName:      "API Admin",
			Period:         "15.5s",
			BirdName:       "Blue Jay",
			BirdLatin:      "Cyanocitta cristata",
			BirdConfidence: 0.87,
			KeyShotURL:     "https://example.com/api-keyshot.jpg",
			ImageURL:       "https://example.com/api-image.jpg",
			VideoURL:       "https://example.com/api-video.mp4",
		}

		// Unmarshal JSON to Event
		var actualEvent Event
		err := json.Unmarshal([]byte(apiJSON), &actualEvent)
		if err != nil {
			t.Fatalf("Failed to map API JSON to Event: %v", err)
		}

		// Verify mapping is correct
		if !reflect.DeepEqual(actualEvent, expectedEvent) {
			t.Errorf("API mapping produced incorrect Event.\nGot: %+v\nExpected: %+v", 
				actualEvent, expectedEvent)
		}
	})

	// Test partial API response mapping (missing some fields)
	t.Run("Map partial API response to Event", func(t *testing.T) {
		// Sample API response JSON (missing some fields)
		partialJSON := `{
			"traceId": "partial-trace-456",
			"timestamp": "2023-05-10T08:15:00Z",
			"deviceName": "Partial Device",
			"serialNumber": "PSN456",
			"birdName": "Cardinal",
			"birdLatin": "Cardinalis cardinalis",
			"birdConfidence": 0.92
		}`

		// Expected Event after mapping (missing fields should be empty strings or zero values)
		expectedEvent := Event{
			TraceID:        "partial-trace-456",
			Timestamp:      "2023-05-10T08:15:00Z",
			DeviceName:     "Partial Device",
			SerialNumber:   "PSN456",
			AdminName:      "", // Missing in JSON
			Period:         "", // Missing in JSON
			BirdName:       "Cardinal",
			BirdLatin:      "Cardinalis cardinalis",
			BirdConfidence: 0.92,
			KeyShotURL:     "", // Missing in JSON
			ImageURL:       "", // Missing in JSON
			VideoURL:       "", // Missing in JSON
		}

		// Unmarshal JSON to Event
		var actualEvent Event
		err := json.Unmarshal([]byte(partialJSON), &actualEvent)
		if err != nil {
			t.Fatalf("Failed to map partial API JSON to Event: %v", err)
		}

		// Verify mapping is correct
		if !reflect.DeepEqual(actualEvent, expectedEvent) {
			t.Errorf("Partial API mapping produced incorrect Event.\nGot: %+v\nExpected: %+v", 
				actualEvent, expectedEvent)
		}
	})

	// Test handling of invalid field types
	t.Run("Handle invalid field types in API response", func(t *testing.T) {
		// Sample API response with invalid field types
		invalidJSON := `{
			"traceId": "invalid-trace-789",
			"timestamp": "2023-05-11T12:30:00Z",
			"deviceName": "Invalid Device",
			"serialNumber": "ISN789",
			"birdName": "Sparrow",
			"birdLatin": "Passer domesticus",
			"birdConfidence": "0.75" 
		}`

		// Try to unmarshal - this should fail due to birdConfidence being a string instead of float64
		var event Event
		err := json.Unmarshal([]byte(invalidJSON), &event)
		
		// For Go's json package, type mismatches do cause unmarshaling errors
		if err == nil {
			t.Error("Expected error when unmarshaling JSON with invalid field types, but got no error")
		}

		// If we did get an error, verify it's related to the field type mismatch
		if err != nil && err.Error() == "" {
			t.Errorf("Got an error when unmarshaling JSON with invalid field types, but it wasn't descriptive: %v", err)
		}
	})
}

// TestEventValidation tests any validation logic for Event objects.
// For this model, we're not implementing explicit validation, but we'll test
// that events with missing or empty fields can be created and used.
func TestEventValidation(t *testing.T) {
	// Test event with all fields empty
	t.Run("Create event with all fields empty", func(t *testing.T) {
		emptyEvent := Event{}

		// Verify all fields have zero values
		if emptyEvent.TraceID != "" || 
		   emptyEvent.Timestamp != "" || 
		   emptyEvent.DeviceName != "" || 
		   emptyEvent.SerialNumber != "" ||
		   emptyEvent.AdminName != "" || 
		   emptyEvent.Period != "" || 
		   emptyEvent.BirdName != "" ||
		   emptyEvent.BirdLatin != "" || 
		   emptyEvent.BirdConfidence != 0 || 
		   emptyEvent.KeyShotURL != "" ||
		   emptyEvent.ImageURL != "" || 
		   emptyEvent.VideoURL != "" {
			t.Error("Empty event should have all fields set to zero values")
		}

		// Verify we can marshal an empty event to JSON without errors
		_, err := json.Marshal(emptyEvent)
		if err != nil {
			t.Errorf("Failed to marshal empty Event to JSON: %v", err)
		}
	})

	// Test event with minimal fields (only the ones that might be required in practice)
	t.Run("Create event with only essential fields", func(t *testing.T) {
		minimalEvent := Event{
			TraceID:    "min-trace-123",
			Timestamp:  "2023-05-12T09:45:00Z",
			DeviceName: "Minimal Device",
		}

		// Verify we can marshal a minimal event to JSON without errors
		jsonData, err := json.Marshal(minimalEvent)
		if err != nil {
			t.Errorf("Failed to marshal minimal Event to JSON: %v", err)
		}

		// Verify the JSON contains the fields we provided
		jsonString := string(jsonData)
		if !contains(jsonString, "min-trace-123") || 
		   !contains(jsonString, "2023-05-12T09:45:00Z") || 
		   !contains(jsonString, "Minimal Device") {
			t.Errorf("JSON for minimal event doesn't contain expected values: %s", jsonString)
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}