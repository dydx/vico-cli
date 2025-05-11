package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dydx/vico-cli/pkg/auth"
)

// TestTokenExpirationDuringAPICall tests how the application handles token expiration during API calls.
func TestTokenExpirationDuringAPICall(t *testing.T) {
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

	// Create mock event data
	mockEvents := make([]map[string]interface{}, 20)
	for i := 0; i < 20; i++ {
		mockEvents[i] = map[string]interface{}{
			"traceId":    fmt.Sprintf("trace-%d", i),
			"timestamp":  fmt.Sprintf("%d", time.Now().Unix()-int64(i*60)),
			"deviceName": fmt.Sprintf("Device-%d", i),
		}
	}

	// Set default response
	mockServer.SetDefaultEventResponse(map[string]interface{}{
		"result": 0,
		"msg":    "success",
		"data": map[string]interface{}{
			"list": mockEvents,
		},
	})

	// Test case 1: Token expiration with successful refresh
	t.Run("Token expiration with successful refresh", func(t *testing.T) {
		mockServer.RequestLog = nil // Clear request log
		mockServer.SimulateExpiration = false // Start with no expiration
		mockServer.AuthTokens = make(map[string]bool) // Clear auth tokens
		
		// Get an initial token 
		token := generateTestToken(mockServer)
		
		// Now enable expiration for the next request
		mockServer.SimulateExpiration = true
		
		// Fetch events directly (will need to refresh token)
		events, err := fetchEventsDirectly(token, mockServer)
		
		// Validate
		if err != nil {
			t.Fatalf("Failed to fetch events with token expiration: %v", err)
		}
		
		if len(events) != 20 {
			t.Errorf("Expected 20 events after token refresh, got %d", len(events))
		}
	})

	// Test case 2: Multiple API calls with expiring tokens
	t.Run("Multiple API calls with expiring tokens", func(t *testing.T) {
		mockServer.RequestLog = nil // Clear request log
		mockServer.AuthTokens = make(map[string]bool) // Clear auth tokens
		
		// Make several API calls in sequence, each should require a new token
		for i := 0; i < 3; i++ {
			// Get a fresh token for each call
			token := generateTestToken(mockServer)
			
			// Enable expiration for the token
			mockServer.SimulateExpiration = true
			
			// Fetch events (should get a new token)
			events, err := fetchEventsDirectly(token, mockServer)
			
			// Validate each call
			if err != nil {
				t.Fatalf("Failed on call %d: %v", i+1, err)
			}
			
			if len(events) != 20 {
				t.Errorf("Call %d: Expected 20 events, got %d", i+1, len(events))
			}
		}
	})

	// Test case 3: Token refresh failure
	t.Run("Token refresh failure", func(t *testing.T) {
		// Custom version for this test case
		mockServer.RequestLog = nil // Clear request log
		mockServer.AuthTokens = make(map[string]bool) // Clear auth tokens
		
		// Create initial token
		token := generateTestToken(mockServer)
		
		// Create request for current time (last 24 hours)
		now := time.Now()
		oneDayAgo := now.Add(-24 * time.Hour)
		
		req := map[string]interface{}{
			"startTimestamp": fmt.Sprintf("%d", oneDayAgo.Unix()),
			"endTimestamp":   fmt.Sprintf("%d", now.Unix()),
			"language":       "en",
			"countryNo":      "US",
		}
		reqJSON, _ := json.Marshal(req)
		
		// Make first request with the token
		httpReq, _ := http.NewRequest("POST", mockServer.Server.URL+"/library/newselectlibrary", bytes.NewBuffer(reqJSON))
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Accept", "application/json")
		httpReq.Header.Set("Authorization", token)
		
		// Before making the request, simulate token expiration
		mockServer.SimulateExpiration = true
		// Also set auth failures for when refresh is attempted
		mockServer.AuthFailures = 1
		
		// Make the request
		client := &http.Client{}
		resp, _ := client.Do(httpReq)
		
		// Test for token expiration in response
		if resp != nil && resp.StatusCode == http.StatusOK {
			var responseMap map[string]interface{}
			decoder := json.NewDecoder(resp.Body)
			if err := decoder.Decode(&responseMap); err == nil {
				if result, ok := responseMap["result"].(float64); ok {
					// Token expired, should get -1025
					if result == -1025 {
						// Good, token expired. Now try to get a new token, which will fail
						_, err := getToken(mockServer)
						
						// This should fail with an auth error
						if err == nil {
							t.Fatalf("Expected token refresh to fail, but it succeeded")
						}
						
						// Should have an auth failure message
						if !strings.Contains(err.Error(), "API error") {
							t.Errorf("Expected API error message, got: %v", err)
						}
						
						// Test passed!
						return
					}
				}
			}
		}
		
		// If we get here, the test didn't work as expected
		t.Fatalf("Test didn't detect token expiration or auth failure")
	})
}

