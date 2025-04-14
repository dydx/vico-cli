// Package config provides configuration structures for output handlers.
package config

import (
	"os"
)

// Config holds configuration for output handlers.
type Config struct {
	InfluxURL    string
	InfluxOrg    string
	InfluxBucket string
	InfluxToken  string
}

// LoadFromEnv loads configuration from environment variables.
func LoadFromEnv() Config {
	return Config{
		InfluxURL:    getEnv("INFLUX_URL", ""),
		InfluxOrg:    getEnv("INFLUX_ORG", ""),
		InfluxBucket: getEnv("INFLUX_BUCKET", ""),
		InfluxToken:  getEnv("INFLUX_TOKEN", ""),
	}
}

// getEnv gets an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
