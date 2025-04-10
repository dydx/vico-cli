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

type DeviceListRequest struct {
	Language  string `json:"language"`
	CountryNo string `json:"countryNo"`
}

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

func listDevices(token string) ([]Device, error) {
	req := DeviceListRequest{
		Language:  "en",
		CountryNo: "US",
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", "https://api-us.vicohome.io/device/listuserdevices", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse response
	var responseMap map[string]interface{}
	if err := json.Unmarshal(respBody, &responseMap); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w\nResponse: %s", err, string(respBody))
	}

	// Check for API errors
	if result, ok := responseMap["result"].(float64); ok && result != 0 {
		msg, _ := responseMap["msg"].(string)
		return nil, fmt.Errorf("API error: %s (code: %.0f)", msg, result)
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