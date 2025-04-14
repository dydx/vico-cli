// Package models provides data models for the Vicohome CLI application.
package models

// Event represents a Vicohome event with its properties as returned by the API.
// This structure contains information about bird sightings, including metadata
// about the device that captured the event, the bird identified, and media URLs.
type Event struct {
	TraceID        string  `json:"traceId"`
	Timestamp      string  `json:"timestamp"`
	DeviceName     string  `json:"deviceName"`
	SerialNumber   string  `json:"serialNumber"`
	AdminName      string  `json:"adminName"`
	Period         string  `json:"period"`
	BirdName       string  `json:"birdName"`
	BirdLatin      string  `json:"birdLatin"`
	BirdConfidence float64 `json:"birdConfidence"`
	KeyShotURL     string  `json:"keyShotUrl"`
	ImageURL       string  `json:"imageUrl"`
	VideoURL       string  `json:"videoUrl"`

	// Internal field - not exported to JSON
	keyshots []map[string]interface{} `json:"-"`
}
