package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
)

// TestDeviceListingFlows tests the device listing functionality.
func TestDeviceListingFlows(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TESTS=true to run.")
	}

	// Create a mock API server
	mockServer := NewMockAPIServer()
	defer mockServer.Close()

	// Override the API endpoint to use our mock server
	originalAPIEndpoint := os.Getenv("VICOHOME_API_ENDPOINT")
	defer func() {
		if originalAPIEndpoint != "" {
			os.Setenv("VICOHOME_API_ENDPOINT", originalAPIEndpoint)
		} else {
			os.Unsetenv("VICOHOME_API_ENDPOINT")
		}
	}()
	os.Setenv("VICOHOME_API_ENDPOINT", mockServer.Server.URL)

	// Set up mock credentials
	os.Setenv("VICOHOME_EMAIL", "test@example.com")
	os.Setenv("VICOHOME_PASSWORD", "password123")

	// Create mock device data
	mockDevices := generateMockDevices(10) // Generate 10 devices

	// Set up response for English/US locale
	mockServer.AddDeviceResponse("en", "US", map[string]interface{}{
		"result": 0,
		"msg":    "success",
		"data": map[string]interface{}{
			"list": mockDevices,
		},
	})

	// Set up empty response for other locales
	mockServer.AddDeviceResponse("fr", "FR", map[string]interface{}{
		"result": 0,
		"msg":    "success",
		"data": map[string]interface{}{
			"list": []interface{}{},
		},
	})

	// Test case 1: Fetch devices with English/US locale
	t.Run("Fetch devices with English/US locale", func(t *testing.T) {
		mockServer.RequestLog = nil // Clear request log
		
		// Fetch devices
		devices, err := fetchDevices("en", "US", mockServer)
		
		// Validate
		if err != nil {
			t.Fatalf("Error fetching devices: %v", err)
		}
		if len(devices) != 10 {
			t.Errorf("Expected 10 devices, got %d", len(devices))
		}
		validateDevices(t, devices)
	})

	// Test case 2: Fetch devices with French/France locale (empty result)
	t.Run("Fetch devices with French/France locale", func(t *testing.T) {
		mockServer.RequestLog = nil // Clear request log
		
		// Fetch devices
		devices, err := fetchDevices("fr", "FR", mockServer)
		
		// Validate
		if err != nil {
			t.Fatalf("Error fetching devices: %v", err)
		}
		if len(devices) != 0 {
			t.Errorf("Expected 0 devices, got %d", len(devices))
		}
	})

	// Test case 3: Test authentication failure
	t.Run("Authentication failure", func(t *testing.T) {
		mockServer.RequestLog = nil // Clear request log
		mockServer.AuthFailures = 1 // Simulate auth failure
		mockServer.AuthTokens = make(map[string]bool) // Clear tokens to force login
		
		// Fetch devices
		_, err := fetchDevices("en", "US", mockServer)
		
		// Validate
		if err == nil {
			t.Fatalf("Expected authentication error, but got none")
		}
		if !strings.Contains(err.Error(), "API error") && !strings.Contains(err.Error(), "failed") {
			t.Errorf("Expected authentication error message, got: %v", err)
		}
	})
}

// MockDevice represents a simplified device structure for testing.
type MockDevice struct {
	SerialNumber string
	ModelNo      string
	DeviceName   string
	NetworkName  string
	IP           string
	BatteryLevel int
	MacAddress   string
}