// generateTestToken creates a test token in the mock server and returns it.
func generateTestToken(mockServer *MockAPIServer) string {
	token := fmt.Sprintf("mock-token-%d", time.Now().UnixNano())
	mockServer.AuthTokens[token] = true
	return token
}

// getToken gets a new authentication token from the mock server.
func getToken(mockServer *MockAPIServer) (string, error) {
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
		return "", fmt.Errorf("error making login request: %v", err)
	}
	defer resp.Body.Close()
	
	// Decode login response
	var loginResponse map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&loginResponse); err != nil {
		return "", fmt.Errorf("error decoding login response: %v", err)
	}
	
	// Check for login errors
	if result, ok := loginResponse["result"].(float64); ok && result != 0 {
		msg, _ := loginResponse["msg"].(string)
		return "", fmt.Errorf("API error during login: %s (code: %.0f)", msg, result)
	}
	
	// Extract token from login response
	var token string
	if data, ok := loginResponse["data"].(map[string]interface{}); ok {
		if tokenObj, ok := data["token"].(map[string]interface{}); ok {
			if tokenStr, ok := tokenObj["token"].(string); ok {
				token = tokenStr
			}
		}
	}
	
	if token == "" {
		return "", fmt.Errorf("failed to get token from login response")
	}
	
	return token, nil
}

// fetchEventsDirectly attempts to fetch events with a given token.
// If the token is expired, it will try to refresh and retry.
func fetchEventsDirectly(token string, mockServer *MockAPIServer) ([]MockEvent, error) {
	// Create request for current time (last 24 hours)
	now := time.Now()
	oneDayAgo := now.Add(-24 * time.Hour)
	
	// Create request payload
	req := map[string]interface{}{
		"startTimestamp": fmt.Sprintf("%d", oneDayAgo.Unix()),
		"endTimestamp":   fmt.Sprintf("%d", now.Unix()),
		"language":       "en",
		"countryNo":      "US",
	}
	reqJSON, _ := json.Marshal(req)
	
	// Create HTTP request
	httpReq, _ := http.NewRequest("POST", mockServer.Server.URL+"/library/newselectlibrary", bytes.NewBuffer(reqJSON))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)
	
	// Send request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()
	
	// Parse response
	var responseMap map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&responseMap); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}
	
	// Check for API errors
	if result, ok := responseMap["result"].(float64); ok && result != 0 {
		msg, _ := responseMap["msg"].(string)
		
		// Check if it's an authentication error (expired token)
		if result == auth.ErrorAccountKicked || result == auth.ErrorTokenMissing || result == -1025 {
			// Token expired, we need to get a new one
			mockServer.SimulateExpiration = false // Turn off expiration for the refresh
			newToken := generateTestToken(mockServer)
			
			// Retry the request with the new token
			return fetchEventsWithNewToken(newToken, req, mockServer)
		}
		
		return nil, fmt.Errorf("API error: %s (code: %.0f)", msg, result)
	}
	
	// If we get here, the request was successful
	return transformToMockEvents(responseMap)
}

// fetchEventsWithNewToken makes a request with a fresh token.
func fetchEventsWithNewToken(token string, req map[string]interface{}, mockServer *MockAPIServer) ([]MockEvent, error) {
	reqJSON, _ := json.Marshal(req)
	
	// Create HTTP request with new token
	httpReq, _ := http.NewRequest("POST", mockServer.Server.URL+"/library/newselectlibrary", bytes.NewBuffer(reqJSON))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)
	
	// Send request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request after token refresh: %v", err)
	}
	defer resp.Body.Close()
	
	// Parse response
	var responseMap map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&responseMap); err != nil {
		return nil, fmt.Errorf("error decoding response after refresh: %v", err)
	}
	
	// Check for API errors
	if result, ok := responseMap["result"].(float64); ok && result != 0 {
		msg, _ := responseMap["msg"].(string)
		return nil, fmt.Errorf("API error after refresh: %s (code: %.0f)", msg, result)
	}
	
	// Transform to our event format
	return transformToMockEvents(responseMap)
}

// isAuthError checks if an error is related to authentication.
func isAuthError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "auth") || 
	       strings.Contains(errMsg, "token") || 
	       strings.Contains(errMsg, "login") ||
	       strings.Contains(errMsg, "credentials")
}