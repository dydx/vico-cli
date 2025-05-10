package stdout

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/dydx/vico-cli/pkg/models"
)

// captureOutput redirects stdout and returns the captured output
func captureOutput(f func() error) (string, error) {
	// Redirect stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the function that writes to stdout
	err := f()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Capture output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String(), err
}

// TestHandlerWrite tests all the output handlers with various event scenarios
func TestHandlerWrite(t *testing.T) {
	// Define common test events
	singleEvent := []models.Event{
		{
			TraceID:        "trace123",
			Timestamp:      "2023-01-01 12:00:00",
			DeviceName:     "TestDevice1",
			SerialNumber:   "SN12345",
			AdminName:      "Admin1",
			Period:         "10.00s",
			BirdName:       "Robin",
			BirdLatin:      "Turdus migratorius",
			BirdConfidence: 0.95,
			KeyShotURL:     "https://example.com/keyshot1.jpg",
			ImageURL:       "https://example.com/image1.jpg",
			VideoURL:       "https://example.com/video1.mp4",
		},
	}

	multipleEvents := append(singleEvent, models.Event{
		TraceID:        "trace456",
		Timestamp:      "2023-01-01 13:00:00",
		DeviceName:     "TestDevice2",
		SerialNumber:   "SN67890",
		AdminName:      "Admin2",
		Period:         "15.00s",
		BirdName:       "Blue Jay",
		BirdLatin:      "Cyanocitta cristata",
		BirdConfidence: 0.87,
		KeyShotURL:     "https://example.com/keyshot2.jpg",
		ImageURL:       "https://example.com/image2.jpg",
		VideoURL:       "https://example.com/video2.mp4",
	})

	emptyEvents := []models.Event{}

	// Define test cases
	tests := []struct {
		name                 string
		handler              interface{} // Either *JSONHandler or *TableHandler
		events               []models.Event
		expectError          bool
		validateOutput       func(t *testing.T, output string, events []models.Event)
	}{
		{
			name:        "JSONHandler with single event",
			handler:     NewJSONHandler(),
			events:      singleEvent,
			expectError: false,
			validateOutput: func(t *testing.T, output string, events []models.Event) {
				// Validate JSON formatting
				var capturedEvents []models.Event
				err := json.Unmarshal([]byte(output), &capturedEvents)
				if err != nil {
					t.Fatalf("Output is not valid JSON: %v", err)
				}

				// Verify content
				if len(capturedEvents) != len(events) {
					t.Errorf("Expected %d events, got %d", len(events), len(capturedEvents))
				}

				// Check specific fields
				if capturedEvents[0].TraceID != events[0].TraceID {
					t.Errorf("Event TraceID expected '%s', got '%s'", events[0].TraceID, capturedEvents[0].TraceID)
				}

				if capturedEvents[0].BirdName != events[0].BirdName {
					t.Errorf("Event BirdName expected '%s', got '%s'", events[0].BirdName, capturedEvents[0].BirdName)
				}
			},
		},
		{
			name:        "JSONHandler with multiple events",
			handler:     NewJSONHandler(),
			events:      multipleEvents,
			expectError: false,
			validateOutput: func(t *testing.T, output string, events []models.Event) {
				// Validate JSON formatting
				var capturedEvents []models.Event
				err := json.Unmarshal([]byte(output), &capturedEvents)
				if err != nil {
					t.Fatalf("Output is not valid JSON: %v", err)
				}

				// Verify content
				if len(capturedEvents) != len(events) {
					t.Errorf("Expected %d events, got %d", len(events), len(capturedEvents))
				}

				// Check fields from multiple events
				if capturedEvents[0].TraceID != events[0].TraceID {
					t.Errorf("First event TraceID expected '%s', got '%s'", events[0].TraceID, capturedEvents[0].TraceID)
				}
				
				if capturedEvents[1].BirdName != events[1].BirdName {
					t.Errorf("Second event BirdName expected '%s', got '%s'", events[1].BirdName, capturedEvents[1].BirdName)
				}
			},
		},
		{
			name:        "TableHandler with events",
			handler:     NewTableHandler(),
			events:      multipleEvents,
			expectError: false,
			validateOutput: func(t *testing.T, output string, events []models.Event) {
				// Verify output contains header
				if !strings.Contains(output, "Trace ID") || !strings.Contains(output, "Timestamp") {
					t.Errorf("Table header not found in output")
				}

				// Verify output contains event data
				for _, event := range events {
					if !strings.Contains(output, event.TraceID) || !strings.Contains(output, event.BirdName) {
						t.Errorf("Event data not found in output: TraceID=%s, BirdName=%s", 
							event.TraceID, event.BirdName)
					}
				}
			},
		},
		{
			name:        "TableHandler with empty events",
			handler:     NewTableHandler(),
			events:      emptyEvents,
			expectError: false,
			validateOutput: func(t *testing.T, output string, events []models.Event) {
				// Verify output for empty table
				expectedEmptyMessage := "No events found in the specified time period."
				if !strings.Contains(output, expectedEmptyMessage) {
					t.Errorf("Expected empty message '%s' not found in output", expectedEmptyMessage)
				}
			},
		},
	}

	// Execute test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			var output string
			
			// Type switch to call the appropriate handler
			switch h := tc.handler.(type) {
			case *JSONHandler:
				output, err = captureOutput(func() error {
					return h.Write(tc.events)
				})
			case *TableHandler:
				output, err = captureOutput(func() error {
					return h.Write(tc.events)
				})
			default:
				t.Fatalf("Unknown handler type: %T", tc.handler)
			}

			// Check error expectation
			if tc.expectError && err == nil {
				t.Error("Expected an error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Validate output
			tc.validateOutput(t, output, tc.events)
		})
	}
}

// TestHandlerClose ensures that the Close method doesn't produce errors
func TestHandlerClose(t *testing.T) {
	tests := []struct {
		name    string
		handler interface{}
	}{
		{"JSONHandler Close", NewJSONHandler()},
		{"TableHandler Close", NewTableHandler()},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Type switch to call the appropriate handler's Close method
			switch h := tc.handler.(type) {
			case *JSONHandler:
				h.Close()
			case *TableHandler:
				h.Close()
			default:
				t.Fatalf("Unknown handler type: %T", tc.handler)
			}
		})
	}
}