// fetchDevices fetches devices with the specified locale.
func fetchDevices(language, countryNo string, mockServer *MockAPIServer) ([]MockDevice, error) {
	// Login to get token
	token := ""
	for k, v := range mockServer.AuthTokens {
		if v {
			token = k
			break
		}
	}
	
	// Generate a new token if needed
	if token == "" {
		// Create login request
		loginReq := map[string]interface{}{
			"email":     "test@example.com",
			"password":  "password123",
			"loginType": 0,
		}
		reqJSON, _ := json.Marshal(loginReq)
		httpReq, _ := http.NewRequest("POST", mockServer.Server.URL+"/account/login", bytes.NewBuffer(reqJSON))
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Accept", "application/json")
		
		// Send login request
		client := &http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			return nil, fmt.Errorf("error making login request: %v", err)
		}
		defer resp.Body.Close()
		
		// Decode login response
		var loginResponse map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&loginResponse); err != nil {
			return nil, fmt.Errorf("error decoding login response: %v", err)
		}
		
		// Check for login errors
		if result, ok := loginResponse["result"].(float64); ok && result != 0 {
			msg, _ := loginResponse["msg"].(string)
			return nil, fmt.Errorf("API error during login: %s (code: %.0f)", msg, result)
		}
		
		// Extract token from login response
		if data, ok := loginResponse["data"].(map[string]interface{}); ok {
			if tokenObj, ok := data["token"].(map[string]interface{}); ok {
				if tokenStr, ok := tokenObj["token"].(string); ok {
					token = tokenStr
				}
			}
		}
		
		if token == "" {
			return nil, fmt.Errorf("failed to get token from login response")
		}
	}
	
	// Create device list request
	req := map[string]interface{}{
		"language":  language,
		"countryNo": countryNo,
	}
	
	// Simulate making a request to the mock server
	reqJSON, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", mockServer.Server.URL+"/device/listuserdevices", bytes.NewBuffer(reqJSON))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)
	
	// Send the request to the mock server
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()
	
	// Check response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from server: %d", resp.StatusCode)
	}
	
	// Decode the response
	var responseMap map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&responseMap); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}
	
	// Check for errors in the response
	if result, ok := responseMap["result"].(float64); ok && result != 0 {
		msg, _ := responseMap["msg"].(string)
		return nil, fmt.Errorf("API error: %s (code: %.0f)", msg, result)
	}
	
	// Now extract the devices from the response
	return transformToMockDevices(responseMap)
}

// transformToMockDevices transforms the API response to a list of test device objects.
func transformToMockDevices(resp map[string]interface{}) ([]MockDevice, error) {
	var devices []MockDevice
	
	if data, ok := resp["data"].(map[string]interface{}); ok {
		if list, ok := data["list"].([]interface{}); ok {
			for _, item := range list {
				if deviceMap, ok := item.(map[string]interface{}); ok {
					// Transform the device to our format (simplified for test)
					device := MockDevice{
						SerialNumber: toString(deviceMap["serialNumber"]),
						ModelNo:      toString(deviceMap["modelNo"]),
						DeviceName:   toString(deviceMap["deviceName"]),
						NetworkName:  toString(deviceMap["networkName"]),
						IP:           toString(deviceMap["ip"]),
						BatteryLevel: toInt(deviceMap["batteryLevel"]),
						MacAddress:   toString(deviceMap["macAddress"]),
					}
					devices = append(devices, device)
				}
			}
		}
	}
	
	return devices, nil
}

// toString safely converts an interface{} to string.
func toString(value interface{}) string {
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}

// toInt safely converts an interface{} to int.
func toInt(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case float32:
		return int(v)
	default:
		return 0
	}
}

// generateMockDevices generates a list of mock devices for testing.
func generateMockDevices(count int) []map[string]interface{} {
	mockDevices := make([]map[string]interface{}, 0, count)
	
	for i := 0; i < count; i++ {
		device := map[string]interface{}{
			"serialNumber":   fmt.Sprintf("SN%06d", i),
			"modelNo":        fmt.Sprintf("MODEL-%d", i%3), // 3 different models
			"deviceName":     fmt.Sprintf("Device-%d", i),
			"networkName":    fmt.Sprintf("Network-%d", i%2), // 2 different networks
			"ip":             fmt.Sprintf("192.168.1.%d", 10+i),
			"batteryLevel":   50 + (i * 5) % 50, // 50-99% battery
			"locationName":   fmt.Sprintf("Location-%d", i%3), // 3 different locations
			"signalStrength": 70 + (i * 3) % 30, // 70-99% signal
			"wifiChannel":    1 + (i % 11), // Channels 1-11
			"isCharging":     i % 2, // 0 or 1
			"chargingMode":   i % 3, // 0, 1, or 2
			"macAddress":     fmt.Sprintf("00:11:22:33:44:%02X", i),
		}
		
		mockDevices = append(mockDevices, device)
	}
	
	return mockDevices
}

// validateDevices checks that devices have valid properties.
func validateDevices(t *testing.T, devices []MockDevice) {
	for i, device := range devices {
		if device.SerialNumber == "" {
			t.Errorf("Device %d has empty serial number", i)
		}
		if device.DeviceName == "" {
			t.Errorf("Device %d has empty name", i)
		}
		if !strings.HasPrefix(device.SerialNumber, "SN") {
			t.Errorf("Device %d has invalid serial number format: %s", i, device.SerialNumber)
		}
	}
}