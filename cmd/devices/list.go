package devices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"github.com/dydx/vico-cli/pkg/auth"
	"github.com/spf13/cobra"
)

// DeviceListRequest represents the JSON request body sent to the Vicohome API
// when listing user devices. It specifies language and country preferences.
type DeviceListRequest struct {
	Language  string `json:"language"`  // Language code (e.g., "en" for English)
	CountryNo string `json:"countryNo"` // Country code (e.g., "US" for United States)
}

// Device represents a Vicohome device with its properties as returned by the API.
// This structure contains essential information about a device that can be displayed
// to the user or used for further API calls.
type Device struct {
	SerialNumber   string `json:"serialNumber"`
	ModelNo        string `json:"modelNo"`
	DeviceName     string `json:"deviceName"`
	NetworkName    string `json:"networkName"`
	IP             string `json:"ip"`
	BatteryLevel   int    `json:"batteryLevel"`
	LocationName   string `json:"locationName"`
	SignalStrength int    `json:"signalStrength"`
	WifiChannel    int    `json:"wifiChannel"`
	IsCharging     int    `json:"isCharging"`
	ChargingMode   int    `json:"chargingMode"`
	MacAddress     string `json:"macAddress"`
}

var outputFormat string

// listCmd represents the command to list all devices associated with the user's account.
// It supports output in both table and JSON formats.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all user devices",
	Long:  `Fetch and display all devices associated with your Vicohome account.`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := auth.Authenticate()
		if err != nil {
			fmt.Printf("Authentication failed: %v\n", err)
			return
		}

		devices, err := listDevices(token)
		if err != nil {
			fmt.Printf("Error fetching devices: %v\n", err)
			return
		}

		// Display devices
		if len(devices) == 0 {
			fmt.Println("No devices found.")
			return
		}

		if outputFormat == "json" {
			// Output JSON format
			prettyJSON, err := json.MarshalIndent(devices, "", "  ")
			if err != nil {
				fmt.Printf("Error formatting JSON: %v\n", err)
				return
			}
			fmt.Println(string(prettyJSON))
		} else {
			// Output table format
			fmt.Printf("%-36s %-20s %-20s %-15s %-15s %-5s\n",
				"Serial Number", "Model", "Name", "Network", "IP", "Battery")
			fmt.Println("----------------------------------------------------------------------------------------------------------------")
			for _, device := range devices {
				fmt.Printf("%-36s %-20s %-20s %-15s %-15s %d%%\n",
					device.SerialNumber,
					device.ModelNo,
					device.DeviceName,
					device.NetworkName,
					device.IP,
					device.BatteryLevel)
			}
		}
	},
}

func init() {
	listCmd.Flags().StringVar(&outputFormat, "format", "table", "Output format (table or json)")
}
func GetCountry() string {
if v := os.Getenv("VICOHOME_COUNTRY"); v != "" {
return v
}
return "US"
}
// listDevices fetches all devices associated with the user's account from the Vicohome API.
// It takes an authentication token and returns a slice of Device objects and any error encountered.
// This function handles the API request, response parsing, and error handling including
// authentication refreshes when needed.
func listDevices(token string) ([]Device, error) {
	req := DeviceListRequest{
		Language:  "en",
		CountryNo: GetCountry(),
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", GetBaseURL()+"/device/listuserdevices", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)

	// Use ExecuteWithRetry for automatic token refresh
	respBody, err := auth.ExecuteWithRetry(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	// Parse response
	var responseMap map[string]interface{}
	if err := json.Unmarshal(respBody, &responseMap); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w\nResponse: %s", err, string(respBody))
	}

	// Check if we need to refresh the token
	needsRefresh, apiError := auth.ValidateResponse(respBody)
	if apiError != nil {
		// There was an API error, but it's not a auth error requiring a retry
		if !needsRefresh {
			return nil, apiError
		}
		// Auth error was handled by ValidateResponse, but we should return with the error
		return nil, fmt.Errorf("authentication error: %v", apiError)
	}

	// Extract device list
	data, ok := responseMap["data"].(map[string]interface{})
	if !ok {
		return []Device{}, nil
	}

	deviceList, ok := data["list"].([]interface{})
	if !ok {
		return []Device{}, nil
	}

	// Transform devices to our simpler format
	devices := make([]Device, 0, len(deviceList))
	for _, item := range deviceList {
		if deviceMap, ok := item.(map[string]interface{}); ok {
			device := transformToDevice(deviceMap)
			devices = append(devices, device)
		}
	}

	return devices, nil
}

// transformToDevice converts a map of device data from the API response into a Device struct.
// It safely extracts and type-converts the various device properties from the dynamic map
// into the strongly-typed Device structure. Missing fields in the map will result in
// zero values in the returned Device structure.
func transformToDevice(deviceMap map[string]interface{}) Device {
	device := Device{}

	// Extract string fields
	if val, ok := deviceMap["serialNumber"].(string); ok {
		device.SerialNumber = val
	}
	if val, ok := deviceMap["modelNo"].(string); ok {
		device.ModelNo = val
	}
	if val, ok := deviceMap["deviceName"].(string); ok {
		device.DeviceName = val
	}
	if val, ok := deviceMap["networkName"].(string); ok {
		device.NetworkName = val
	}
	if val, ok := deviceMap["ip"].(string); ok {
		device.IP = val
	}
	if val, ok := deviceMap["locationName"].(string); ok {
		device.LocationName = val
	}
	if val, ok := deviceMap["macAddress"].(string); ok {
		device.MacAddress = val
	}

	// Extract numeric fields
	if val, ok := deviceMap["batteryLevel"].(float64); ok {
		device.BatteryLevel = int(val)
	}
	if val, ok := deviceMap["signalStrength"].(float64); ok {
		device.SignalStrength = int(val)
	}
	if val, ok := deviceMap["wifiChannel"].(float64); ok {
		device.WifiChannel = int(val)
	}
	if val, ok := deviceMap["isCharging"].(float64); ok {
		device.IsCharging = int(val)
	}
	if val, ok := deviceMap["chargingMode"].(float64); ok {
		device.ChargingMode = int(val)
	}

	return device
}
