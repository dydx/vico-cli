// Package stdout provides implementations for outputting events to standard output.
package stdout

import (
	"encoding/json"
	"fmt"

	"github.com/dydx/vico-cli/pkg/models"
)

// JSONHandler outputs events in JSON format to stdout.
type JSONHandler struct{}

// NewJSONHandler creates a new JSON stdout handler.
func NewJSONHandler() *JSONHandler {
	return &JSONHandler{}
}

// Write outputs the events in JSON format to stdout.
func (h *JSONHandler) Write(events []models.Event) error {
	prettyJSON, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return fmt.Errorf("error formatting JSON: %w", err)
	}
	fmt.Println(string(prettyJSON))
	return nil
}

// Close is a no-op for stdout handlers.
func (h *JSONHandler) Close() {
	// No resources to release
}

// TableHandler outputs events in table format to stdout.
type TableHandler struct{}

// NewTableHandler creates a new table stdout handler.
func NewTableHandler() *TableHandler {
	return &TableHandler{}
}

// Write outputs the events in table format to stdout.
func (h *TableHandler) Write(events []models.Event) error {
	if len(events) == 0 {
		fmt.Println("No events found in the specified time period.")
		return nil
	}

	// Print table header
	fmt.Printf("%-36s %-20s %-25s %-25s %-25s\n",
		"Trace ID", "Timestamp", "Device Name", "Bird Name", "Bird Latin")
	fmt.Println("--------------------------------------------------------------------------------------------------")

	// Print table rows
	for _, event := range events {
		fmt.Printf("%-36s %-20s %-25s %-25s %-25s\n",
			event.TraceID,
			event.Timestamp,
			event.DeviceName,
			event.BirdName,
			event.BirdLatin)
	}

	return nil
}

// Close is a no-op for stdout handlers.
func (h *TableHandler) Close() {
	// No resources to release
}
