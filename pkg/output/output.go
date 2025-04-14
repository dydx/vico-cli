// Package output provides interfaces and implementations for different output formats and destinations.
package output

import (
	"fmt"

	"github.com/dydx/vico-cli/pkg/models"
	"github.com/dydx/vico-cli/pkg/output/config"
	"github.com/dydx/vico-cli/pkg/output/influxdb"
	"github.com/dydx/vico-cli/pkg/output/stdout"
)

// OutputHandler defines the interface for handling event output.
// Implementations of this interface can output events to different destinations
// (e.g., stdout, InfluxDB) in different formats (e.g., table, JSON).
type OutputHandler interface {
	// Write outputs the events using the configured format and destination.
	Write(events []models.Event) error

	// Close releases any resources used by the handler.
	Close()
}

// Factory creates an OutputHandler based on the specified destination and format.
func Factory(destination, format string, cfg Config) (OutputHandler, error) {
	switch destination {
	case "stdout":
		return NewStdoutHandler(format), nil
	case "influxdb":
		// Import the influxdb package and create a handler
		handler, err := influxdb.NewHandler(cfg.InfluxURL, cfg.InfluxOrg, cfg.InfluxBucket, cfg.InfluxToken)
		if err != nil {
			return nil, fmt.Errorf("failed to create InfluxDB handler: %w", err)
		}
		return handler, nil
	default:
		return nil, fmt.Errorf("unsupported output destination: %s", destination)
	}
}

// Config is an alias for config.Config to maintain backward compatibility
type Config = config.Config

// LoadConfigFromEnv loads configuration from environment variables.
func LoadConfigFromEnv() Config {
	return config.LoadFromEnv()
}

// NewStdoutHandler creates a new stdout output handler.
func NewStdoutHandler(format string) OutputHandler {
	switch format {
	case "json":
		return stdout.NewJSONHandler()
	default:
		return stdout.NewTableHandler()
	}
}

// NewInfluxDBHandler creates a new InfluxDB output handler.
func NewInfluxDBHandler(cfg Config) (OutputHandler, error) {
	// Validate required configuration
	if cfg.InfluxURL == "" {
		return nil, fmt.Errorf("InfluxDB URL is required")
	}
	if cfg.InfluxOrg == "" {
		return nil, fmt.Errorf("InfluxDB organization is required")
	}
	if cfg.InfluxBucket == "" {
		return nil, fmt.Errorf("InfluxDB bucket is required")
	}
	if cfg.InfluxToken == "" {
		return nil, fmt.Errorf("InfluxDB token is required")
	}

	// We'll use a function in the main package to create the handler
	// to avoid import cycles
	return nil, fmt.Errorf("InfluxDB handler creation is handled in the main package")
}
