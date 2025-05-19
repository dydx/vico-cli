package events

import (
	"fmt"
	"time"
)

// parseTimeParameters validates and parses the start and end time parameters
func parseTimeParameters(startTime, endTime string) (time.Time, time.Time, error) {
	// List of supported time formats
	formats := []string{
		"2006-01-02 15:04:05", // Standard format
		time.RFC3339,          // ISO 8601 format
	}

	var start, end time.Time
	var err error
	startParsed := false
	endParsed := false

	// Try to parse start time with different formats
	for _, format := range formats {
		start, err = time.Parse(format, startTime)
		if err == nil {
			startParsed = true
			break
		}
	}

	if !startParsed {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start time format: %v", err)
	}

	// Try to parse end time with different formats
	for _, format := range formats {
		end, err = time.Parse(format, endTime)
		if err == nil {
			endParsed = true
			break
		}
	}

	if !endParsed {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end time format: %v", err)
	}

	// Validate that start is before end
	if start.After(end) {
		return time.Time{}, time.Time{}, fmt.Errorf("start time must be before end time")
	}

	return start, end, nil
}
