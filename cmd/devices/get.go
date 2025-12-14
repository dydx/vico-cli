package devices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dydx/vico-cli/pkg/auth"
	"github.com/spf13/cobra"
)

// DeviceRequest represents the JSON request body sent to the Vicohome API
// when fetching a specific device by serial number.
type DeviceRequest struct {
	SerialNumber string `json:"serialNumber"` // Unique identifier for the device
	Language     string `json:"language"`     // Language code (e.g., "en" for English)
	CountryNo    string `json:"countryNo"`    // Country code (e.g., "US" for United States)
}

// getCmd represents the command to retrieve details for a specific device by its serial number.
// It supports output in both table and JSON formats.
var getCmd = &cobra.Command{
	Use:   "get [serialNumber]",
	Short: "Get details for a specific device",
	Long:  `Fetch and display detailed information for a specific Vicohome device by its serial number.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serialNumber := args[0]

		token, err := auth.Authenticate()
		if err != nil {
			fmt.Printf("Authentication failed: %v\n", err)
			return
		}

		device, err := getDevice(token, serialNumber)
		if err != nil {
			fmt.Printf("Error fetching device: %v\n", err)
			return
		}

		// Display device details
		if outputFormat == "json" {
			// Output JSON format
			prettyJSON, err := json.MarshalIndent(device, "", "  ")
			if err != nil {
				fmt.Printf("Error formatting JSON: %v\n", err)
				return
			}
			fmt.Println(string(prettyJSON))
		} else {
			// Output formatted table
			fmt.Println("Device Details:")
			fmt.Println("------------------------------")
			fmt.Printf("Serial Number:   %s\n", device.SerialNumber)
			fmt.Printf("Model Number:    %s\n", device.ModelNo)
			fmt.Printf("Device Name:     %s\n", device.DeviceName)
			fmt.Printf("Network Name:    %s\n", device.NetworkName)
			fmt.Printf("IP Address:      %s\n", device.IP)
			fmt.Printf("Battery Level:   %d%%\n", device.BatteryLevel)
			fmt.Printf("Location:        %s\n", device.LocationName)
			fmt.Printf("Signal Strength: %d dBm\n", device.SignalStrength)
			fmt.Printf("WiFi Channel:    %d\n", device.WifiChannel)
			fmt.Printf("Is Charging:     %s\n", boolFromInt(device.IsCharging))
			fmt.Printf("Charging Mode:   %d\n", device.ChargingMode)
			fmt.Printf("MAC Address:     %s\n", device.MacAddress)
		}
	},
}

func init() {
	getCmd.Flags().StringVar(&outputFormat, "format", "table", "Output format (table or json)")
}

// getDevice fetches detailed information for a specific device from the Vicohome API.
// It takes an authentication token and the device's serial number, and returns
// a Device object and any error encountered.
// This function handles the API request, response parsing, and error handling including
// authentication refreshes when needed.
func getDevice(token string, serialNumber string) (Device, error) {
	req := DeviceRequest{
		SerialNumber: serialNumber,
		Language:     "en",
		CountryNo:    "US",
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return Device{}, fmt.Errorf("error marshaling request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", GetBaseURL()+"/device/selectsingledevice", bytes.NewBuffer(reqBody))
	if err != nil {
		return Device{}, fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)

	// Use ExecuteWithRetry for automatic token refresh
	respBody, err := auth.ExecuteWithRetry(httpReq)
	if err != nil {
		return Device{}, fmt.Errorf("error making request: %w", err)
	}

	// Parse response
	var responseMap map[string]interface{}
	if err := json.Unmarshal(respBody, &responseMap); err != nil {
		return Device{}, fmt.Errorf("error unmarshaling response: %w\nResponse: %s", err, string(respBody))
	}

	// Check if we need to refresh the token
	needsRefresh, apiError := auth.ValidateResponse(respBody)
	if apiError != nil {
		// There was an API error, but it's not a auth error requiring a retry
		if !needsRefresh {
			return Device{}, apiError
		}
		// Auth error was handled by ValidateResponse, but we should return with the error
		return Device{}, fmt.Errorf("authentication error: %v", apiError)
	}

	// Extract device data
	data, ok := responseMap["data"].(map[string]interface{})
	if !ok {
		return Device{}, fmt.Errorf("no device data found")
	}

	return transformToDevice(data), nil
}

// boolFromInt converts an integer value to a human-readable string representation
// of a boolean value. Any value greater than 0 returns "Yes", otherwise "No".
// This is used for display purposes when showing boolean properties from the API.
func boolFromInt(val int) string {
	if val > 0 {
		return "Yes"
	}
	return "No"
}
