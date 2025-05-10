// Package output provides interfaces and implementations for different output formats and destinations.
package output

import (
	"github.com/dydx/vico-cli/pkg/models"
	"github.com/dydx/vico-cli/pkg/output/stdout"
)

// Handler defines the interface for handling event output.
// Implementations of this interface can output events in different formats (e.g., table, JSON).
type Handler interface {
	// Write outputs the events using the configured format and destination.
	Write(events []models.Event) error

	// Close releases any resources used by the handler.
	Close()
}

// Factory creates a Handler based on the specified format.
func Factory(format string) (Handler, error) {
	return NewStdoutHandler(format), nil
}

// NewStdoutHandler creates a new stdout output handler.
func NewStdoutHandler(format string) Handler {
	switch format {
	case "json":
		return stdout.NewJSONHandler()
	default:
		return stdout.NewTableHandler()
	}
}
