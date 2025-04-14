// Package influxdb provides functionality for writing events to InfluxDB.
package influxdb

import (
	"fmt"
	"time"

	"github.com/dydx/vico-cli/pkg/models"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// Handler implements the output.Handler interface for InfluxDB.
type Handler struct {
	client     influxdb2.Client
	writeAPI   api.WriteAPI
	org        string
	bucket     string
	maxRetries int
	retryDelay time.Duration
}

// Config holds the configuration for the InfluxDB handler.
type Config struct {
	URL        string
	Org        string
	Bucket     string
	Token      string
	MaxRetries int
	RetryDelay time.Duration
}

// NewHandler creates a new InfluxDB handler.
func NewHandler(url, org, bucket, token string) (*Handler, error) {
	// Create client
	client := influxdb2.NewClient(url, token)
	writeAPI := client.WriteAPI(org, bucket)

	return &Handler{
		client:     client,
		writeAPI:   writeAPI,
		org:        org,
		bucket:     bucket,
		maxRetries: 3,                      // Default to 3 retries
		retryDelay: 500 * time.Millisecond, // Default to 500ms delay
	}, nil
}

// Write writes events to InfluxDB.
func (h *Handler) Write(events []models.Event) error {
	if len(events) == 0 {
		fmt.Println("No events to write to InfluxDB")
		return nil
	}

	fmt.Printf("Writing %d events to InfluxDB...\n", len(events))

	// Create error channel to capture async errors
	errorsCh := h.writeAPI.Errors()

	// Create a channel to signal completion of error handling
	done := make(chan bool)

	// Track if we've encountered any errors
	var writeError error

	// Start a goroutine to handle errors
	go func() {
		for err := range errorsCh {
			if writeError == nil {
				writeError = err
			} else {
				writeError = fmt.Errorf("%w; additional error: %v", writeError, err)
			}
		}
		done <- true
	}()

	// Write all points
	for _, event := range events {
		// Parse timestamp
		timestamp, err := parseTimestamp(event.Timestamp)
		if err != nil {
			fmt.Printf("Warning: Invalid timestamp '%s' for event with TraceID '%s': %v\n",
				event.Timestamp, event.TraceID, err)
			continue
		}

		// Create a point with measurement "bird_sighting"
		point := influxdb2.NewPoint(
			"bird_sighting",
			map[string]string{
				"device":     event.DeviceName,
				"serial":     event.SerialNumber,
				"bird_name":  event.BirdName,
				"bird_latin": event.BirdLatin,
				"trace_id":   event.TraceID,
			},
			map[string]interface{}{
				"confidence": event.BirdConfidence,
			},
			timestamp,
		)

		// Write the point
		h.writeAPI.WritePoint(point)
	}

	// Force writing of buffered points
	h.writeAPI.Flush()

	// Wait for any errors to be processed
	<-done

	// If no errors occurred, we're done
	if writeError == nil {
		fmt.Printf("Successfully wrote %d events to InfluxDB\n", len(events))
		return nil
	}

	return fmt.Errorf("error writing to InfluxDB: %w", writeError)
}

// Close closes the InfluxDB client.
func (h *Handler) Close() {
	h.client.Close()
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
