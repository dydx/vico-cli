package devices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dydx/vico-cli/pkg/auth"
	"github.com/spf13/cobra"
)

type DeviceRequest struct {
	SerialNumber string `json:"serialNumber"`
	Language     string `json:"language"`
	CountryNo    string `json:"countryNo"`
}

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

	httpReq, err := http.NewRequest("POST", "https://api-us.vicohome.io/device/selectsingledevice", bytes.NewBuffer(reqBody))
	if err != nil {
		return Device{}, fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return Device{}, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Device{}, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse response
	var responseMap map[string]interface{}
	if err := json.Unmarshal(respBody, &responseMap); err != nil {
		return Device{}, fmt.Errorf("error unmarshaling response: %w\nResponse: %s", err, string(respBody))
	}

	// Check for API errors
	if result, ok := responseMap["result"].(float64); ok && result != 0 {
		msg, _ := responseMap["msg"].(string)
		return Device{}, fmt.Errorf("API error: %s (code: %.0f)", msg, result)
	}

	// Extract device data
	data, ok := responseMap["data"].(map[string]interface{})
	if !ok {
		return Device{}, fmt.Errorf("no device data found")
	}

	return transformToDevice(data), nil
}

func boolFromInt(val int) string {
	if val > 0 {
		return "Yes"
	}
	return "No"
}